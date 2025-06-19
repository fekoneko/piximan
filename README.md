# piximan - Pixiv Manager

Pixiv batch **downloader** and local collection **viewer**. Preserve your favorite art with ease!

> [!NOTE]
> The GUI for viewer is yet to be implemented. By now you can use CLI tool to download works.

## Installation

Go to [Releases](https://github.com/fekoneko/piximan/releases) page

## Downloader Features

- Download illustrations / manga / ugoira / novels
- Download user bookmarks by ID
- Download by ID or from list
- Infer work IDs from existing collection paths
- Download different sizes (resolutions) of the illustrations / manga
- Interactive mode for download and configuration with `piximan` CLI
- Store work metadata with downloaded work in _YAML_ format
- Use substitutions in download path: `{title}` / `{id}` / `{user}` / `{user-id}` / etc.
- Make requests concurrently when it's possible without bothering the Pixiv servers too much
- Authorize requests with your session ID, `piximan` will try to use it as few as possible
- Encrypt your session ID with a master password
- Adjust request delays and concurrency limits

## Getting started with `piximan` CLI

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

> For downloading any bookmarks you need to be authorized (configuration the session ID).

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

You can also specify a user-assigned tag or download only a specified chunk. This example
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

### Inferring work IDs

You can infer the IDs of works from the given path. For example, this is useful for updating
the metadata in the existing collection when coupled with the `--only-meta` flag:

```shell
piximan download \
  --infer-id './artworks/*/* ({id})' \
  --only-meta
```

### Help

To see other options and examples use the `help` command in your terminal:

```shell
piximan help download
piximan help config
```

## Development

Use `make` to run and build the project:

```shell
make run                 # Run piximan GUI
make run ARGS='download' # Run piximan CLI tool with the arguments
make build               # Build both for all platforms
make build               # Build for all platforms
make build:current       # Build for current platform
make build:$PLATFORM     # Build for $PLATFORM
```

## Related projects

This project is the next iteration on my way to perfect local Pixiv app. So far I tried:

- [React Native for Windows](https://github.com/fekoneko/pixiv-powerful-viewer-legacy)
- [Electron + React](https://github.com/fekoneko/pixiv-powerful-viewer/tree/v1.0.0-alpha.2)
- [Tauri + React + Rust](https://github.com/fekoneko/pixiv-powerful-viewer)
- [Rust + GTK](https://github.com/fekoneko/pixiv-powerful-viewer-gtk)
- [Go + GTK](https://github.com/fekoneko/piximan) <- this one
