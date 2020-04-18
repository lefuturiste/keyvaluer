package commands

import (
	"encoding/json"
)

// RPush -
func RPush(cmd CmdInterface) {
	addElement(cmd, false)
}

// LPop - Will remove and return the first element of a list
func LPop(cmd CmdInterface) {
	var list []string
	if value, ok := cmd.State[cmd.Parts[1]]; ok {
		err := json.Unmarshal([]byte(value), &list)
		if err != nil {
			returnNull(cmd)
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
				delete(cmd.State, cmd.Parts[1])
			} else {
				jsonEncoding, _ := json.Marshal(list)
				cmd.State[cmd.Parts[1]] = string(jsonEncoding)
			}
			if isFirstNull {
				returnNull(cmd)
			} else {
				returnString(cmd, first)
			}
		}
	} else {
		returnNull(cmd)
	}
}

// LPush - Will push an element into a list
func LPush(cmd CmdInterface) {
	addElement(cmd, false)
}
