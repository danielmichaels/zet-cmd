// Copyright 2022 zet-cmd Authors
// SPDX-License-Identifier: Apache-2.0

package zet

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/charmbracelet/glamour"
	"github.com/danielmichaels/zet-cmd/internal/term"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	zetRegex    = "^[0-9]{14,}$"
	zetWordWrap = 73
)

var (
	Pager       = os.Getenv("PAGER")
	Editor      = os.Getenv("EDITOR")
	GitUser     = os.Getenv("GITUSER")
	RepoName    = "zet"
	Repo        = filepath.Join(os.Getenv("ZETDIR"))
	Pictures    = filepath.Join(os.Getenv("HOME"), "Pictures", "zet")
	Screenshots = filepath.Join(os.Getenv("HOME"), "Pictures", "zet")
	Downloads   = filepath.Join(os.Getenv("HOME"), "Downloads")
)

// Zet is the struct to hang methods from which are used to create, edit, find
// and delete Zet's.
type Zet struct {
	Title  string
	Path   string
	Latest string
}

func (z *Zet) render(arg string) error {
	err := z.searchScanner(arg)
	if err != nil {
		return err
	}
	p := z.GetReadme(filepath.Join(Repo, z.Path))
	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(), glamour.WithWordWrap(zetWordWrap),
	)
	if err != nil {
		return err
	}
	c, err := os.ReadFile(p)
	if err != nil {
		return err
	}

	out, err := r.Render(string(c))
	if err != nil {
		return err
	}
	fmt.Print(out)
	return nil
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

// GetRepo returns the GitRepo and RepoName as a filepath.
func (z *Zet) GetRepo() string { return Repo }

// GetReadme returns a filepath with README.md appended using a filepath join. This
// is used to retrieve the full path to the README.md being written to or read from.
func (z *Zet) GetReadme(path string) string { return filepath.Join(path, "README.md") }

// SearchTags scans an open file for a tag and uses regexp to find any matches.
// Matching tags return a truthy boolean.
func (z *Zet) SearchTags(tag string) (bool, error) {
	reg := regexp.MustCompile(fmt.Sprintf("(#%s+)", tag))
	f, _ := os.Open(z.GetReadme(z.Path))
	defer func(f *os.File) {
		_ = f.Close()
	}(f)
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
	defer func(f *os.File) {
		_ = f.Close()
	}(f)
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

// CreateReadme builds the zet README.md file structure, sets the permissions
// and the path to the file on a Zet struct.
func (z *Zet) CreateReadme(r Zet, path string) error {
	f := []byte(fmt.Sprintf("# %s\n\n", z.Title))
	err := os.WriteFile(r.GetReadme(path), f, 0664)
	if err != nil {
		return err
	}
	z.Path = path
	return nil
}

// ReadDir reads all files within a given path excluding the .git repo
func (z *Zet) ReadDir(path string) ([]string, error) {
	r := regexp.MustCompile(zetRegex)
	var files []string
	fileInfo, err := os.ReadDir(path)
	if err != nil {
		return files, err
	}

	for _, f := range fileInfo {
		if !r.MatchString(f.Name()) {
			continue
		}
		files = append(files, f.Name())
	}
	return files, nil
}

// CreateDir creates a directory inside the zet repository using the Isosec
// function to create the directory using the returned timestamp.
func (z *Zet) CreateDir() (string, error) {
	path := filepath.Join(Repo, Isosec())
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
		return fmt.Errorf("file does not exist %q", path)
	}
	return nil
}

// Last inspects the Zet repo directories (Isosec folders) and returns the
// most recent directory.
func (z *Zet) Last() (string, error) {
	files, _ := os.ReadDir(z.GetRepo())
	r := regexp.MustCompile(zetRegex)
	var last string
	var newest int64 = 0
	for _, f := range files {
		if !r.MatchString(f.Name()) {
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
	// Set path on Zet now as its used everywhere
	z.Path = filepath.Join(Repo, last)
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
	fmt.Println(term.Blue + "Repo: " + term.Reset + Repo)
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
	fmt.Println(term.Blue + "Zet Git Repo Exists: " + term.Reset + zetRepo)
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
	fmt.Println(term.Blue + "Zet GitHub Remote: " + term.Reset + zetRemote)
	return nil
}

// Isosec returns the GMT current time in ISO8601 (RFC3339) without
// any punctuation or the T.  This is frequently a very good unique
// suffix that has the added advantage of being chronologically sortable
// and more readable than the epoch.
func Isosec() string {
	return fmt.Sprintf("%v", time.Now().In(time.UTC).Format("20060102150405"))
}

// mkdir is the functional equivalent of 'mkdir -p' and is used to create new
// folders recursively.
func mkdir(path string) error {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return err
	}
	return nil
}
