package download

import "os"

type options struct {
	Ids          *[]uint64 `short:"i" long:"id"`
	Bookmarks    *string   `short:"b" long:"bookmarks"`
	DownloadList *string   `short:"l" long:"list"`
	InferId      *string   `short:"I" long:"infer-id"`
	Kind         *string   `short:"t" long:"type"`
	Size         *uint     `short:"s" long:"size"`
	OnlyMeta     *bool     `short:"m" long:"only-meta"`
	Tag          *string   `short:"G" long:"tag"`
	FromOffset   *uint64   `short:"F" long:"from"`
	ToOffset     *uint64   `short:"T" long:"to"`
	Private      *bool     `short:"R" long:"private"`
	LowMeta      *bool     `short:"M" long:"low-meta"`
	Ignore       *string   `short:"n" long:"ignore"` // TODO: interactive mode + validation
	Path         *string   `short:"p" long:"path"`
	Password     *string   `short:"P" long:"password"`
}

func Run() {
	if len(os.Args) <= 1 {
		interactive()
	} else {
		nonInteractive()
	}
}
