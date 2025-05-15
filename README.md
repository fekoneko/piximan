# piximan - Pixiv Manager

Pixiv batch **downloader** and local collection **viewer**. Preserve your favorite art with ease!

> [!NOTE]
> The GUI for viewer is yet to be implemented. By now you can use CLI tool to download works.

## Downloader Features

- Download illustrations / manga / ugoira / novels
- Download by ID or from list
- Infer work IDs from existing collection paths
- Download different sizes (resolutions) of the illustrations / manga
- Interractive mode for download and configuration with `piximanctl` tool
- Store work metadata with downloaded work in _YAML_ format
- Use substitutions in download path: `{title}` / `{id}` / `{user}` / `{userid}` / `{restrict}`
- Make requests concurrently when it's possible without bothering the Pixiv servers too much
- Authorize requests with your session ID, `piximanctl` will try to use it as few as possible
- Encrypt your session ID with a master password

## Getting started with `piximanctl` CLI tool

### Authorization

> If you only download works without restriction (without R-18, R-18G) you can skip this section.

For some requests Pixiv requires you to be authorized. For example, to fetch frames for R-18 ugoira
you must have the R-18 option checked in your profile. To authorize these requests you need to
configure the _session ID_.

You can get session ID from your browser _cookies_ right now:

- Go to [https://www.pixiv.net](https://wwww.pixiv.net)
- On the website press `F12` to access the devtools panel
- In the devtools panel switch to the _Application_ tab (for Chrome) / _Storage_ tab (for Firefox)
- Expand the _Cookies_ section and select `https://wwww.pixiv.net` origin
- Find the row named `PHPSESSID` - this is your cookie
- Copy the value of the cookie to the clipboard

Now open the terminal and run the command to enter interractive configuration mode:

```shell
piximanctl config
```

Paste the copied session ID and then specify the master password if you want.

### Downloading with interactive mode

The easiest way to use the tool is with interactive mode. To enter it run the command,
then answer the questions about what to download and where to save the files:

```shell
piximanctl download
```

### Downloading a work by ID

You're ready to go! Try out `piximanctl` by downloading an artwork from pixiv:

```shell
piximanctl download \
  --id 584231 \
  --path './artworks/{user} ({userid})/{title} ({id})'
```

Downloading a novel is as simple:

```shell
piximanctl download \
  --id 584231 \
  --type novel \
  --path './novels/{user} ({userid})/{title} ({id})'
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
  onlymeta: false
  paths: ['./special artwork']
- id: 45678
  type: novel
  onlymeta: true
  paths: ['./special novel']
```

Start downloading with the command:

```shell
piximanctl download \
  --list './list.yaml' \
  --path './artworks/{user} ({userid})/{title} ({id})'
```

### Inferring work IDs

You can infer the IDs of works from the given path. For example, this is useful for updating
the metadata in the existing collection when coupled with the `--onlymeta` flag:

```shell
piximanctl download \
  --inferid './artworks/*/* ({id})' \
  --onlymeta
```

### Help

To see other options and examples use the `help` command in your terminal:

```shell
piximanctl help download
piximanctl help config
```

## Development

Use `make` to run and build the project:

```shell
make run:piximan                    # Run piximan GUI
make run:piximanctl ARGS="download" # Run piximanctl CLI tool with the arguments
make build:piximan                  # Build piximan GUI
make build:piximanctl               # Build piximanctl CLI tool
make build                          # Build both
```

## Related projects

This project is the next iteration on my way to perfect local Pixiv app. So far I tried:

- [React Native for Windows](https://github.com/fekoneko/pixiv-powerful-viewer-legacy)
- [Electron + React](https://github.com/fekoneko/pixiv-powerful-viewer/tree/v1.0.0-alpha.2)
- [Tauri + React + Rust](https://github.com/fekoneko/pixiv-powerful-viewer)
- [Rust + GTK](https://github.com/fekoneko/pixiv-powerful-viewer-gtk)
- [Go + GTK](https://github.com/fekoneko/piximan) <- this one
