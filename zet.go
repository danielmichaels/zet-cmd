// Copyright 2022 zet-cmd Authors
// SPDX-License-Identifier: Apache-2.0

package zet

import (
	"bufio"
	"errors"
	"fmt"
	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/help"
	"github.com/rwxrob/term"
	"io/ioutil"
	url2 "net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var (
	Pager       = os.Getenv("PAGER")
	Editor      = os.Getenv("EDITOR")
	GitUser     = os.Getenv("GITUSER")
	GitBranch   = "main"
	RepoName    = "zet"
	GitRepo     = filepath.Join(os.Getenv("HOME"), os.Getenv("REPOS"))
	Pictures    = filepath.Join(os.Getenv("HOME"), "Pictures", "zet")
	Screenshots = filepath.Join(os.Getenv("HOME"), "Pictures", "zet")
	Downloads   = filepath.Join(os.Getenv("HOME"), "Downloads")
)

var Create = &Z.Cmd{
	Name:     `create`,
	Aliases:  []string{"new", "c"},
	Summary:  `Create a new zet`,
	Commands: []*Z.Cmd{help.Cmd},
	Usage:    `must provide a title for each zet`,
	MinArgs:  1,
	Call: func(caller *Z.Cmd, args ...string) error {
		z := Zet{Title: args[0]}

		dir, err := z.CreateDir()
		if err != nil {
			return err
		}
		z.Path = dir

		println(dir)
		err = z.CreateReadme(z, dir)
		if err != nil {
			return err
		}

		// Drop into vim and write Zet contents
		err = Z.Exec(Editor, z.GetReadme(z.Path))
		if err != nil {
			return err
		}
		err = z.PullAddCommitPush()
		if err != nil {
			return err
		}
		fmt.Printf("Committing %q\n", z.Title)
		return nil
	},
}
var Get = &Z.Cmd{
	Name:     `get`,
	Aliases:  []string{"g"},
	Summary:  `Retrieve a zet for editing`,
	MinArgs:  1,
	Usage:    `must provide a zet isosec value`,
	Commands: []*Z.Cmd{help.Cmd},
	Call: func(caller *Z.Cmd, args ...string) error {
		z := new(Zet)
		zet, err := z.GetZet(args[0])
		if err != nil {
			return err
		}
		fmt.Printf("%s", zet)

		return nil
	},
}

var Edit = &Z.Cmd{
	Name:     `edit`,
	Aliases:  []string{"e"},
	Summary:  `Edit a zet`,
	MinArgs:  1,
	Usage:    `must provide a zet isosec value`,
	Commands: []*Z.Cmd{help.Cmd},
	Call: func(caller *Z.Cmd, args ...string) error {
		z := new(Zet)
		zet, err := z.GetZet(args[0])
		if err != nil {
			return err
		}
		file := filepath.Join(zet, "README.md")
		err = Z.Exec(Editor, file)
		if err != nil {
			return err
		}

		var r string
		fmt.Printf("Commit? y/N ")
		_, err = fmt.Scanln(&r)
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
	},
}

