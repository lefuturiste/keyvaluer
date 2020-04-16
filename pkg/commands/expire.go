package commands

// Expire -
func Expire(cmd CmdInterface) {
	if _, ok := cmd.State[cmd.Parts[1]]; ok {
		returnInt(cmd, 1)
	} else {
		returnInt(cmd, 0)
	}
}
