# [WIP] Jq completion

[WIP] Jq completion tool.
Hobby project to learn Go by creating some tools.

## TODO
- [ ] Better to use fzf
- [ ] Accept stdin
- [ ] Handle `|` in query
- [ ] Add jq keyword completion, e.g., length, select, etc.

## Installation
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