var Last = &Z.Cmd{
	Name:     `last`,
	Aliases:  []string{"l", "latest"},
	Summary:  `Get the most recent zet`,
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
		z.Latest = last
		fmt.Printf("%s", z.Latest)
		return nil
	},
}
var Query = &Z.Cmd{
	Name:     `query`,
	Aliases:  []string{"q"},
	Summary:  `Create a searchable URL with a query string`,
	MinArgs:  1,
	Usage:    `must provide a search term`,
	Commands: []*Z.Cmd{help.Cmd},
	Call: func(caller *Z.Cmd, args ...string) error {
		term := url2.QueryEscape(args[0])
		url := fmt.Sprintf("https://github.com/%s/%s/search?q=%s", GitUser, RepoName, term)
		fmt.Println(url)
		return nil
	},
}
var Find = &Z.Cmd{
	Name:     `find`,
	Aliases:  []string{"f"},
	Summary:  `Find a zet title by search term`,
	MinArgs:  1,
	Usage:    `must provide a search term`,
	Commands: []*Z.Cmd{help.Cmd},
	Call: func(caller *Z.Cmd, args ...string) error {
		z := new(Zet)

		err := z.ChangeDir(z.GetRepo())
		if err != nil {
			return err
		}
		dir, _ := os.Getwd()
		files, err := z.ReadDir(dir)
		if err != nil {
			return err
		}
		titles, err := z.FindTitles(files)
		if err != nil {
			return err
		}
		results, err := z.SearchTitles(args[0], titles)
		if err != nil {
			return err
		}
		for _, v := range results {
			fmt.Println(v.Id, v.Title)
		}
		return nil
	},
}
var Tags = &Z.Cmd{
	Name:     `tags`,
	Aliases:  []string{"t"},
	Summary:  `Find zet(s) by tag'`,
	MinArgs:  1,
	Usage:    `must provide a search term`,
	Commands: []*Z.Cmd{help.Cmd},
	Call: func(caller *Z.Cmd, args ...string) error {
		z := new(Zet)
		err := z.ChangeDir(z.GetRepo())
		if err != nil {
			return err
		}
		dir, _ := os.Getwd()
		files, err := z.ReadDir(dir)
		if err != nil {
			return err
		}
		results, err := z.FindTags(args[0], files)
		if err != nil {
			return err
		}
		for _, v := range results {
			fmt.Println(v.Id, v.Title)
		}

		return nil
	},
}

// FindTags takes a tag and array of files and then searches the files for the
// tag. Any matches are returned as a slice of Title structs containing the
// Id and Title of the file with the match.
func (z *Zet) FindTags(tag string, files []string) ([]Title, error) {
	var titles []Title
	for _, t := range files {
		z.Path = t
		tag, err := z.SearchTags(tag)
		if err != nil {
			return nil, err
		}
		if !tag {
			continue
		}
		err = z.GetTitle()
		if err != nil {
			return nil, err
		}
		title := Title{
			Id:    t,
			Title: z.Title,
		}
		titles = append(titles, title)
	}
	return titles, nil
}

var Check = &Z.Cmd{
	Name:     `check`,
	Summary:  `Check environment variables and configuration`,
	Commands: []*Z.Cmd{help.Cmd},
	Call: func(caller *Z.Cmd, args ...string) error {
		z := new(Zet)
		err := z.CheckZetConfig()
		if err != nil {
			return err
		}
		return nil
	},
}

// screenshot

// Title holds the id and title for a given Zet when searching the filesystem
type Title struct {
	Id    string
	Title string
}

// FindTitles searches through a slice of files inspecting the title (in Zet
// parlance this is the first line of a Readme) and returns a slice of Title
func (z *Zet) FindTitles(files []string) ([]Title, error) {
	var titles []Title
	for _, t := range files {
		z.Path = t
		err := z.GetTitle()
		if err != nil {
			return nil, err
		}
		title := Title{
			Id:    t,
			Title: z.Title,
		}
		titles = append(titles, title)
	}
	return titles, nil
}

// SearchTitles searches through a slice of Title for any matching query element.
// The search uses strings.Contains and will match partials within a string.
func (z *Zet) SearchTitles(query string, titles []Title) ([]Title, error) {
	var results []Title
	for i, t := range titles {
		if strings.Contains(strings.ToLower(titles[i].Title), query) {
			r := Title{
				Id:    t.Id,
				Title: t.Title,
			}
			results = append(results, r)
		}
	}
	return results, nil
}

// Zet is the struct to hang methods from which are used to create, edit, find
// and delete Zet's.
type Zet struct {
	Title  string
	Path   string
	Latest string
}

// GetRepo returns the GitRepo and RepoName as a filepath.
func (z *Zet) GetRepo() string { return filepath.Join(GitRepo, RepoName) }

