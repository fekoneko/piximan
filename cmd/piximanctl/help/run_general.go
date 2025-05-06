package help

import "fmt"

const GENERAL_HELP = //
`Configure:   piximanctl config         # Run in interactive mode
             piximanctl help config    # Run for more information

Download:    piximanctl download       # Run in interactive mode
             piximanctl help download  # Run for more information
`

func RunGeneral() {
	fmt.Print(GENERAL_HELP)
}
