package download

import "os"

type options struct {
	Ids        *[]uint64 `short:"i" long:"id"`
	Bookmarks  *string   `short:"b" long:"bookmarks"`
	Lists      *[]string `short:"l" long:"list"`
	InferIds   *[]string `short:"I" long:"infer-id"`
	Kind       *string   `short:"t" long:"type"`
	Size       *uint     `short:"s" long:"size"`
	Language   *string   `short:"L" long:"language"`
	OnlyMeta   *bool     `short:"m" long:"only-meta"`
	Rules      *[]string `short:"r" long:"rules"`
	Skips      *[]string `short:"S" long:"skip"`
	Tags       *[]string `short:"G" long:"tag"`
	FromOffset *uint64   `short:"F" long:"from"`
	ToOffset   *uint64   `short:"T" long:"to"`
	Private    *bool     `short:"R" long:"private"`
	LowMeta    *bool     `short:"M" long:"low-meta"`
	UntilSkip  *bool     `short:"U" long:"until-skip"`
	Paths      *[]string `short:"p" long:"path"`
	Password   *string   `short:"P" long:"password"`
}

func Run() {
	if len(os.Args) <= 1 {
		interactive()
	} else {
		nonInteractive()
	}
}
