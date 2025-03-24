package flagext

import "flag"

var seenFlags map[string]bool

func Provided(flagName string) bool {
	if !flag.Parsed() {
		return false
	}

	if seenFlags == nil {
		seenFlags = make(map[string]bool)
		flag.Visit(func(f *flag.Flag) { seenFlags[f.Name] = true })
	}

	return seenFlags[flagName]
}
