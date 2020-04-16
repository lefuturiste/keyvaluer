package commands

import (
	"net"
	"strconv"
)

// CmdInterface -
type CmdInterface struct {
	Parts map[int]string
	State map[string]string
	Conn  net.Conn
}

func returnSuccess(cmd CmdInterface) {
	cmd.Conn.Write([]byte("+OK"))
	endLine(cmd)
}

func returnInt(cmd CmdInterface, numb int) {
	cmd.Conn.Write([]byte(":" + strconv.Itoa(numb)))
	endLine(cmd)
}

func returnIntFromStr(cmd CmdInterface, numb string) {
	cmd.Conn.Write([]byte(":" + numb))
	endLine(cmd)
}

func returnString(cmd CmdInterface, str string) {
	cmd.Conn.Write([]byte("$" + strconv.Itoa(len(str))))
	endLine(cmd)
	cmd.Conn.Write([]byte(str))
	endLine(cmd)
}

func returnNull(cmd CmdInterface) {
	cmd.Conn.Write([]byte("$-1"))
	endLine(cmd)
}

func returnEmptyArr(cmd CmdInterface) {
	cmd.Conn.Write([]byte("*0"))
	endLine(cmd)
}

func returnArr(cmd CmdInterface, arr []string) {
	cmd.Conn.Write([]byte("*" + strconv.Itoa(len(arr))))
	endLine(cmd)
	for _, val := range arr {
		returnString(cmd, val)
	}
}

func endLine(cmd CmdInterface) {
	cmd.Conn.Write([]byte("\r\n"))
}
