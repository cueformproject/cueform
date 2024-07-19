package manifest

import "strings"

#Config: {
	command: *"terraform" | string
	output:  *"./providers" | string

	providers: [S=string]: {
		name:     strings.Split(S, "/")[1]
		source:   S
		version:  string
		filename: *"\(output)/\(source)/\(version)/\(name).gen.cue" | string
	}
}
