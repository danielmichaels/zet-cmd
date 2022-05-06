# Zet-cmd

[![License](https://img.shields.io/badge/license-Apache2-brightgreen.svg)](LICENSE)

## Install

This command can be installed as a standalone program or composed into a
Bonzai command tree.

Standalone

```
go install github.com/danielmichaels/zet-cmd/cmd/zet@latest
```
Composed

```go
package zet

import (
	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/help"
	"github.com/danielmichaels/zet-cmd"
)

var Cmd = &Z.Cmd{
	Name:     `zet`,
	Commands: []*Z.Cmd{help.Cmd, zet.Cmd},
}
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

## Tab Completion

To activate bash completion just use the `complete -C` option from your
`.bashrc` or command line. There is no messy sourcing required. All the
completion is done by the program itself.

```
complete -C zet zet
```

If you don't have bash or tab completion check use the shortcut
commands instead.

## Embedded Documentation

All documentation (like manual pages) has been embedded into the source
code of the application. See the source or run the program with help to
access it.

## Other Examples

* <https://github.com/rwxrob/cmd-zet> - with heavy inspiration
* <https://github.com/danielmichaels/ds> - my personal Bonzai commander
 
[ghzet]: https://github.com/danielmichaels/zet
