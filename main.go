package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "6379"
	CRLF      = "\r\n"
)

var state map[string]string
var password string

func startServer(port string, pass string) {
	if port != "0" && port != "" {
		os.Setenv("PORT", port)
	}
	os.Setenv("REQUIRED_PASS", pass)
	go main()
}

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

type listMemberType []string

func handleClient(conn net.Conn) {
	fmt.Println("New client:", conn.RemoteAddr().String())
	var message bool
	var input string
	var componentIndex int
	var components map[int]string
	var authenticated bool
	var password string = os.Getenv("REQUIRED_PASS")
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

					if password != "" && name != "AUTH" && name != "QUIT" && !authenticated {
						conn.Write([]byte("-NOAUTH Authentication required."))
						conn.Write([]byte(CRLF))
						name = ""
					}

					var removeTwin bool = true
					if name == "RPUSH" {
						name = "SADD"
						removeTwin = false
					}

					switch name {
					case "PING":
						conn.Write([]byte("+PONG"))
						conn.Write([]byte(CRLF))

					case "SELECT":
						success(conn)
						conn.Write([]byte(CRLF))

					case "AUTH":
						fmt.Println("Password:", password)
						if components[1] == password {
							authenticated = true
							success(conn)
						} else {
							authenticated = false
							time.Sleep(1 * time.Second)
							conn.Write([]byte("-ERR invalid password"))
						}
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

					case "SET":
						state[components[1]] = components[2]
						success(conn)
						conn.Write([]byte(CRLF))

					case "EXISTS":
						if _, ok := state[components[1]]; ok {
							conn.Write([]byte(":1"))
						} else {
							conn.Write([]byte(":0"))
						}
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

					case "APPEND":
						if _, ok := state[components[1]]; ok {
							state[components[1]] += components[2]
						} else {
							state[components[1]] = components[2]
						}
						conn.Write([]byte(":" + strconv.Itoa(len(state[components[1]]))))
						conn.Write([]byte(CRLF))

					case "INCR":
						if value, ok := state[components[1]]; ok {
							valueParsed, _ := strconv.ParseInt(value, 10, 64)
							state[components[1]] = strconv.FormatInt(valueParsed+1, 10)
						} else {
							state[components[1]] = components[2]
						}
						conn.Write([]byte(":" + state[components[1]]))
						conn.Write([]byte(CRLF))

					case "INCRBY":
						if value, ok := state[components[1]]; ok {
							valueA, _ := strconv.ParseInt(components[2], 10, 64)
							valueB, _ := strconv.ParseInt(value, 10, 64)
							state[components[1]] = strconv.FormatInt(valueA+valueB, 10)
						} else {
							state[components[1]] = components[2]
						}
						conn.Write([]byte(":" + state[components[1]]))
						conn.Write([]byte(CRLF))

					/**
					 * Will push an element into a set
					 */
					case "SADD":
						var list []string
						if value, ok := state[components[1]]; ok {
							err = json.Unmarshal([]byte(value), &list)
							if err != nil {
								fmt.Println("Fatal JSON decoding err")
							}
						} else {
							list = make([]string, 0)
						}

						// append the members to the list
						var appended int = 0
						for key, value := range components {
							if key != 0 && key != 1 {
								// check if the value to append is already in the list
								var notAppended bool = true
								for _, v := range list {
									if v == value {
										notAppended = false
									}
								}
								if (notAppended && removeTwin) || (!removeTwin) {
									list = append(list, value)
									appended++
								}
							}
						}
						// encode the array as JSON
						jsonEncoding, _ := json.Marshal(list)
						state[components[1]] = string(jsonEncoding)
						conn.Write([]byte(":" + strconv.FormatInt(int64(appended), 10)))
						conn.Write([]byte(CRLF))

					/**
					 * Will return an array of string
					 * Will return all the items of a set
					 */
					case "SMEMBERS":
						var list []string
						if value, ok := state[components[1]]; ok {
							err = json.Unmarshal([]byte(value), &list)
							if err != nil {
								fmt.Println("Fatal JSON decoding err")
							}
							conn.Write([]byte("*" + strconv.Itoa(len(list))))
							conn.Write([]byte(CRLF))
							for _, val := range list {
								conn.Write([]byte("$" + strconv.Itoa(len(val))))
								conn.Write([]byte(CRLF))
								conn.Write([]byte(val))
								conn.Write([]byte(CRLF))
							}
						} else {
							conn.Write([]byte("*0"))
							conn.Write([]byte(CRLF))
						}

					/**
					 * Will return a boolean
					 * True if the key belongs to the set
					 * False if the key don't belongs to the set
					 */
					case "SISMEMBER":
						var list []string
						if value, ok := state[components[1]]; ok {
							err = json.Unmarshal([]byte(value), &list)
							if err != nil {
								fmt.Println("Fatal JSON decoding err")
							} else {
								var exists bool = false
								for _, val := range list {
									if val == components[2] {
										exists = true
									}
								}
								if exists {
									conn.Write([]byte(":1"))
								} else {
									conn.Write([]byte(":0"))
								}
							}
						} else {
							conn.Write([]byte(":0"))
						}
						conn.Write([]byte(CRLF))

					/**
					 * Will remove and return the first element of a set
					 */
					case "LPOP":
						var list []string
						if value, ok := state[components[1]]; ok {
							err = json.Unmarshal([]byte(value), &list)
							if err != nil {
								fmt.Println("Fatal JSON decoding err")
							} else {
								var isFirstNull bool = true
								var first string
								if len(list) > 0 {
									isFirstNull = false
									first = list[0]
									list = list[1:]
								}
								if len(list) == 0 {
									// delete directly the array if empty at this point
									delete(state, components[1])
								} else {
									jsonEncoding, _ := json.Marshal(list)
									state[components[1]] = string(jsonEncoding)
								}
								if isFirstNull {
									conn.Write([]byte("$-1"))
								} else {
									conn.Write([]byte("$" + strconv.Itoa(len(first))))
									conn.Write([]byte(CRLF))
									conn.Write([]byte(first))
								}
							}
						} else {
							conn.Write([]byte("$-1"))
						}
						conn.Write([]byte(CRLF))

					/**
					 * Will remove a value from a set
					 */
					case "SREM":
						var list []string
						if value, ok := state[components[1]]; ok {
							err = json.Unmarshal([]byte(value), &list)
							if err != nil {
								fmt.Println("Fatal JSON decoding err")
							} else {
								var newList []string
								var removedCount int = 0
								for key, val := range list {
									if val != components[2] {
										newList[key] = val
									} else {
										removedCount++
									}
								}
								if len(newList) == 0 {
									// delete directly the array if empty at this point
									delete(state, components[1])
								} else {
									jsonEncoding, _ := json.Marshal(newList)
									state[components[1]] = string(jsonEncoding)
								}
								conn.Write([]byte(":" + strconv.Itoa(removedCount)))
							}
						} else {
							conn.Write([]byte("$-1"))
						}
						conn.Write([]byte(CRLF))

					case "FLUSHALL":
						state = make(map[string]string)
						success(conn)
						conn.Write([]byte(CRLF))

					case "EXPIRE":
						if _, ok := state[components[1]]; ok {
							conn.Write([]byte(":1"))
						} else {
							conn.Write([]byte(":0"))
						}
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

					case "QUIT":
						success(conn)
						conn.Write([]byte(CRLF))
						conn.Close()
						return
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
