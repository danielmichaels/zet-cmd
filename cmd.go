// Copyright 2022 zet-cmd Authors
// SPDX-License-Identifier: Apache-2.0

package zet

import (
	"text/template"

	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/conf"
	"github.com/rwxrob/help"
	"github.com/rwxrob/vars"
)

// Cmd provides a Bonzai branch command that can be composed into Bonzai
// trees or used as a standalone with light wrapper (see cmd/).
var Cmd = &Z.Cmd{

	Name:      `zet`,
	Summary:   `zettelkasten commander`,
	Version:   `v0.0.1`,
	Copyright: `Copyright 2022 Daniel Michaels`,
	License:   `Apache-2.0`,
	Site:      `danielms.site`,
	Source:    `git@github.com:danielmichaels/zet-cmd.git`,
	Issues:    `github.com/danielmichaels/zet-cmd/issues`,

	Commands: []*Z.Cmd{
		// standard external branch imports (see rwxrob/{help,conf,vars})
		help.Cmd, conf.Cmd, vars.Cmd,

		// local commands (in this module)
		Create, Latest, Edit, Get, Query, Find,
	},

	// Add custom BonzaiMark template extensions (or overwrite existing ones).
	Dynamic: template.FuncMap{
		"uname": func(_ *Z.Cmd) string { return Z.Out("uname", "-a") },
		"dir":   func() string { return Z.Out("dir") },
	},

	Description: `
		The **{{.Name}}** command branch is a well-documented example to get
		you started.  You can start the description here and wrap it to look
		nice and it will just work.  Descriptions are written in BonzaiMark,
		a simplified combination of CommonMark, "go doc", and text/template
		that uses the Cmd itself as a data source and has a rich set of
		builtin template functions ({{pre "pre"}}, {{pre "exename"}}, {{pre
		"indent"}}, etc.). There are four block types and four span types in
		BonzaiMark:

		Spans

		    Plain
		    *Italic*
		    **Bold**
		    ***BoldItalic***
		    <Under> (brackets remain)

		Note that on most terminals italic is rendered as underlining and
		depending on how old the terminal, other formatting might not appear
		as expected. If you know how to set LESS_TERMCAP_* variables they
		will be observed when output is to the terminal.

		Blocks

		1. Paragraph
		2. Verbatim (block begins with '    ', never first)
		3. Numbered (block begins with '* ')
		4. Bulleted (block begins with '1. ')

		Currently, a verbatim block must never be first because of the
		stripping of initial white space.

		Templates

		Anything from Cmd that fulfills the requirement to be included in
		a Go text/template may be used. This includes {{ "{{ .Name }}" }}
		and the rest. A number of builtin template functions have also been
		added (such as {{ "indent" }}) which can receive piped input. You
		can add your own functions (or overwrite existing ones) by adding
		your own Dynamic template.FuncMap (see text/template for more about
		Go templates). Note that verbatim blocks will need to indented to work:

		    {{ "{{ dir | indent 4 }}" }}

		Produces a nice verbatim block:

		{{ dir | indent 4 }}

		Note this is different for every user and their specific system. The
		ability to incorporate dynamic data into any help documentation is
		a game-changer not only for creating very consumable tools, but
		creating intelligent, interactive training and education materials
	 	as well.

		Templates Within Templates

		Sometimes you will need more text than can easily fit within
		a single action. (Actions may not span new lines.) For such things
		defining a template with that text is required and they you can
		include it with the {{pre "template"}} tag.

		    {{define "long" -}}
		    Here is something
		    that spans multiple
		    lines that would otherwise be too long for a single action.
		    {{- end}}

		    The {{ "**{{.Name}}**" }} branch is for everything to help with
		    development, use, and discovery of Bonzai branches and leaf
		    commands ({{ "{{- template \"long\" \"\" | pre -}}" }}).

		The help documentation can scan the state of the system and give
		specific pointers and instruction based on elements of the host
		system that are missing or misconfigured.  Such was *never* possible
		with simple "man" pages and still is not possible with Cobra,
		urfave/cli, or any other commander framework in use today. In fact,
		Bonzai branch commands can be considered portable, dynamic web
		servers (once the planned support for embedded fs assets is
		added).`,

	Other: []Z.Section{
		{`Custom Sections`, `
			Additional sections can be added to the Other field.

			A Z.Section is just a Title and Body and can be assigned using
			composite notation (without the key names) for cleaner, in-code
			documentation.

			The Title will be capitalized for terminal output if using the
			common help.Cmd, but should use a suitable case for appearing in
			a book for other output renderers later (HTML, PDF, etc.)`,
		},
	},
}
