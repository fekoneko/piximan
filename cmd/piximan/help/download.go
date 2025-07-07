package help

import (
	"fmt"
	"os"
)

// TODO: split download help into multiple categories same as `piximan download rules`
//       make `piximan download` only give quick overview and link to the categories

const downloadHelp = //
`Run without arguments to enter interactive mode.

> piximan download [--id        ...] [--type  ...]      [--tag   ...] [--path     ...]
                   [--bookmarks ...] [--size  ...]      [--from  ...] [--password ...]
                   [--list      ...] [--only-meta]      [--to    ...]
                   [--infer-id  ...] [--rules ...]      [--low-meta ]
                                     [--collection ...] [--fresh    ]

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

--infer-id   Infer the IDs of works from the given path. Useful for updating
 -I          the metadata in existing collection when coupled with -only-meta flag.
             The path may contain the following patterns:
             - {id}         : the ID of the work - required.
             - *            : matches any sequence of non-separator characters.
             - {<anything>} : will be treated as *.

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

--only-meta  Only download the metadata.yaml file for the work. Useful for
 -m          updating the metadata of existing works.

--rules      Path to YAML file with download rules. The download rules are used to
 -r          filter wich works should be downloaded. Run 'piximan help rules' for more info.

--collection If provided, all works already present in the collection directory will
 -c          be skipped when downloading.

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

--low-meta   Specify to skip fetching the full metadata for each work. This will
 -M          significantly reduce the number of pixiv.net API calls.
             These options will be missing in the metadata.yaml files:
             - original, views , bookmarks, likes, comments, uploaded,
             - series_id, series_title, series_order
             When downloading novels without --low-meta flag, the full metadata will be
             downloaded without any request overhead, so --low-meta should be omitted.

--fresh      Useful if you already have all of your bookmarks downloaded in the collection and
 -f          only want to sync the new ones. This option tells the downloader to stop crawling
             new bookmark pages once it encounters a fully ignored one. This may greatly reduce
             the number of authorized requests to pixiv.net.

                              Other parameters
                              ----------------
--path       Directory to save the files into. Defaults to the current directory
 -p          or the one found with -infer-id flag.
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

# Use download rulesto filter user's bookmarks
> piximan download --bookmarks 12345 --rules './rules.yaml' --path './{user-id}/{id}'

# Sync your bookmarks with the existing collection
> piximan download --bookmarks my --collection '.' --fresh --path './{user-id}/{id}'
`

const downloadRulesHelp = //
`Rules are used to filter wich works should be downloaded and defined in YAML format.
All rules are optional, and if multiple rules are defined, the work should match all of
them to be downloaded (AND). Any array matches any of its elements (OR).

If you download bookmarks with --low-meta flag, be aware that rules, related to the
missing metadata fields will be ignored. You need to download full metadata to match them.

Here's an example of all available rules:

  ids:                       [12345, 23456]
  not_ids:                   [34567, 45678]
  title_contains:            ['cute', 'cat']
  title_not_contains:        ['ugly', 'dog']
  title_regexp:              '^.*[0-9]+$'
  kinds:                     ['illust', 'manga', 'ugoira', 'novel']
  description_contains:      ['hello', 'world']
  description_not_contains:  ['goodbye', 'universe']
  description_regexp:        '^.*[0-9]+$'
  user_ids:                  [12345, 23456]
  not_user_ids:              [34567, 45678]
  user_names:                ['fekoneko', 'somecoolartist']
  not_user_names:            ['notsocoolartist', 'notme']
  restrictions:              ['none', 'R-18', 'R-18G']
  ai:                        false
  original:                  true
  pages_less_than:           50
  pages_more_than:           3
  views_less_than:           10000
  views_more_than:           1000
  bookmarks_less_than:       1000
  bookmarks_more_than:       100
  likes_less_than:           500
  likes_more_than:           50
  comments_less_than:        10
  comments_more_than:        2
  uploaded_before:           2022-01-01T00:00:00Z00:00
  uploaded_after:            2010-01-01T00:00:00Z00:00
  series:                    true
  series_ids:                [12345, 23456]
  not_series_ids:            [34567, 45678]
  series_title_contains:     ['cute', 'cat']
  series_title_not_contains: ['ugly', 'dog']
  series_title_regexp:       '^.*[0-9]+$'
  tags:                      ['お気に入り', '東方']
  not_tags:                  ['おっぱい', 'AI生成']
`

func RunDownload() {
	if len(os.Args) > 2 && os.Args[2] == "rules" {
		fmt.Print(downloadRulesHelp)
	} else {
		fmt.Print(downloadHelp)
	}
}

// TODO: ability to provide infer id pattern to --collection

// TODO: download user's works ('my' or by id)
// TODO: --save-list option to only save crawl results as a yaml list
// TODO: --log, -L option to log the output to a file (-L should be reserved for language actually)
// TODO: choose language to download work metadata in (flag and config)

// TODO: piximan sort - sort based on metadata

// TODO: test if session ID is valid after configuring by sending a request
// TODO: check if session ID is set in interractive download mode and prompt user if they want to set it
// TODO: detailed descriptions in interactive mode

// TODO: multilpe sources in one command:
//       piximan download --id 12345 --path "./{id}" \
//                        --bookmarks my --path "./{id}"
//       or in list.yaml
