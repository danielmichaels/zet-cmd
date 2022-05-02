package zet

import (
	"errors"
	"fmt"
	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/help"
	"os"
	"strings"
)

var GitCmd = &Z.Cmd{
	Name:     `git`,
	Summary:  `run git commands over the zet repo`,
	MinArgs:  1,
	Usage:    `must provide a git command`,
	Commands: []*Z.Cmd{help.Cmd, gitPull},
}

var gitPull = &Z.Cmd{
	Name:     `pull`,
	Summary:  `retrieve upstream latest commits`,
	Commands: []*Z.Cmd{help.Cmd},
	Call: func(caller *Z.Cmd, args ...string) error {
		z := new(Zet)
		err := z.ChangeDir(z.GetRepo())
		if err != nil {
			return err
		}
		err = Z.Exec("git", "pull")
		if err != nil {
			return err
		}
		return nil
	},
}

// scanAndCommit checks that the user wants to commit their work to the VCS
// and pushes the commit if they accept.
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
	err = z.PullAddCommitPush()
	if err != nil {
		return err
	}
	return nil
}

// PullAddCommitPush is a helper method which flows through a Git workflow
// and is called often in Commands such as `create` and `edit`.
func (z *Zet) PullAddCommitPush() error {
	if z.Title == "" {
		err := z.GetTitle()
		if err != nil {
			return err
		}
	}
	err := z.Pull()
	if err != nil {
		return err
	}
	err = z.Add()
	if err != nil {
		return err
	}
	err = z.Commit()
	if err != nil {
		return err
	}
	err = z.Push()
	if err != nil {
		return err
	}
	return nil
}

// GitRemote checks for the existence of a non-empty `git remote -v` response.
func (z *Zet) GitRemote() error {
	if os.Getenv("GIT_REMOTE") != "" {
		return nil
	}
	s := Z.Out("git", "remote", "-v")
	if s == "" {
		return errors.New("no git remote found")
	}
	return nil
}

func (z *Zet) Pull() error {
	err := z.GitRemote()
	if err != nil {
		return err
	}
	err = z.ChangeDir(ZetRepo)
	if err != nil {
		return err
	}
	err = Z.Exec("git", "pull", "-q")
	if err != nil {
		return err
	}
	return nil
}

func (z *Zet) Add() error {
	err := z.ChangeDir(ZetRepo)
	if err != nil {
		return err
	}
	err = Z.Exec("git", "add", "-A", z.Path)
	if err != nil {
		return err
	}
	return nil
}

func (z *Zet) Commit() error {
	err := z.ChangeDir(ZetRepo)
	if err != nil {
		return err
	}
	err = Z.Exec("git", "commit", "-m", z.Title)
	if err != nil {
		return err
	}
	fmt.Printf("Committed %q\n", z.Title)
	return nil
}

func (z *Zet) Push() error {
	err := z.GitRemote()
	if err != nil {
		return err
	}
	err = z.ChangeDir(ZetRepo)
	if err != nil {
		return err
	}
	err = Z.Exec("git", "push", "--quiet")
	if err != nil {
		return err
	}
	return nil
}
