package help

import (
	"fmt"
	"os"
)

const generalHelp = //
`Usage:

Configure:      piximanctl config         # Run in interactive mode
                piximanctl help config    # Run for more information

Download:       piximanctl download       # Run in interactive mode
                piximanctl help download  # Run for more information
`

const configHelp = //
`Usage:         Run without arguments to enter interactive mode.
                piximanctl config [ -sessionid <...> ]
                                  [ -password <...> ]

Description:    Change permanent configuration for piximan. The configured settings will be
                used for all future commands by default.
Options:
    -sessionid  The session ID to use for pixiv.net API autorization.
                You can get this id from browser cookies on https://www.pixiv.net.
                Search for a cookie named "PHPSESSID".
                DO NOT paste the value directly in the command line as it could
                be logged in the terminal history (e.g. ~/.bash_history).
                Session ID will be encrypted and stored in ~/.piximan/sessionid.

    -password   The master password that can be set to encrypt the session ID.
                If omited the password will be set to an empty string.
                Similarly to the session ID, avoid pasting the value directly.

Examples:       Pass the session ID value from the clipboard:
                - On Linux (X11)        piximanctl config -sessionid $(xclip -o)
                - On Linux (Wayland)    piximanctl config -sessionid $(wl-paste)
                - On Windows            piximanctl config -sessionid $(Get-Clipboard)
                - On MacOS              piximanctl config -sessionid $(pbpaste)
`

const downloadHelp = //
`Usage:         Run without arguments to enter interactive mode.
                piximanctl download [ -id <...> ]
                                    [ -type <artwork|novel> ]
                                    [ -path <...> ]
                                    [ -size <0-3> ]
                                    [ -sessionid <...> ]

Description:    Download the work files and metadata from pixiv.net to the given directory.
                Session ID must be configued prior to this command.
Options:
    -id         ID of the downloaded work. You can found it in the work URI:
                https://www.pixiv.net/artworks/12345 <- 12345 is the ID here.

    -type       The type of work to download. Defaults to artwork.
                Available options are:
                - artwork
                - novel

    -size       Size (resolution) of the image to download. This Option doesn't
                apply to ugoira and novels. Defaults to original size.
                Available options are:
                - 0 thumbnail
                - 1 small
                - 2 medium
                - 3 original

    -path       Directory to save the files into. Defaults to the current directory.
                You can use these substitutions in the pathname:
                - {user}    the username of the work author.
                - {title}   the title of the work.
                - {id}      the ID of the work.
                - {userid}  the ID of the work author.
                Be aware that any Windows / NTFS reserved names will be automaticaly
                padded with underscores, reserved characters - replaced and any dots
                or spaces in front or end of the filenames will be trimmed.

    -sessionid  Will default to the session ID stored in config.
                For additional information, run 'piximanctl config'.

    -password   The master password to access the session ID from configuration.
                Avoid pasting the value directly in the terminalas it could be
                logged in the history.

Examples:       piximanctl download -id 12345 -size 1 -path ~/Downloads/work
                piximanctl download -type novel -id 12345 -path ./{user}/{title} -nowin
`

func Run() {
	var section string
	if len(os.Args) > 1 {
		section = os.Args[1]
	}

	switch section {
	case "config":
		RunConfig()
	case "download":
		RunDownload()
	default:
		RunGeneral()
	}
}

func RunGeneral() {
	fmt.Print(generalHelp)
}

func RunConfig() {
	fmt.Print(configHelp)
}

func RunDownload() {
	fmt.Print(downloadHelp)
}
