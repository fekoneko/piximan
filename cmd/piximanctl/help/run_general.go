package help

import "fmt"

const generalHelp = //
`Configure:   piximanctl config         # Run in interactive mode
             piximanctl help config    # Run for more information

Download:    piximanctl download       # Run in interactive mode
             piximanctl help download  # Run for more information
`

func RunGeneral() {
	fmt.Print(generalHelp)
}
