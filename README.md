# Dango 🍡

store file path to files in a shared text based clipboard
copy the path and create sym or mv or delete

## Install structions for now (might make an installer)


- mkdir ~/.config/dango
- touch ~/.config/dango/dango.toml
- git cloen this repo
- git mod tidy
- go build ./main.go
- add alias dango="path/to/dango/main" to your terminal rc


```
$ dango pickup

What file do you want to pickup?
Press space to add.
Press c to copy.
Press q to close.

> [✨] /Users/al/0_projects/stagfoo/dango/.gitignore
  [✨] /Users/al/0_projects/stagfoo/chuchu/README.md
  [✨] /Users/al/0_projects/stagfoo/dango/README.md
  [🍡] /Users/al/0_projects/stagfoo/dango/go.mod
  [✨] /Users/al/1_daily/notes.md

```

space is add (🍡) and remove (✨)

## Supported Commands

- list: lists the current paths and allows you to copy and remove
- pickup: lists the files in the current dir and lets you add and remove items


## Issues
- currently uses pbclip instead of piping
