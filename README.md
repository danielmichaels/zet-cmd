# Zet-cmd

> A command line tool for managing a personal Zettelkasten

This is a CLI tool for managing zet's or notes. It allows me to capture ideas and thoughts from the command line.
GitHub is used as the "database" for note storage. This has the added benefit of being publicly available without
the need for the CLI.

These *zet's* are also available on my [blog](https://danielms.site/zet). I have written a [blog post](https://danielms.site/blog/github-actions-auto-publish-zettelkasten-notes/)
about how I use GitHub Actions to automatically publish these notes.

[![License](https://img.shields.io/badge/license-Apache2-brightgreen.svg)](LICENSE)

## Install

```bash
go install github.com/danielmichaels/zet-cmd/cmd/zet@latest
```
## Requirements

`zet-cmd` must have a GitHub repository to push commits to named `zet`. For instance, my 
personal `zet` repository is [github.com/danielmichaels/zet][ghzet]. Without this `zet` will not 
have a remote repository to commit to.

On a new machine (but existing `zet` repo), you will need to `git clone` to the new device first.

### Environment Variables

- `EDITOR` must be set to create and edit Zet's.
- `GITUSER` must be your GitHub account username
- `ZETDIR` should point to the `zet` repo on your system e.g. `$HOME/Code/github/zet`. Without this `zet` cannot find the directory or files

**📣 Note**

`zet-cmd` has a `check` command which will output the required environment variables and directory
paths. Any `false` values or empty `Repo` entries will need to be rectified or your `zet-cmd` may
not function as expected, or at all.

### Inspiration

A lot of this was inspired by [rwxrob](https://github.com/rwxrob/)

[ghzet]: https://github.com/danielmichaels/zet
