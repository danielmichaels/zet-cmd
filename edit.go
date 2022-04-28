package zet

import (
	"fmt"
	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/help"
	"path/filepath"
	"strings"
)

var EditCmd = &Z.Cmd{
	Name:     `edit`,
	Aliases:  []string{"e"},
	Summary:  `Edit a zet`,
	MinArgs:  1,
	Usage:    `must provide a zet isosec value`,
	Commands: []*Z.Cmd{help.Cmd, editLast},
	Call: func(caller *Z.Cmd, args ...string) error {
		z := new(Zet)

		zet, err := z.GetZet(args[0])
		if err != nil {
			return err
		}
		err = z.editZet(zet)
		if err != nil {
			return err
		}

		err = z.scanAndCommit(zet)
		if err != nil {
			return err
		}
		return nil
	},
}

func (z *Zet) editZet(zet string) error {
	file := filepath.Join(zet, "README.md")
	err := Z.Exec(Editor, file)
	if err != nil {
		return err
	}
	return nil
}

var editLast = &Z.Cmd{
	Name:     `last`,
	Summary:  `edit the last modified zet entry from the git repo`,
	Commands: []*Z.Cmd{help.Cmd},
	Call: func(caller *Z.Cmd, args ...string) error {
		z := new(Zet)
		err := z.ChangeDir(z.GetRepo())
		if err != nil {
			return err
		}

		last, err := z.Last()
		if err != nil {
			return err
		}

		err = z.editZet(last)
		if err != nil {
			return err
		}

		err = z.scanAndCommit(last)
		if err != nil {
			return err
		}

		return nil
	},
}

func (z *Zet) scanAndCommit(zet string) error {
	var r string
	fmt.Printf("Commit? y/N ")
	_, err := fmt.Scanln(&r)
	if err != nil {
		// <Enter> will return the following error, so we mark it as "N".
		if err.Error() != "unexpected newline" {
			return err
		}
		r = "N"
	}
	r = strings.TrimSpace(r)
	r = strings.ToLower(r)

	if r != "y" {
		fmt.Printf("%q not commited but modified\n", zet)
		return nil
	}
	z.Path = filepath.Join(z.GetRepo(), zet)
	err = z.PullAddCommitPush()
	if err != nil {
		return err
	}
	return nil
}
