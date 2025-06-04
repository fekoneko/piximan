package help

import "fmt"

const DOWNLOAD_HELP = //
`Run without arguments to enter interactive mode.

> piximan download [--id        ...] [--type ...] [--tag  ...] [--path     ...]
                      [--bookmarks ...] [--size ...] [--from ...] [--password ...]
                      [--list      ...] [--only-meta] [--to   ...]
                      [--infer-id   ...]              [--low-meta ]

                              Download sources
                              ----------------
--id         ID of the downloaded work. You can found it in the work URI:
 -i          https://www.pixiv.net/artworks/12345 <- 12345 is the ID here.
             Can be provided multiple times.

--bookmarks  Download your bookmarks or the bookmarks of the given user.
 -b          Authorization is required for this source. See 'piximan help config'
             Available options are:
             - my         - download bookmarks of the authorized user
             - <user ID>  - numeric ID of the user to download bookmarks from

--list       Path to a file with information about which works to download.
 -l          The file must contain a list in YAML format, for example:
             - id: 12345                  # required
               type: artwork              # defaults to the --type argument
               size: 1                    # defaults to the --size argument
               only-meta: true             # defaults to the --only-meta argument
               paths: ['./{userid}/{id}'] # defaults to the --path argument
             - id: 23456
               type: novel

--infer-id    Infer the IDs of works from the given path. Useful for updating
 -I          the metadata in existing collection when coupled with -only-meta flag.
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

--only-meta   Only download the metadata.yaml file for the work. Useful for
 -m          updating the metadata of existing works.

                         Bookmarks-specific options
                         --------------------------
--tag        User-assigned tag to filter the bookmarks by. You can see those on the
 -G          bookmarks page. You can only specify one tag or omit to download all.

--from       Crawl bookmarks starting with N'th latest bookmark. Zero-based.
 -F          Omit this option to crawl from the latest bookmark.

--to         Crawl bookmarks up to N'th latest bookmark. Zero-based, non-inclusive.
 -T          Omit this option to crawl up to the oldest bookmark.

--private    Download private bookmarks. You only have access to private bookmarks
 -R          of the authorized user so you probably want to use it with --bookmarks my.
             If not provided, public bookmarks will be downloaded instead.

--low-meta    Specify to skip fetching the full metadata for each work. This will
 -M          significantly reduce the number of pixiv.net API calls.
             These options will be missing in the metadata.yaml files:
             - original, views , bookmarks, likes, comments, uploaded,
             - series_id, series_title, series_order
             When downloading novels without --low-meta flag, the full metadata will be
             downloaded without any request overhead, so --low-meta should be omitted.

                              Other parameters
                              ----------------
--path       Directory to save the files into. Defaults to the current directory
 -p          or the one found with -infer-id flag.
             You can use these substitutions in the pathname:
             - {title}      : the title of the work
             - {id}         : the ID of the work
             - {user}       : the username of the work author
             - {userid}     : the ID of the work author
             - {type}       : the type of the work (Illustrations, Manga, Ugoira, Novels)
             - {restrict}   : age restriction of the work (All Ages, R-18, R-18G)
             - {ai}         : whether the work is worth your attention (Human, AI)
             - {original}   : whether the work is original (Original, Not Original)
             - {series}     : the title of the series the work belongs to
             - {seriesid}   : the ID of the series the work belongs to
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
> piximan download --id 10000 --path "$HOME/My Collection/{userid}/{id}"

# Download novels with ID 10000 and 20000
> piximan download --id 10000 --id 20000 --type novel --path "./{userid}/{id}"

# Download 101th to 200th of your latest public artwork bookmarks with full metadata
> piximan download --bookmarks my --from 100 --to 200 --path "./{id}"

# Download all of your private artwork bookmarks with partial metadata saving requests
> piximan download --bookmarks my --private --low-meta --path "./{id}"

# Download all novel bookmarks from user 10000 with tag 'お気に入り'
> piximan download --bookmarks 10000 --type novel --tag "お気に入り" --path "./{id}"

# Download works from list.yaml to the current directory with fallback path
> piximan download --list "./list.yaml" --path "./{userid}/{id}"

# Update metadata in the collection saved earlier
> piximan download --infer-id --only-meta "$HOME/My Collection/*/{id}"
`

func RunDownload() {
	fmt.Print(DOWNLOAD_HELP)
}

// TODO: bookmarks --newer, --older than date
// TODO: download user's works ('my' or by id)

// TODO: --log, -L option to log the output to a file (-L should be reserved for language actually)
// TODO: summary about downloaded / not downloaded works at the end of download
// TODO: total download progress at the buttom
// TODO: count errors / warnings

// TODO: choose language to download work metadata in (flag and config)

// TODO: piximan dedupe --newer --path './* ({userid})/* ({id})' + interactive mode
//       utility for merging authors and / or works with duplicate IDs
// TODO: --dedupe option for piximan download as well
// TODO: replace unused patterns like {title} in --infer-id with *, do the same with piximan dedupe

// TODO: test if session ID is valid after configuring by sending a request
// TODO: check if session ID is set in interractive download mode and prompt user if they want to set it

// TODO: look at app api ->
//       https://github.com/akameco/pixiv-app-api/blob/d153118b62da1e1f17c8287d7f73ce72848aaaf9/src/index.ts#L154
//       https://hanshsieh.github.io/pixiv-api-doc/
//       https://github.com/piglig/pixiv-token/blob/main/pixiv_token_fetcher.py

// TODO: multilpe sources in one command:
//       piximan download --id 12345 --path "./{id}" \
//                        --bookmarks my --path "./{id}"
//       or in list.yaml
