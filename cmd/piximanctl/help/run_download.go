package help

import "fmt"

const DOWNLOAD_HELP = //
`Run without arguments to enter interactive mode.

> piximanctl download [--id      ...] [--type ...] [--path     ...]
                      [--list    ...] [--size ...] [--password ...]
                      [--inferid ...] [--onlymeta]

                              Download sources
                              ----------------
--id         ID of the downloaded work. You can found it in the work URI:
 -i          https://www.pixiv.net/artworks/12345 <- 12345 is the ID here.
             Can be provided multiple times.

--list       Path to a file with information about which works to download.
 -l          The file must contain a list in YAML format, for example:
             - id: 12345                  # required
               type: artwork              # defaults to the --type argument
               size: 1                    # defaults to the --size argument
               onlymeta: true             # defaults to the --onlymeta argument
               paths: ['./{userid}/{id}'] # defaults to the --path argument
             - id: 23456
               type: novel

--inferid    Infer the IDs of works from the given path. Useful for updating
 -I          the metadata in existing collection when coupled with -onlymeta flag.
             The path may contain the following patterns:
             - {id}         : the ID of the work - required.
             - *            : matches any sequence of non-separator characters.

                              Download options
                              ----------------
--type       The type of work to download. Defaults to artwork.
 -t          Available options are:
             - artwork      - novel

--size       Size (resolution) of the image to download. This Option doesn't
 -s          apply to ugoira and novels. Defaults to original size.
             Available options are:
             - 0 thumbnail  - 2 medium
             - 1 small      - 3 original

--onlymeta   Only download the metadata.yaml file for the work. Useful for
 -m          updating the metadata of existing works.

                              Other parameters
                              ----------------
--path       Directory to save the files into. Defaults to the current directory
 -p          or the one found with -inferid flag.
             You can use these substitutions in the pathname:
             - {title}      : the title of the work
             - {id}         : the ID of the work
             - {user}       : the username of the work author
             - {userid}     : the ID of the work author
             - {restrict}   : age restriction of the work (all-ages, R-18, R-18G)
             Be aware that any Windows / NTFS reserved names will be automaticaly
             padded with underscores, reserved characters - replaced and any dots
             or spaces in front or end of the filenames will be trimmed.

--password   The master password that is used to decrypt session ID if one has been
 -P          set. If omited, you will be prompted for the password when needed.
             Do not paste the value directly in the command line as it could
             be logged in the terminal history (e.g. ~/.bash_history).

                                  Examples
                                  --------
# Download artwork with ID 10000 to your possible collection directory
> piximanctl download --id 10000 --path "$HOME/My Collection/{userid}/{id}"

# Download novels with ID 10000 and 20000
> piximanctl download --id 10000 --id 20000 --type novel --path "./{userid}/{id}"

# Download works from list.yaml to the current directory with fallback path
> piximanctl download --list "./list.yaml" --path "./{userid}/{id}"

# Update metadata in the collection saved earlier
> piximanctl download --inferid "$HOME/My Collection/*/{id}" --onlymeta
`

func RunDownload() {
	fmt.Print(DOWNLOAD_HELP)
}

// TODO: authorized user bookmarks download (+private) + write about it in the help + example
// TODO: bookmarks --newer, --older than date
// TODO: help and examples for bookmarks download
// TODO: {tag} pattern for user-assigned tags when downloading bookmarks

// TODO: just download user's works ('my' or by id)

// TODO: --log, -L option to log the output to a file (-L should be for language actually)
// TODO: summary about downloaded / not downloaded works at the end of download
// TODO: total download progress at the buttom
// TODO: count errors / warnings

// TODO: choose language to download work metadata in (flag and config)

// TODO: piximanctl dedupe --newer --path './* ({userid})/* ({id})' + interactive mode
//       utility for merging authors and / or works with duplicate IDs
// TODO: --dedupe option for piximanctl download as well
// TODO: replace unused patterns like {title} in --inferid with *, do the same with piximanctl dedupe

// TODO: test if session ID is valid after configuring by sending a request
