package help

import "fmt"

const appHelp = //
`> piximan      # Run piximan collection viewer
> piximan app  # Allows to provide arquments to the viewer
`

func RunApp() {
	fmt.Print(appHelp)
}
