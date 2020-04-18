package commands

import (
	"encoding/json"
)

func addElement(cmd CmdInterface, removeTwin bool) {
	var list []string
	if value, ok := cmd.State[cmd.Parts[1]]; ok {
		err := json.Unmarshal([]byte(value), &list)
		if err != nil {
			returnNull(cmd)
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
