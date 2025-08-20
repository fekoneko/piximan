package help

import "fmt"

const APP_HELP = //
`> piximan      # Run piximan collection viewer
> piximan app  # Allows to provide arquments to the viewer
`

func RunApp() {
	fmt.Print(APP_HELP)
}
