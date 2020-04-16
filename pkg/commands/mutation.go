package commands

import (
	"strconv"
)

// Set -
func Set(cmd CmdInterface) {
	cmd.State[cmd.Parts[1]] = cmd.Parts[2]
	returnSuccess(cmd)
}

// Del -
func Del(cmd CmdInterface) {
	if _, ok := cmd.State[cmd.Parts[1]]; ok {
		delete(cmd.State, cmd.Parts[1])
		returnInt(cmd, 1)
	} else {
		returnInt(cmd, 0)
	}
}

// Append -
func Append(cmd CmdInterface) {
	if _, ok := cmd.State[cmd.Parts[1]]; ok {
		cmd.State[cmd.Parts[1]] += cmd.Parts[2]
	} else {
		cmd.State[cmd.Parts[1]] = cmd.Parts[2]
	}
	returnInt(cmd, len(cmd.State[cmd.Parts[1]]))
}

// Incr -
func Incr(cmd CmdInterface) {
	if value, ok := cmd.State[cmd.Parts[1]]; ok {
		valueParsed, _ := strconv.ParseInt(value, 10, 64)
		cmd.State[cmd.Parts[1]] = strconv.FormatInt(valueParsed+1, 10)
	} else {
		cmd.State[cmd.Parts[1]] = cmd.Parts[2]
	}
	returnIntFromStr(cmd, cmd.State[cmd.Parts[1]])
}

// IncrBy -
func IncrBy(cmd CmdInterface) {
	if value, ok := cmd.State[cmd.Parts[1]]; ok {
		valueA, _ := strconv.ParseInt(cmd.Parts[2], 10, 64)
		valueB, _ := strconv.ParseInt(value, 10, 64)
		cmd.State[cmd.Parts[1]] = strconv.FormatInt(valueA+valueB, 10)
	} else {
		cmd.State[cmd.Parts[1]] = cmd.Parts[2]
	}
	returnIntFromStr(cmd, cmd.State[cmd.Parts[1]])
}

// FlushAll -
func FlushAll(cmd CmdInterface) {
	cmd.State = make(map[string]string)
	returnSuccess(cmd)
}
