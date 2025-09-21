# tweet-me

**Generate a tweet-styled image from terminal**

Because screenshots and design tools are boring. **tweet-me** gives you clean, beautiful tweet images you can drop anywhere, blogs, slides, social posts, or just to share your thoughts. 

![Screenshot](./tweet.png)

### Features
- Type out a tweet in your terminal
- Export beautiful, high quality images
- Add your name, handle, and profile picture
- Auto-generate interaction counts for authentic look
- Switch between light and dark mode

### Quick Install

Option 1 (Go toolchain):
```
go install github.com/indiecodermm/tweet-me@latest
```
(Ensure `$(go env GOPATH)/bin` is on your PATH.)

Option 2 (Make):
```
git clone https://github.com/indiecodermm/tweet-me.git
cd tweet-me
make install
```

Option 3 (Script):
```
git clone https://github.com/indiecodermm/tweet-me.git
cd tweet-me
./install.sh    # add --system for /usr/local/bin
```

### Usage
```
tweet-me             # interactive prompts
tweet-me config      # show current config
tweet-me config -user "Alice" -handle "alice" -dark true -outdir ./out
tweet-me --version
```

### Config File
Stored at:
```
~/.config/tweet-me/config.json
```
Example:
```json
{ "user": "Alice", "handle": "alice", "output_directory": "/home/alice/pics", "dark": true }
```

Place an avatar at:
```
~/.config/tweet-me/profile.png
```

### Building Manually
```
go build -trimpath -ldflags "-s -w -X main.version=$(git describe --tags --always --dirty 2>/dev/null || echo dev)" -o tweet-me .
```

### Releasing (multi-platform artifacts)
```
make release
ls dist/
```

### Dependencies
* Go 1.21+ (module)
* wkhtmltoimage (PNG generation). If absent, HTML is still produced.

### Dark Mode
Set once:
```
tweet-me config -dark true
```
Generates with `.dark` theme variables applied.

### Uninstall
```
rm -f ~/.local/bin/tweet-me
```

### License

This project is licensed under [GLWTPL](./LICENSE). 

*Good Luck!*
