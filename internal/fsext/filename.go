package fsext

import "strings"

var filenameReplacer = strings.NewReplacer(
	"/", "／", "\\", "＼", ":", "：",
	"*", "＊", "?", "？", "<", "＜",
	">", "＞", "|", "｜", "\"", "＂",
	"\x00", "", "\x01", "", "\x02", "", "\x03", "", "\x04", "", "\x05", "", "\x06", "",
	"\x07", "", "\x08", "", "\x09", "", "\x0a", "", "\x0b", "", "\x0c", "", "\x0d", "",
	"\x0e", "", "\x0f", "", "\x10", "", "\x11", "", "\x12", "", "\x13", "", "\x14", "",
	"\x15", "", "\x16", "", "\x17", "", "\x18", "", "\x19", "", "\x1a", "", "\x1b", "",
	"\x1c", "", "\x1d", "", "\x1e", "", "\x1f", "",
)

func ToValidFilename(filename string) string {

	switch strings.ToUpper(filename) {
	case ".", "..", "CON", "PRN", "AUX", "NUL",
		"COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9",
		"LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9":
		return "_" + filename
	default:
		filename := filenameReplacer.Replace(filename)
		return strings.Trim(filename, ". ")
	}
}
