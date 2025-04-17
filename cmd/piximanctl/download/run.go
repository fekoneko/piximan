package download

import "os"

type options struct {
	Id       *uint64 `short:"i" long:"id"`
	King     *string `short:"t" long:"type"`
	Size     *uint   `short:"s" long:"size"`
	Path     *string `short:"p" long:"path"`
	InferId  *string `short:"I" long:"inferid"`
	OnlyMeta *bool   `short:"m" long:"onlymeta"`
}

func Run() {
	if len(os.Args) <= 1 {
		interactive()
	} else {
		nonInteractive()
	}
}
