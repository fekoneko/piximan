package help

import "fmt"

const generalHelp = //
`Viewer GUI:     > piximan                # Run the main application

Downloader CLI: > piximan download       # Run in interactive mode
                > piximan help download  # More information

Config CLI:     > piximan config         # Run in interactive mode
                > piximan help config    # More information
`

func RunGeneral() {
	fmt.Print(generalHelp)
}
