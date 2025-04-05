package zet

import (
	"errors"
	"fmt"
	"github.com/danielmichaels/zet-cmd/internal/term"
	"os"
)

type GitCmd struct {
	Operation string `arg:"" help:"git command to run" enum:"pull,status,push,log,stash,stash-pop,lazygit,lg"`
}

func (c *GitCmd) Run() error {
	z := new(Zet)
	err := z.ChangeDir(z.GetRepo())
	if err != nil {
		return err
	}

	var cmdArgs []string
	switch c.Operation {
	case "pull":
		cmdArgs = []string{"git", "pull"}
	case "status":
		cmdArgs = []string{"git", "status"}
	case "push":
		cmdArgs = []string{"git", "push"}
	case "log":
		cmdArgs = []string{"git", "log"}
	case "stash":
		cmdArgs = []string{"git", "stash"}
	case "stash-pop":
		cmdArgs = []string{"git", "stash", "pop"}
	case "lazygit", "lg":
		cmdArgs = []string{"lazygit"}
	default:
		return errors.New("invalid git command")
	}
	if err := term.Exec(cmdArgs...); err != nil {
		return err
	}
	return nil
}

// scanAndCommit checks that the user wants to commit their work to the VCS
// and pushes the commit if they accept.
func (z *Zet) scanAndCommit(zet string) error {
	if term.Prompt("Commit? (y/N) ") != "y" {
		fmt.Printf("%q not commited but modified\n", zet)
		os.Exit(0)
	}
	err := z.PullAddCommitPush()
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
			return errors.New("failed to ascertain zet title")
		}
	}
	err := z.Pull()
	if err != nil {
		return fmt.Errorf("failed to pull from git remote: %w", err)
	}
	err = z.Add()
	if err != nil {
		return fmt.Errorf("failed to add files to git: %w", err)
	}
	err = z.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit files to git: %w", err)
	}
	err = z.Push()
	if err != nil {
		return fmt.Errorf("failed to push files to git: %w", err)
	}
	return nil
}

// GitRemote checks for the existence of a non-empty `git remote -v` response.
func (z *Zet) GitRemote() error {
	if os.Getenv("GIT_REMOTE") != "" {
		return nil
	}
	s := term.Out("git", "remote", "-v")
	if s == "" {
		return errors.New("no git remote found")
	}
	return nil
}

// Pull changes the current directory to the repository, verifies the git remote, and pulls the latest changes
// from the remote repository with quiet output.
func (z *Zet) Pull() error {
	err := z.ChangeDir(Repo)
	if err != nil {
		return err
	}
	err = z.GitRemote()
	if err != nil {
		return err
	}
	err = term.Exec("git", "pull", "-q")
	if err != nil {
		return err
	}
	return nil
}

// Add stages all files in the Zet's path to the git repository by changing to the repository directory
// and executing a git add command with the all (-A) flag.
func (z *Zet) Add() error {
	err := z.ChangeDir(Repo)
	if err != nil {
		return err
	}
	err = term.Exec("git", "add", "-A", z.Path)
	if err != nil {
		return err
	}
	return nil
}

// Commit stages and commits the current changes in the Zet repository with the Zet's title as the commit message.
// It changes the current directory to the repository, executes a git commit command, and prints a confirmation message.
func (z *Zet) Commit() error {
	err := z.ChangeDir(Repo)
	if err != nil {
		return err
	}
	err = term.Exec("git", "commit", "-m", z.Title)
	if err != nil {
		return err
	}
	fmt.Printf("Committed %q\n", z.Title)
	return nil
}

// Push changes the current directory to the repository, verifies the git remote, and pushes the current branch
// to the remote repository with quiet output.
func (z *Zet) Push() error {
	err := z.ChangeDir(Repo)
	if err != nil {
		return err
	}
	err = z.GitRemote()
	if err != nil {
		return err
	}
	err = term.Exec("git", "push", "--quiet")
	if err != nil {
		return err
	}
	return nil
}
