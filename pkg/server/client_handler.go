package server

import (
	"net"
	"os"
	"strings"

	"github.com/lefuturiste/keyvaluer/pkg/commands"
	log "github.com/sirupsen/logrus"
)

var state map[string]string
var password string

func handleClient(conn net.Conn) {
	log.Debug("New client: ", conn.RemoteAddr().String())
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
					log.Debug("New cmd: ", components)
					var name string = strings.ToUpper(components[0])

					if password != "" && name != "AUTH" && name != "QUIT" && !authenticated {
						conn.Write([]byte("-NOAUTH Authentication required."))
						conn.Write([]byte("\r\n"))
						name = ""
					} else {
						cmd := commands.CmdInterface{
							Parts: components,
							State: state,
							Conn:  conn,
						}
						if name == "AUTH" {
							authenticated = commands.Auth(cmd)
						} else {
							commandMap[name].(func(commands.CmdInterface))(cmd)
						}
					}
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
