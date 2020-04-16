package commands

import (
	"os"
	"time"
)

// Ping -
func Ping(cmd CmdInterface) {
	cmd.Conn.Write([]byte("+PONG"))
	endLine(cmd)
}

// Quit -
func Quit(cmd CmdInterface) {
	returnSuccess(cmd)
	cmd.Conn.Close()
}

// Command -
func Command(cmd CmdInterface) {
	cmd.Conn.Write([]byte("*0"))
	endLine(cmd)
}

// Select -
func Select(cmd CmdInterface) {
	returnSuccess(cmd)
}

// Auth -
func Auth(cmd CmdInterface) bool {
	if cmd.Parts[1] == os.Getenv("REQUIRED_PASS") {
		returnSuccess(cmd)
		return true
	} else {
		time.Sleep(1 * time.Second)
		cmd.Conn.Write([]byte("-ERR invalid password"))
		endLine(cmd)
		return false
	}
}
