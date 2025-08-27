package config

import "strings"

// TODO: move to internal/promptuiext

func parseStrings(input string) []string {
	strs := []string{}
	builder := strings.Builder{}
	runes := []rune(input)

	for i := 0; i < len(runes); i++ {
		if i > 0 && runes[i-1] == '\\' {
			builder.WriteRune(runes[i])
		} else if runes[i] == '\\' {
		} else if runes[i] == ',' || runes[i] == ';' || runes[i] == '、' || runes[i] == '；' {
			str := strings.Trim(builder.String(), " 　\t\r\n")
			if str != "" {
				strs = append(strs, str)
			}
			builder.Reset()
		} else {
			builder.WriteRune(runes[i])
		}
	}
	str := strings.Trim(builder.String(), " 　\t\r\n")
	if str != "" {
		strs = append(strs, str)
	}

	return strs
}