// GetReadme returns a filepath with README.md appended using a filepath join. This
// is used to retrieve the full path to the README.md being written to or read from.
func (z *Zet) GetReadme(path string) string { return filepath.Join(path, "README.md") }

// SearchTags scans an open file for a tag and uses regexp to find any matches.
// Matching tags return a truthy boolean.
func (z *Zet) SearchTags(tag string) (bool, error) {
	reg := regexp.MustCompile(fmt.Sprintf("(#%s+)", tag))
	f, _ := os.Open(z.GetReadme(z.Path))
	defer f.Close()
	s := bufio.NewScanner(f)
	hit := false
	for s.Scan() {
		match := reg.MatchString(s.Text())
		if match {
			hit = true
		}
	}
	if err := s.Err(); err != nil {
		return false, err
	}
	return hit, nil
}

// GetTitle inspects the Zet README.md from the z.Path and retrieves the
// h1 title. This ensures that the title is up-to-date as it may have been
// altered after its initial creation.
func (z *Zet) GetTitle() error {
	f, _ := os.Open(z.GetReadme(z.Path))
	defer f.Close()
	s := bufio.NewScanner(f)
	var line int
	for s.Scan() {
		// only get first line and skip everything else
		if line >= 1 {
			break
		}
		t := strings.Replace(s.Text(), "#", "", -1)
		t = strings.TrimSpace(t)
		z.Title = t
		line++
	}
	if err := s.Err(); err != nil {
		return err
	}
	return nil
}

func (z *Zet) GetZet(zet string) (string, error) {
	r := regexp.MustCompile("^[0-9]{14,}$")
	l := regexp.MustCompile("last")
	switch {
	case l.MatchString(zet):
		err := z.ChangeDir(z.GetRepo())
		if err != nil {
			return "", err
		}
		l, err := z.Last()
		if err != nil {
			return "", err
		}
		return l, nil
	case r.MatchString(zet):
		err := z.ChangeDir(z.GetRepo())
		if err != nil {
			return "", err
		}
		return zet, nil
	default:
		return "", errors.New("invalid entry, or zet not found")
	}
}

func (z *Zet) CreateReadme(r Zet, path string) error {
	f := []byte(fmt.Sprintf("# %s\n\n", z.Title))
	err := os.WriteFile(r.GetReadme(path), f, 0755)
	if err != nil {
		return err
	}
	z.Path = path
	return nil
}

// ReadDir reads all files within a given path excluding the .git repo
func (z *Zet) ReadDir(path string) ([]string, error) {
	var files []string
	fileInfo, err := ioutil.ReadDir(path)
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		if file.Name() == ".git" {
			continue
		}
		files = append(files, file.Name())
	}
	return files, nil
}

// CreateDir creates a directory inside the zet repository using the Isosec
// function to create the directory using the returned timestamp.
func (z *Zet) CreateDir() (string, error) {
	path := filepath.Join(GitRepo, RepoName, Isosec())
	err := mkdir(path)
	if err != nil {
		return "", err
	}
	return path, nil
}

// ChangeDir must be called during any git operation otherwise the git command
// cannot be called within the zet repo reliably
func (z *Zet) ChangeDir(path string) error {
	err := os.Chdir(path)
	if err != nil {
		return errors.New(fmt.Sprintf("file does not exist %q", path))
	}
	return nil
}

func (z *Zet) Pull() error {
	// sanity check git remote
	err := z.GitRemote()
	if err != nil {
		return err
	}
	err = z.ChangeDir(z.Path)
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
	err := z.ChangeDir(z.Path)
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
	err := z.ChangeDir(z.Path)
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
	// sanity check git remote
	err := z.GitRemote()
	if err != nil {
		return err
	}
	err = z.ChangeDir(z.Path)
	if err != nil {
		return err
	}
	err = Z.Exec("git", "push", "--quiet")
	if err != nil {
		return err
	}
	return nil
}

