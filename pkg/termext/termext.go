package termext

import (
	"os"

	"golang.org/x/sys/unix"
)

var initialTermois *unix.Termios

// Probably need to use unix.TCGETA for MacOS
func DisableInputEcho() {
	if initialTermois != nil {
		return
	}
	initialTermois, _ = unix.IoctlGetTermios(int(os.Stdin.Fd()), unix.TCGETS)
	termois := *initialTermois
	termois.Lflag &^= unix.ECHO
	unix.IoctlSetTermios(unix.Stdout, unix.TCSETS, &termois)
}

// Probably need to use unix.TCGETA for MacOS
func RestoreInputEcho() {
	if initialTermois != nil {
		unix.IoctlSetTermios(unix.Stdout, unix.TCSETS, initialTermois)
	}
}
