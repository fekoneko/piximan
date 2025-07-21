# piximan - Pixiv Manager

Pixiv batch **downloader** and local collection **viewer**. Preserve your favorite art with ease!

## Installation

Go to [Releases](https://github.com/fekoneko/piximan/releases) page

## Viewer features

> [!TODO]
> If I forget to mention something here before the release, my bad

## Downloader features

- Download illustrations / manga / ugoira / novels
- Download different sizes (resolutions) of the illustrations / manga
- Download your bookmarks or bookmarks of another user
- Download by work ID
- Download from list
- Filter downloaded works with download rules
- Avoid downloading duplicate works by providing path to your collection
- Infer work IDs from existing collection paths for easy migration to piximan metadata
- Interactive mode for download and configuration CLI
- Use substitutions in download path: `{title}` / `{id}` / `{user}` / `{user-id}` / etc.
- Make requests concurrently when it's possible without bothering the Pixiv servers too much
- Authorize requests with your session ID, `piximan` will try to use it as few as possible
- Encrypt your session ID with a master password
- Adjust request delays and concurrency limits
- Retry failed requests

## Getting started with `piximan` viewer GUI

To start the application run the command:

```shell
piximan
```

That's it, enjoy!

## Getting started with `piximan` downloader CLI

### Authorization

> If you only download works without age restriction and don't need to fetch user bookmarks,
> the downloader is usable without authorization.

For some requests Pixiv requires you to be authorized. For example, to fetch frames for R-18 ugoira
you must have the R-18 option checked in your profile. To authorize these requests you need to
configure the _session ID_.

You can get session ID from your browser _cookies_ right now:

- Go to [https://www.pixiv.net](https://wwww.pixiv.net)
- On the website press `F12` to access the devtools panel
- In the devtools panel switch to the _Application_ tab (for Chrome) / _Storage_ tab (for Firefox)
- Expand the _Cookies_ section and select `https://www.pixiv.net` origin
- Find the row named `PHPSESSID` - this is your cookie
- Copy the value of the cookie to the clipboard

Now open the terminal and run the command to enter interractive configuration mode:

```shell
piximan config
```

Paste the copied session ID and then specify the master password if you want.

### Downloading with interactive mode

The easiest way to use the tool is with interactive mode. To enter it run the command,
then answer the questions about what to download and where to save the files:

```shell
piximan download
```

### Downloading a work by ID

You're ready to go! Try out `piximan` by downloading an artwork from pixiv:

```shell
piximan download \
  --id 584231 \
  --path './artworks/{user} ({user-id})/{title} ({id})'
```

Downloading a novel is as simple:

```shell
piximan download \
  --id 584231 \
  --type novel \
  --path './novels/{user} ({user-id})/{title} ({id})'
```

### Downloading bookmarks

> For downloading any bookmarks you need to be authorized (configure the session ID).

You can download your public artwork bookmarks like this:

```shell
piximan download \
  --bookmarks my \
  --path './bookmarks/{user} ({user-id})/{title} ({id})'
```

You can also specify `--type novel` to download novel bookmarks and `--private` to download
private bookmarks.

```shell
piximan download \
  --bookmarks my \
  --type novel \
  --private \
  --path './bookmarks/{user} ({user-id})/{title} ({id})'
```

You can also specify a user-assigned tag or download only a specified range. This example
will download your public bookmarks with tag 'お気に入り' from 101th to 200th latest:

```shell
piximan download \
  --bookmarks my \
  --tag 'お気に入り' \
  --from 100 \
  --to 200 \
  --path './bookmarks/{user} ({user-id})/{title} ({id})'
```

You can also download public bookmarks of any user knowing their ID.
For example, this will download novel bookmarks of user 12345:

```shell
piximan download \
  --bookmarks 12345 \
  --type novel \
  --path './bookmarks/{user} ({user-id})/{title} ({id})'
```

### Downloading from list

You can specify a queue for downloader using YAML format such as:

> `./list.yaml`

```yaml
# This will download two artworks with ID 12345 and 23456
- { id: 12345, type: artwork }
- { id: 23456, type: artwork }

# This will override provided downloader arguments
- id: 34567
  type: artwork
  size: 1
  only-meta: false
  paths: ['./special artwork']
- id: 45678
  type: novel
  only-meta: true
  paths: ['./special novel']
```

Start downloading with the command:

```shell
piximan download \
  --list './list.yaml' \
  --path './artworks/{user} ({user-id})/{title} ({id})'
```

### Inferring work IDs from path

You can infer the IDs of works from the given path. For example, this is useful for updating
the metadata in the existing collection when coupled with the `--only-meta` flag:

```shell
piximan download \
  --infer-id './artworks/*/* ({id})' \
  --only-meta
```

### Downloading rules

Rules are used to filter wich works should be downloaded and defined in YAML format.
All rules are optional, and if multiple rules are defined, the work should match all of
them to be downloaded (AND). Any array matches any of its elements (OR).

Here's an example of all available rules:

```yaml
ids: [12345, 23456]
not_ids: [34567, 45678]
title_contains: ['cute', 'cat']
title_not_contains: ['ugly', 'dog']
title_regexp: '^.*[0-9]+$'
kinds: ['illust', 'manga', 'ugoira', 'novel']
description_contains: ['hello', 'world']
description_not_contains: ['goodbye', 'universe']
description_regexp: '^.*[0-9]+$'
user_ids: [12345, 23456]
not_user_ids: [34567, 45678]
user_names: ['fekoneko', 'somecoolartist']
not_user_names: ['notsocoolartist', 'notme']
restrictions: ['none', 'R-18', 'R-18G']
ai: false
original: true
pages_less_than: 50
pages_more_than: 3
views_less_than: 10000
views_more_than: 1000
bookmarks_less_than: 1000
bookmarks_more_than: 100
likes_less_than: 500
likes_more_than: 50
comments_less_than: 10
comments_more_than: 2
uploaded_before: 2022-01-01T00:00:00Z00:00
uploaded_after: 2010-01-01T00:00:00Z00:00
series: true
series_ids: [12345, 23456]
not_series_ids: [34567, 45678]
series_title_contains: ['cute', 'cat']
series_title_not_contains: ['ugly', 'dog']
series_title_regexp: '^.*[0-9]+$'
tags: ['お気に入り', '東方']
not_tags: ['おっぱい', 'AI生成']
```

When downloading, specify the rules with `--rules` flag:

```shell
piximan download --id 12345 --rules './rules.yaml'
```

### Syncing bookmarks with existing collection

You can skip works already present in the collection with `--collection` flag:

```shell
piximan download --bookmarks my --collection '.' --path './{user-id}/{id}'
```

Infer ID pattern can be provided here as well (see `--infer-id` flag). Note tat in this case
all matched work IDs will be assumed to be of type provided with `--type` flag:

```shell
piximan download --bookmarks my --collection './*/{id}' --path './{user-id}/{id}'
```

While the above commands will skip already downloaded works, it will need to fetch
the list of all your bookmarks to ensure there isn't some older one that isn't present
in the collection.

Flag `--fresh` will tell the downloader to stop crawling new bookmark pages once it encounters
a fully ignored one. This may greatly reduce the number of authorized requests to pixiv.net.

So for syncing your new bookmarks once in a while you can use the downloader like this:

```shell
piximan download --bookmarks my --collection '.' --fresh --path './{user-id}/{id}'
```

### Help

To see other options and examples use the `help` command in your terminal:

```shell
piximan help download
piximan help config
```

## Development

Use these Bash scripts to run and build the project:

```shell
./run.sh                    # Run piximan
./run.sh download -i 10000  # Run piximan with the args
./build.sh                  # Build for the current platform
./build.sh $os $arch        # Build for specified OS and architecture
```

## Related projects

This project is the next iteration on my way to perfect local Pixiv app. So far I tried:

- [React Native for Windows](https://github.com/fekoneko/pixiv-powerful-viewer-legacy)
- [Electron + React](https://github.com/fekoneko/pixiv-powerful-viewer/tree/v1.0.0-alpha.2)
- [Tauri + React + Rust](https://github.com/fekoneko/pixiv-powerful-viewer)
- [Rust + GTK](https://github.com/fekoneko/pixiv-powerful-viewer-gtk)
- [Go + GTK](https://github.com/fekoneko/piximan) <- this one