// PullAddCommitPush is a helper method which flows through a Git workflow
// and is called often in Commands such as `create` and `edit`.
func (z *Zet) PullAddCommitPush() error {
	err := z.GetTitle()
	if err != nil {
		return err
	}
	err = z.Pull()
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

// Last inspects the Zet repo directories (Isosec folders) and returns the
// most recent directory.
func (z *Zet) Last() (string, error) {
	files, _ := ioutil.ReadDir(z.GetRepo())
	var last string
	var newest int64 = 0
	for _, f := range files {
		if f.Name() == ".git" {
			continue
		}
		fi, err := os.Stat(f.Name())
		if err != nil {
			return "", err
		}
		currTime := fi.ModTime().Unix()
		if currTime > newest {
			newest = currTime
			last = f.Name()
		}
	}
	return last, nil
}

// CheckZetConfig outputs important information about the executable's
// configuration such as environment variables and directory paths. This is
// useful for debugging issues with the host system or failures for the exe
// to commit to GitHub successfully.
func (z *Zet) CheckZetConfig() error {
	fmt.Println(term.U + term.Green + "Checking Zet Config" + term.Reset)
	// System variables
	fmt.Println(term.Blue + "Editor: " + term.Reset + Editor)
	fmt.Println(term.Blue + "Pager: " + term.Reset + Pager)
	// Git/Repo variables
	fmt.Println(term.U + term.Yellow + "Repository Variables" + term.Reset)
	fmt.Println(term.Blue + "RepoName: " + term.Reset + RepoName)
	REPOS := os.Getenv("REPOS")
	if REPOS == "" {
		REPOS = term.Red + "Variable not set. Must point to the `zet` git repo locally." + term.Reset
	}
	fmt.Println(term.Blue + "Repos Variable: " + term.Reset + REPOS)
	fmt.Println(term.Blue + "GitUser: " + term.Reset + GitUser)
	fmt.Println(term.Blue + "GitRepo: " + term.Reset + GitRepo)
	fmt.Println(term.Blue + "System Zet Repo: " + term.Reset + z.GetRepo())
	// Future use case info
	fmt.Println(term.U + term.Yellow + "Utility Directories" + term.Reset)
	fmt.Println(term.Blue + "Pictures Directory: " + term.Reset + Pictures)
	fmt.Println(term.Blue + "Screenshots Directory:" + term.Reset + Screenshots)
	fmt.Println(term.Blue + "Downloads Directory: " + term.Reset + Downloads)
	// Check directories exist
	fmt.Println(term.U + term.Yellow + "Directories Exist" + term.Reset)
	_, err := os.Stat(z.GetRepo())
	zetRepo := "true"
	if err != nil {
		zetRepo = term.Red + "false" + term.Reset
	}
	fmt.Println(term.Blue + "Zet Git Repo Exists: " + term.Reset + fmt.Sprintf("%s", zetRepo))
	if zetRepo == term.Red+"false"+term.Reset {
		fmt.Println(term.Blue + "Zet GitHub Remote: " + term.Red + "Zet repo does not exist on host" + term.Reset)
		return nil
	}
	err = z.ChangeDir(z.GetRepo())
	if err != nil {
		return err
	}
	err = z.GitRemote()
	zetRemote := "true"
	if err != nil {
		zetRemote = term.Red + "false" + term.Reset
	}
	fmt.Println(term.Blue + "Zet GitHub Remote: " + term.Reset + fmt.Sprintf("%s", zetRemote))
	return nil
}

// Isosec returns the GMT current time in ISO8601 (RFC3339) without
// any punctuation or the T.  This is frequently a very good unique
// suffix that has the added advantage of being chronologically sortable
// and more readable than the epoch.
func Isosec() string {
	return fmt.Sprintf("%v", time.Now().In(time.UTC).Format("20060102150405"))
}

// mkdir is the functional equivlent of 'mkdir -p' and is used to create new
// folders recursively.
func mkdir(path string) error {
	err := os.MkdirAll(path, 0664)
	if err != nil {
		return err
	}
	return nil
}
