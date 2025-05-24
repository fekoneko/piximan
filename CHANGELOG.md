# v0.1.0

## piximanctl

### New features

- Download illustrations / manga / ugoira / novels
- Download by ID or from list
- Infer work IDs from existing collection paths
- Download different sizes (resolutions) of the illustrations / manga
- Interactive mode for download and configuration
- Progress bars and verbose logs while downloading
- Store work metadata with downloaded work in YAML format
- Use substitutions in download path: {title} / {id} / {user} / {userid} / {restrict}
- Make requests concurrently when it's possible without bothering the Pixiv servers too much
- Authorize requests with your session ID, piximanctl will try to use it as few as possible
- Encrypt your session ID with a master password

---

# v0.1.1

## piximanctl

### Bug fixes

- Fix downloading empty images if the first guessed extension was incorrect
