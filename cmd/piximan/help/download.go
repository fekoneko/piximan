package help

import (
	"fmt"
)

// TODO: split download help into multiple categories same as `piximan download rules`
//       make `piximan download` only give quick overview and link to the categories

const downloadHelp = //
`Run without arguments to enter interactive mode.

> piximan download [--id        ...] [--type  ...] [--tag    ...] [--path     ...]
                   [--bookmarks ...] [--size  ...] [--from   ...] [--password ...]
                   [--list      ...] [--only-meta] [--to     ...]
                   [--infer-id  ...] [--rules ...] [--low-meta  ]
                                     [--skip  ...] [--until-skip]

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
             - id: 12345                   # required
               type: artwork               # defaults to the --type argument
               size: 1                     # defaults to the --size argument
               only-meta: true             # defaults to the --only-meta argument
               paths: ['./{user-id}/{id}'] # defaults to the --path argument
             - id: 23456
               type: novel
             May be provided multiple times.

--infer-id   Infer the IDs of works from the given path. Useful for updating
 -I          the metadata in existing collection when coupled with -only-meta flag.
             If this flag is provided, --path can be omitted and downloaded works will replace
             ones in the original location. The path may contain the following patterns:
             - {id}         : the ID of the work
             - *            : matches any sequence of non-separator characters
             - {<anything>} : will be treated as *
             If the provided argument doesn't contain any patterns, piximan will recursively
             look for metadata.yaml files in the provided directory and infer IDs and work types
             from there. May be provided multiple times.

                              Download options
                              ----------------
--type       The type of work to download. Defaults to artwork.
 -t          Available options are:
             - artwork      - novel

--size       Size (resolution) of downloaded images. This Option doesn't apply to ugoira.
 -s          Defaults to original size.
             Available options are:
             - 0 thumbnail  - 2 medium
             - 1 small      - 3 original

--language   Japanese is the default language for tags and work titles / descriptions on pixiv,
 -L          translations may also be provided. Available values are:
             - ja - Japanese (original language), applies to tags and descriptions / titles
             - en - English, applies to tags and work descriptions / titles
             - zh - Chinese, applies to tags, descriptions / titles will use the original language
             - ko - Korean, applies to tags, descriptions / titles will use the original language
             The default option is ja, but this can be configured, see 'piximan help config'

--only-meta  Only download the metadata.yaml file for the work. Useful for
 -m          updating the metadata of existing works.

--rules      Path to YAML file with download rules. The download rules are used to
 -r          filter wich works should be downloaded. Run 'piximan help rules' for more info.
             May be provided multiple times.

--skip       All works already present in the provided directory will be skipped when downloading.
 -S          The search is recursive. If you don't use metadata.yaml files in your collection,
             you can also provide infer ID pattern here (see --infer-id). Note that this way, the
             type of the inferred works will be assumed to be the same as provided with --type flag.
             Can be provided multiple times.

                         Bookmarks-specific options
                         --------------------------
--tag        User-assigned tags to filter the bookmarks by. You can see these on the bookmarks
 -G          page. You may specify this option multiple times or omit to download all.

--from       Crawl bookmarks starting with N'th latest bookmark. Zero-based.
 -F          Omit this option to crawl from the latest bookmark.

--to         Crawl bookmarks up to N'th latest bookmark. Zero-based, non-inclusive.
 -T          Omit this option to crawl up to the oldest bookmark.

--private    Download private bookmarks. You only have access to private bookmarks
 -R          of the authorized user so you probably want to use it with --bookmarks my.
             If not provided, public bookmarks will be downloaded instead.

--low-meta   Specify to skip fetching the full metadata for each work. This will
 -M          significantly reduce the number of pixiv.net API calls.
             These options will be missing in the metadata.yaml files:
             - original, views , bookmarks, likes, comments, uploaded,
             - series_id, series_title, series_order
             When downloading novels without --low-meta flag, the full metadata will be
             downloaded without any request overhead, so --low-meta should be omitted.

--until-skip Useful if you already have all of your bookmarks downloaded in the collection and
 -U          only want to sync new ones. This option tells the downloader to stop crawling
             new bookmark pages once it encounters a fully skipped one. This may greatly reduce
             the number of authorized requests to pixiv.net.
             May only be used when coupled with --skip flag.

                              Other parameters
                              ----------------
--path       Directory to save the files into.
 -p          Defaults to current directory or may be inferred from provided --infer-id.
             You can use these substitutions in the pathname:
             - {title}       : the title of the work
             - {id}          : the ID of the work
             - {user}        : the username of the work author
             - {user-id}     : the ID of the work author
             - {type}        : the type of the work (Illustrations, Manga, Ugoira, Novels)
             - {restriction} : age restriction of the work (All Ages, R-18, R-18G)
             - {ai}          : the humanity is doomed (Human, AI)
             - {original}    : whether the work is original (Original, Not Original)
             - {series}      : the title of the series the work belongs to
             - {series-id}   : the ID of the series the work belongs to
             Be aware that any Windows / NTFS reserved names will be automaticaly
             padded with underscores, reserved characters - replaced and any dots
             or spaces in front or end of the filenames will be trimmed.
             May be specified multiple times.

--password   The master password that is used to decrypt session ID if one has been
 -P          set. If omited, you will be prompted for the password when needed.
             Do not paste the value directly in the command line as it could
             be logged in the terminal history (e.g. ~/.bash_history).

                                  Examples
                                  --------
# Download artwork with ID 10000 to your possible collection directory
> piximan download --id 10000 --path "$HOME/My Collection/{user-id}/{id}"

# Download novels with ID 10000 and 20000
> piximan download --id 10000 --id 20000 --type novel --path './{user-id}/{id}'

# Download 101th to 200th of your latest public artwork bookmarks with full metadata
> piximan download --bookmarks my --from 100 --to 200 --path './{id}'

# Download all of your private artwork bookmarks with partial metadata saving requests
> piximan download --bookmarks my --private --low-meta --path './{id}'

# Download all novel bookmarks from user 10000 with tag 'お気に入り'
> piximan download --bookmarks 10000 --type novel --tag 'お気に入り' --path './{id}'

# Download works from list.yaml to the current directory with fallback path
> piximan download --list './list.yaml' --path './{user-id}/{id}'

# Update metadata for artworks in the collection saved earlier
> piximan download --infer-id --only-meta "$HOME/My Collection/*/{id}"

# Use download rules to filter user's bookmarks
> piximan download --bookmarks 12345 --rules './rules.yaml' --path './{user-id}/{id}'

# Sync your bookmarks with the existing collection
> piximan download --bookmarks my --skip '.' --until-skip --path './{user-id}/{id}'
`

func RunDownload() {
	fmt.Print(downloadHelp)
}

// TODO: download user's works ('my' or by id)
// TODO: --save-list option to only save crawl results as a yaml list
// TODO: --log, -o option to log the output to a file (-L should be reserved for language actually)

// TODO: choose language to download work metadata in (flag and config)

// TODO: piximan sort - sort based on metadata

// TODO: test if session ID is valid after configuring by sending a request
// TODO: check if session ID is set in interractive download mode and prompt user if they want to set it
// TODO: detailed descriptions in interactive mode

// TODO: multilpe sources in one command:
//       piximan download --id 12345 --path "./{id}" -- \
//                        --bookmarks my --path "./{id}"

// TODO: configure substitution words
