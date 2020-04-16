package commands

import (
	"encoding/json"
	"fmt"
)

func addElement(cmd CmdInterface, removeTwin bool) {
	var list []string
	if value, ok := cmd.State[cmd.Parts[1]]; ok {
		err := json.Unmarshal([]byte(value), &list)
		if err != nil {
			fmt.Println("Fatal JSON decoding err")
		}
	} else {
		list = make([]string, 0)
	}

	// append the members to the list
	var appended int = 0
	for key, value := range cmd.Parts {
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
	cmd.State[cmd.Parts[1]] = string(jsonEncoding)
	returnInt(cmd, appended)
}

// SAdd - Will push an element into a set
func SAdd(cmd CmdInterface) {
	addElement(cmd, true)
}

// RPush -
func RPush(cmd CmdInterface) {
	addElement(cmd, false)
}

// SMembers - Return all elements of a set
func SMembers(cmd CmdInterface) {
	var list []string
	if value, ok := cmd.State[cmd.Parts[1]]; ok {
		err := json.Unmarshal([]byte(value), &list)
		if err != nil {
			fmt.Println("Fatal JSON decoding err")
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
			fmt.Println("Fatal JSON decoding err")
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

// LPop - Will remove and return the first element of a set
func LPop(cmd CmdInterface) {
	var list []string
	if value, ok := cmd.State[cmd.Parts[1]]; ok {
		err := json.Unmarshal([]byte(value), &list)
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

// SRem - Will remove a value from a set
func SRem(cmd CmdInterface) {
	var list []string
	if value, ok := cmd.State[cmd.Parts[1]]; ok {
		err := json.Unmarshal([]byte(value), &list)
		if err != nil {
			fmt.Println("Fatal JSON decoding err")
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
