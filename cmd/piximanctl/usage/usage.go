package usage

import "fmt"

const generalUsage = //
`Usage:

Configure:  piximanctl config [-sessionid <sessionid>]
            Run 'piximanctl config' for more information.

Download:   piximanctl download -id <id> [-path <path>]
            Run 'piximanctl download' for more information.
`

const configUsage = //
`Usage:
    piximanctl config [-sessionid <sessionid>]

Description:
    Configure piximan.

Options:
    -sessionid  The session ID to use for pixiv.net API autorization.
                You can get this id from browser cookies on https://www.pixiv.net.
                Search for a cookie named "PHPSESSID".
                DO NOT paste the value directly in the command line as it could
                be logged in the terminal history (e.g. ~/.bash_history).
                Session ID will be stored in ~/.piximan/sessionid.

Examples:
    Pass the session ID value from the clipboard:
    - On Linux (X11)        piximanctl config -sessionid $(xclip -o)
    - On Linux (Wayland)    piximanctl config -sessionid $(wl-paste)
    - On Windows            piximanctl config -sessionid $(Get-Clipboard)
    - On MacOS              piximanctl config -sessionid $(pbpaste)
`

const downloadUsage = //
`Usage:
    piximanctl download -id <id> [-path <path>]

Description:
    Download the work files and metadata from pixiv.net to the given directory.
    Session ID must be configued prior to this command.

Options:
    -id     ID of the downloaded work. You can found it in the work URI:
            https://www.pixiv.net/artworks/12345 <- 12345 is the ID here.
    -path   The directory to save the files into. Defaults to the current directory.
            You can use this substitutions in the pathname:
            - {user}    the username of the work author.
            - {title}   the title of the work.
            - {id}      the ID of the work.
            - {userid}  the ID of the work author.

Examples:
    piximanctl download -id 12345 -path ~/Downloads/work
    piximanctl download -id 12345 -path "./{user} ({userid})/{title} ({id})"
`

func General() {
	fmt.Print(generalUsage)
}

func Config() {
	fmt.Print(configUsage)
}

func Download() {
	fmt.Print(downloadUsage)
}
