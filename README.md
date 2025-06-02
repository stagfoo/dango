# Dango 🍡 (WIP)

store file path to files in a shared text based clipboard
copy the path and create sym or mv or delete

## Goal

I wanted to make a little util that allowed me to get just pick the files I wanted from a folder in the terminal and pipe those paths into other clis like mv or copy or clipboard
things like keys or templates or envs between miroservices

## Install instructions for now (might make an installer)

- git clone this repo
- git mod tidy
- go build ./main.go
- add alias dango="path/to/dango/main" to your terminal rc


```
$ dango stick

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
stick also accepts piping

space is add (🍡) and remove (✨)

## Supported Commands

- stuck: lists the current paths saved in your config and allows you to copy and remove
- stick: accepts a list of the files from a pipe and lets you save the path to your config
- drop: lists of the files with the full path in the current dir which can be piped to pickup
- serve: echos current list (good for greping)


## Issues
- currently uses `pbclip` instead of piping
- No tests

## Example usage

migrating a pull request template

```
$ dango stick
$ > [✨] /Users/al/0_projects/basal/project_a/.github/PULL_REQUEST_TEMPLATE.md
$ cd new_project

$ cp "$(dango serve | grep 'PULL')" ./.github/PULL_REQUEST_TEMPLATE.md
```
I can copy the new template whenever its missing to the new repo

```
$ dango stick
$ > [✨] /Users/al/0_projects/basal/project_a/alias.sh
$ dango stuck
$ > [✨] /Users/al/0_projects/basal/project_a/alias.sh
$ path copied
$ nvim /Users/al/0_projects/basal/.dotfiles/alias.sh
```
I can copy the new template whenever its missing to the new repo
