package utils

import "os"

func GetFlagValue(flagName string) string {
	for i, arg := range os.Args {
		if arg == "-"+flagName && i+1 < len(os.Args) {
			return os.Args[i+1]
		} else if len(arg) > len(flagName)+1 && arg[:len(flagName)+2] == "-"+flagName+"=" {
			return arg[len(flagName)+2:]
		}
	}
	return ""
}
