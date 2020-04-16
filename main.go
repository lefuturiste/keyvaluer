package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "6379"
	CRLF      = "\r\n"
)

var state map[string]string

func main() {
	var host string = os.Getenv("HOST")
	var port string = os.Getenv("PORT")
	if host == "" {
		host = CONN_HOST
	}
	if port == "" {
		port = CONN_PORT
	}
	l, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()
	fmt.Println("Listening on " + host + ":" + port)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		go handleClient(conn)
	}
}

func parseCommand(input string) map[int]string {
	//var input = `SET hello-world '{"random": true}'`
	var cmdComponents = strings.Split(input, " ")
	var quotedComponent = ""
	var inQuoted bool = false
	var components map[int]string = make(map[int]string)
	var componentIndex int = 0
	for _, value := range cmdComponents {
		if strings.Contains(value, "'") && value[0:1] == "'" {
			// we encountered a new quoted component
			// fmt.Println("Start of quoted ", value)
			if value[len(value)-1:] == "'" {
				components[componentIndex] = value[1 : len(value)-1]
				componentIndex++
			} else {
				inQuoted = true
				quotedComponent = value
			}
		} else if strings.Contains(value, "'") && value[len(value)-1:] == "'" {
			// fmt.Println("End of quoted ", value)
			// we reached the end of the quoted component
			quotedComponent += value
			quotedComponent = quotedComponent[1 : len(quotedComponent)-1]
			inQuoted = false
			components[componentIndex] = quotedComponent
			componentIndex++
		} else if inQuoted {
			quotedComponent += value
		} else {
			components[componentIndex] = value
			componentIndex++
		}
	}

	return components
}

func success(conn net.Conn) {
	conn.Write([]byte("+OK"))
}

func error(conn net.Conn, msg string) {
	conn.Write([]byte("-ERR " + msg))
}

func handleClient(conn net.Conn) {
	fmt.Println("New client:", conn.RemoteAddr().String())
	var message bool
	var input string
	var componentIndex int
	var components map[int]string
	if state == nil {
		state = make(map[string]string)
	}
	for {
		message = true
		buf := make([]byte, 1024)
		_, err := conn.Read(buf)
		if err != nil {
			conn.Close()
			break
		}
		for index, value := range buf {
			if 0 == value {
				// 10 mean linefeed, so end of the line
				input += string(buf)[0:index]

				if strings.Count(input, "\r\n") > 1 {
					componentIndex = 0
					components = make(map[int]string)
					for _, component := range strings.Split(input, "\r\n") {
						if len(component) == 0 {
							break
						}
						if !(component[0:1] == "*" || component[0:1] == "$") {
							components[componentIndex] = component
							componentIndex++
						}
					}
				} else {
					components = parseCommand(input[0 : len(input)-1])
				}
				if len(components) > 0 {
					fmt.Println("New cmd: ", components)
					var name string = strings.ToUpper(components[0])
					switch name {
					case "AUTH":
						success(conn)
						conn.Write([]byte(CRLF))

					case "GET":
						if val, ok := state[components[1]]; ok {
							conn.Write([]byte("$" + strconv.Itoa(len(val))))
							conn.Write([]byte(CRLF))
							conn.Write([]byte(val))
							conn.Write([]byte(CRLF))
						} else {
							conn.Write([]byte("$-1"))
							conn.Write([]byte(CRLF))
						}

					case "EXISTS":
						if _, ok := state[components[1]]; ok {
							conn.Write([]byte(":1"))
						} else {
							conn.Write([]byte(":0"))
						}
						conn.Write([]byte(CRLF))

					case "APPEND":
						if _, ok := state[components[1]]; ok {
							state[components[1]] += components[2]
						} else {
							state[components[1]] = components[2]
						}
						conn.Write([]byte(":" + strconv.Itoa(len(state[components[1]]))))
						conn.Write([]byte(CRLF))

					case "SET":
						state[components[1]] = components[2]
						success(conn)
						conn.Write([]byte(CRLF))

					case "DEL":
						if _, ok := state[components[1]]; ok {
							delete(state, components[1])
							conn.Write([]byte(":1"))
							conn.Write([]byte(CRLF))
						} else {
							conn.Write([]byte(":0"))
							conn.Write([]byte(CRLF))
						}

					case "FLUSHALL":
						state = make(map[string]string)
						success(conn)
						conn.Write([]byte(CRLF))

					case "KEYS":
						// TODO: Implement keys command with pattern
						conn.Write([]byte("*" + strconv.Itoa(len(state))))
						conn.Write([]byte(CRLF))
						for key := range state {
							conn.Write([]byte("$" + strconv.Itoa(len(key))))
							conn.Write([]byte(CRLF))
							conn.Write([]byte(key))
							conn.Write([]byte(CRLF))
						}

					case "DBSIZE":
						conn.Write([]byte(":" + strconv.Itoa(len(state))))
						conn.Write([]byte(CRLF))

					case "COMMAND":
						// TODO: Implement COMMAND
						conn.Write([]byte("*0"))
						conn.Write([]byte(CRLF))

					case "PING":
						conn.Write([]byte("+PONG"))
						conn.Write([]byte(CRLF))
					}
					// var output string = ""
					// for _, value := range commandParsed {
					// 	output += value + " "
					// }
					// fmt.Println("cmd: " + output)
					// //conn.Write([]byte(output + "\n"))
				}

				input = ""
				message = false
				break
			}
		}
		if message {
			input = input + string(buf)
		}
	}
}
