package commands

import (
	"encoding/json"
)

// SAdd - Will push an element into a set
func SAdd(cmd CmdInterface) {
	addElement(cmd, true)
}

// SMembers - Return all elements of a set
func SMembers(cmd CmdInterface) {
	var list []string
	if value, ok := cmd.State[cmd.Parts[1]]; ok {
		err := json.Unmarshal([]byte(value), &list)
		if err != nil {
			returnNull(cmd)
		}
		returnArr(cmd, list)
	} else {
		returnEmptyArr(cmd)
	}
}

// SIsMember - Return true if element is part of the set
func SIsMember(cmd CmdInterface) {
	var list []string
	if value, ok := cmd.State[cmd.Parts[1]]; ok {
		err := json.Unmarshal([]byte(value), &list)
		if err != nil {
			returnNull(cmd)
		} else {
			var exists bool = false
			for _, val := range list {
				if val == cmd.Parts[2] {
					exists = true
				}
			}
			if exists {
				returnInt(cmd, 1)
			} else {
				returnInt(cmd, 0)
			}
		}
	} else {
		returnInt(cmd, 0)
	}
}

// SRem - Will remove a value from a set
func SRem(cmd CmdInterface) {
	var list []string
	if value, ok := cmd.State[cmd.Parts[1]]; ok {
		err := json.Unmarshal([]byte(value), &list)
		if err != nil {
			returnNull(cmd)
		} else {
			var newList []string
			var removedCount int = 0
			for key, val := range list {
				if val != cmd.Parts[2] {
					newList[key] = val
				} else {
					removedCount++
				}
			}
			if len(newList) == 0 {
				// delete directly the array if empty at this point
				delete(cmd.State, cmd.Parts[1])
			} else {
				jsonEncoding, _ := json.Marshal(newList)
				cmd.State[cmd.Parts[1]] = string(jsonEncoding)
			}
			returnInt(cmd, removedCount)
		}
	} else {
		returnNull(cmd)
	}
}
