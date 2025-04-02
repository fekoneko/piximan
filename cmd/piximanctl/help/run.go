package help

import (
	"fmt"
	"os"
)

const generalHelp = //
`Configure:   piximanctl config         # Run in interactive mode
             piximanctl help config    # Run for more information

Download:    piximanctl download       # Run in interactive mode
             piximanctl help download  # Run for more information
`

const configHelp = //
`Usage:       Run without arguments to enter interactive mode.
             piximanctl config [ -sessionid ... ] [ -password ... ]

Description: Change permanent configuration for piximan. The configured settings
             will be used for all future commands by default.

-sessionid   The session ID to use for pixiv.net API autorization.
             The authorization is used only when it's absolutely required, other
             requests will be made anonymously.
             You can get this ID from browser cookies on https://www.pixiv.net.
             Search for a cookie named 'PHPSESSID'.
             Do not paste the value directly in the command line as it could
             be logged in the terminal history (e.g. ~/.bash_history).
             Session ID will be encrypted and stored in ~/.piximan/sessionid.
             You can remove the session ID by providing an empty string.

-password    The master password that can be set to encrypt the provided
             session ID. If omited the password will be set to an empty string.
             Similarly to the session ID, avoid pasting the value directly.

Examples:    piximanctl config -sessionid $(xclip -o)
             piximanctl config -sessionid $PHPSESSID -password $PASSWORD
`

const downloadHelp = //
`Usage:       Run without arguments to enter interactive mode.
             piximanctl download [ -id        ... ] [ -type     artwork ]
                                 [ -size      0-3 ]             novel
                                 [ -path      ... ] [ -onlymeta         ]

Description: Download the work files and metadata from pixiv.net to the given
             directory. This command does not require a session ID. All the
             requests to Pixiv API will be made anonymously.

-id          ID of the downloaded work. You can found it in the work URI:
             https://www.pixiv.net/artworks/12345 <- 12345 is the ID here.

-type        The type of work to download. Defaults to artwork.
             Available options are:
             - artwork      - novel

-size        Size (resolution) of the image to download. This Option doesn't
             apply to ugoira and novels. Defaults to original size.
             Available options are:
             - 0 thumbnail  - 2 medium
             - 1 small      - 3 original

-path        Directory to save the files into. Defaults to the current directory
             or the one found with -inferid flag.
             You can use these substitutions in the pathname:
             - {title}      : the title of the work.
             - {id}         : the ID of the work.
             - {user}       : the username of the work author.
             - {userid}     : the ID of the work author.
             Be aware that any Windows / NTFS reserved names will be automaticaly
             padded with underscores, reserved characters - replaced and any dots
             or spaces in front or end of the filenames will be trimmed.

-inferid     Infer the IDs of works from the given path. Useful for updating
             the metadata in existing collection when coupled with -onlymeta flag.
             The path may contain the following patterns:
             - {id}         : the ID of the work - required.
             - {any}        : any string not containing path separators.
-onlymeta    Only download the metadata.yaml file for the work. Useful for
             updating the metadata of existing works.

Examples:    piximanctl download -id 12345 -size 1 -password $(xclip -o)
             piximanctl download -id 12345 -type novel -path ./{userid}/{id}
             piximanctl download -inferid ./{any}/{any}{id} -onlymeta
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
