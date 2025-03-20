# Jq completion
![GitHub Release](https://img.shields.io/github/v/release/matsuren/jqcompletion)

Jq key completion tui app. This is a hobby project to learn Go by creating some tools.

### Demo
![demo](.README/demo.gif)

## TODO
- [ ] Handle `|` and [0] in query
- [x] print jq query
- [x] open output view in vim (etc.)
- [x] Handle large file -> better than tview version due to async update
- [ ] Better to use fzf
- [ ] Accept stdin
- [ ] Add jq keyword completion, e.g., length, select, etc.

## Installation

Download binary
```
# Get the latest release version
LATEST=$(curl -s https://api.github.com/repos/matsuren/jqcompletion/releases/latest | grep tag_name | cut -d '"' -f 4)

# For linux amd64 architecture
curl -L "https://github.com/matsuren/jqcompletion/releases/download/${LATEST}/jqcompletion_${LATEST#v}_linux_amd64.tar.gz" | tar -xz

# For macOS arm64 architecture
curl -L "https://github.com/matsuren/jqcompletion/releases/download/${LATEST}/jqcompletion_${LATEST#v}_darwin_arm64.tar.gz" | tar -xz

# Move to a directory in your PATH
mv jqcompletion $HOME/.local/bin
```

Or
```
go install github.com/matsuren/jqcompletion
```

## Usage

```
jqcompletion sample.json
```

UI Keybindings
- Ctrl+p/n: Suggestions up/down
- Tab: Sync selected suggestion to Query
- Enter: Evaluate query

## Inspiration

I got inspirations from the followings:

- `echo "" | fzf --print-query --preview 'jq {q} sample.json'`
- https://github.com/ynqa/jnv

## Development
Install mise and run the following to setup the environment
```
mise install
```
