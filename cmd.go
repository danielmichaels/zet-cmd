// Copyright 2022-2024 zet-cmd Authors
// SPDX-License-Identifier: Apache-2.0

package zet

import (
	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/conf"
	"github.com/rwxrob/help"
	"github.com/rwxrob/vars"
)

var Cmd = &Z.Cmd{

	Name:      `zet`,
	Summary:   `zettelkasten commander`,
	Version:   `v0.5.0`,
	Copyright: `Copyright 2022-2024 Daniel Michaels`,
	License:   `Apache-2.0`,
	Site:      `danielms.site`,
	Source:    `git@github.com:danielmichaels/zet-cmd.git`,
	Issues:    `github.com/danielmichaels/zet-cmd/issues`,

	Commands: []*Z.Cmd{
		// standard external branch imports (see rwxrob/{help,conf,vars})
		help.Cmd, conf.Cmd, vars.Cmd,

		// local commands (in this module)
		CreateCmd, LastCmd, EditCmd, GetCmd, QueryCmd,
		FindCmd, CheckCmd, TagsCmd, GitCmd, ViewCmd,
	},
	Description: `
		The **{{.Name}}** command is Zettelkasten Bonzai branch used to create
		small slips of knowledge. Those slips are then uploaded to Github for
		public search-ability and ease of use.
		`,
}
