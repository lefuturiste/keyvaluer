package commands

import "strconv"

// Get -
func Get(cmd CmdInterface) {
	if val, ok := cmd.State[cmd.Parts[1]]; ok {
		returnString(cmd, val)
	} else {
		returnNull(cmd)
	}
}

// Exists -
func Exists(cmd CmdInterface) {
	if _, ok := cmd.State[cmd.Parts[1]]; ok {
		returnInt(cmd, 1)
	} else {
		returnInt(cmd, 0)
	}
}

// Keys -
func Keys(cmd CmdInterface) {
	// TODO: Implement keys command with pattern
	cmd.Conn.Write([]byte("*" + strconv.Itoa(len(cmd.State))))
	endLine(cmd)
	for key := range cmd.State {
		returnString(cmd, key)
	}
}

// DbSize -
func DbSize(cmd CmdInterface) {
	returnInt(cmd, len(cmd.State))
}
