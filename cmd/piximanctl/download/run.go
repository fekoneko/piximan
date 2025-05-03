package download

import "os"

type options struct {
	Ids      *[]uint64 `short:"i" long:"id"`
	King     *string   `short:"t" long:"type"`
	Size     *uint     `short:"s" long:"size"`
	Path     *string   `short:"p" long:"path"`
	InferId  *string   `short:"I" long:"inferid"`
	OnlyMeta *bool     `short:"m" long:"onlymeta"`
	Password *string   `short:"P" long:"password"`
}

func Run() {
	if len(os.Args) <= 1 {
		interactive()
	} else {
		nonInteractive()
	}
}
