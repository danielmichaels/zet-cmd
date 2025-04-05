package main

import (
	"fmt"
	"github.com/danielmichaels/zet-cmd"
	"github.com/danielmichaels/zet-cmd/internal/version"
	"os"

	"github.com/alecthomas/kong"
)

const appName = "zet"

type VersionFlag string

func (v VersionFlag) Decode(_ *kong.DecodeContext) error { return nil }
func (v VersionFlag) IsBool() bool                       { return true }
func (v VersionFlag) BeforeApply(app *kong.Kong, vars kong.Vars) error {
	fmt.Println(vars["version"])
	app.Exit(0)
	return nil
}

type CLI struct {
	Version VersionFlag `help:"Print version information and quit" name:"version"`
	zet.Globals

	// Commands
	Create zet.CreateCmd `cmd:"" aliases:"new" help:"Create a new zet"`
	Last   zet.LastCmd   `cmd:"" help:"Show the last created zet's isosec (location)'"`
	Edit   zet.EditCmd   `cmd:"" help:"Edit a zet"`
	Find   zet.FindCmd   `cmd:"" help:"Search for a zet title and retrieve any matching entry"`
	Check  zet.CheckCmd  `cmd:"" help:"Check zettelkasten for issues"`
	Tags   zet.TagsCmd   `cmd:"" help:"Search for a zet by tag and retrieve any entries with that tag"`
	Git    zet.GitCmd    `cmd:"" help:"Git operations for zettelkasten"`
	View   zet.ViewCmd   `cmd:"" help:"View supports both direct 'isosec' lookup's and keyword searches"`
}

func run() error {
	ver := version.Get()
	if ver == "unavailable" {
		ver = "development"
	}
	cli := CLI{
		Version: VersionFlag(ver),
	}
	// Display help if no args are provided instead of an error message
	if len(os.Args) < 2 {
		os.Args = append(os.Args, "--help")
	}

	ctx := kong.Parse(&cli,
		kong.Name(appName),
		kong.Description(fmt.Sprintf("%s is a zettelkasten tool", appName)),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
		kong.DefaultEnvars(appName),
		kong.Vars{
			"version": string(cli.Version),
		})
	err := ctx.Run(cli.Globals)
	ctx.FatalIfErrorf(err)
	return nil
}
func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}
