package termext

import (
	"os"

	"golang.org/x/sys/unix"
)

var initialTermois *unix.Termios

// TODO: Test for MacOS

func DisableInputEcho() {
	if initialTermois != nil {
		return
	}
	initialTermois, _ = unix.IoctlGetTermios(int(os.Stdin.Fd()), unix.TIOCGETA)
	termois := *initialTermois
	termois.Lflag &^= unix.ECHO
	unix.IoctlSetTermios(unix.Stdout, unix.TIOCGETA, &termois)
}

func RestoreInputEcho() {
	if initialTermois != nil {
		unix.IoctlSetTermios(unix.Stdout, unix.TIOCGETA, initialTermois)
	}
}
