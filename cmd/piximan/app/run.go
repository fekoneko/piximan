package app

import (
	"fmt"
)

const MESSAGE = //
`GUI for this application is not yet implemented.
- To use the available downloader CLI, run 'piximan download'
- Run 'piximan help downloader' or 'piximan help config' to see usage
- Follow the project on https://github.com/fekoneko/piximan to see the updates`

func Run() {
	fmt.Println(MESSAGE)
}
