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
    piximanctl download -id <id> [-path <path> -sessionid <sessionid>]

Description:
    Download the work files and metadata from pixiv.net to the given directory.
    Session ID must be configued prior to this command.

Options:
    -type       The type of work to download. Defaults to artwork.
                Available options are:
                - artwork
                - novel

    -id         ID of the downloaded work. You can found it in the work URI:
                https://www.pixiv.net/artworks/12345 <- 12345 is the ID here.

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

    -sessionid  Will default to the session ID stored in config.
                For additional information, run 'piximanctl config'.

Examples:
    piximanctl download -id 12345 -size 1 -path ~/Downloads/work
    piximanctl download -type novel -id 12345 -path ./{user}/{title}
`

func RunGeneral() {
	fmt.Print(generalUsage)
}

func RunConfig() {
	fmt.Print(configUsage)
}

func RunDownload() {
	fmt.Print(downloadUsage)
}
