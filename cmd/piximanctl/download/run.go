package download

import "os"

type options struct {
	Ids      *[]uint64 `short:"i" long:"id"`
	Kind     *string   `short:"t" long:"type"`
	Size     *uint     `short:"s" long:"size"`
	Path     *string   `short:"p" long:"path"` // TODO: Paths array
	List     *string   `short:"l" long:"list"` // TODO: interactive and non-interactive
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
