package app

import (
	"fmt"
)

const message = //
`GUI for this application is not yet implemented.
- To use available downloader CLI, run 'piximan download'
- Run 'piximan help download' or 'piximan help config' to see usage
- Follow the project on https://github.com/fekoneko/piximan to see updates
`

func Run() {
	fmt.Print(message)
}
