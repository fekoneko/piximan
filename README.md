# piximan - Pixiv Manager

Pixiv batch **downloader** and local collection **viewer**. Preserve your favorite art with ease!

> [!NOTE]
> The project is still in early development

## Included tools

- `piximan` - collection viewer GUI - **not yet implemented**
- `piximanctl` - CLI download tool for your automated scripts

## Features

### Downloader
- Illustrations / manga / ugoira / novels are all supported for download
- You can download different sizes (resolution) of the illustrations / manga
- Useful metadata is saved with the downloaded work in _YAML_ format
- Supperted substitutions in destination path: `{title}` / `{id}` / `{user}` / `{userid}`

## Getting started with `piximanctl` CLI tool

Before using the tool you need to configure the _session ID_. This will authorize you on _pixiv.net_ and
let `piximanctl` fetch work metadata.

You can get session ID with your browser _cookies_ right now:
- Go to [https://www.pixiv.net](https://wwww.pixiv.net)
- On the website press `F12` to access the devtools panel
- In the devtools panel switch to the _Application_ tab (for Chrome) / _Storage_ tab (for Firefox)
- Expand the _Cookies_ section and select `https://wwww.pixiv.net` origin
- Find the row named `PHPSESSID` - this is your cookie
- Copy the value of the cookie to the clipboard

Almost there! Let's configure `piximanctl` to permanently use the copied session ID.
To pass the cookie from the clipboard use one of the following commands:
- On Linux (X11):
```shell
piximanctl config -sessionid $(xclip -o)
```
- On Linux (Wayland):
```shell
piximanctl config -sessionid $(wl-paste)
```
- On Windows:
```powershell
piximanctl config -sessionid $(Get-Clipboard)`
```
- On MacOS:
```shell
piximanctl config -sessionid $(pbpaste)
```

You can always paste the value directly to your terminal, but be sure that it is not logged to the
terminal history (e.g. `~/.bash_history` very commonly on Linux).

You're ready to go! Try out `piximanctl` by downloading some artwork from pixiv:

```shell
piximanctl download -id 584231
```

You can also specify the destination path with the flag `-path`.
With flag `-type novel` you can also download nolels.

There are some extra options that can be used.
The manual is always available to you, just run the program without the arguments:

```shell
piximanctl
```

## Related projects

This project is the next iteration on my way to perfect local Pixiv app. So far I tried:
- [React Native for Windows](https://github.com/fekoneko/pixiv-powerful-viewer-legacy)
- [Electron + React](https://github.com/fekoneko/pixiv-powerful-viewer/tree/v1.0.0-alpha.2)
- [Tauri + React + Rust](https://github.com/fekoneko/pixiv-powerful-viewer)
- [Rust + GTK](https://github.com/fekoneko/pixiv-powerful-viewer-gtk)
- [Go + GTK](https://github.com/fekoneko/piximan) <- this one
