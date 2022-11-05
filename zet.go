// Copyright 2022 zet-cmd Authors
// SPDX-License-Identifier: Apache-2.0

package zet

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/charmbracelet/glamour"
	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/help"
	"github.com/rwxrob/term"
	url2 "net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"time"
)

const (
	zetRegex    = "^[0-9]{14,}$"
	zetWordWrap = 73
	zetUrl      = "https://github.com/danielmichaels/zet/blob/main/%s/README.md"
)

var (
	Pager       = os.Getenv("PAGER")
	Editor      = os.Getenv("EDITOR")
	GitUser     = os.Getenv("GITUSER")
	RepoName    = "zet"
	ZetRepo     = filepath.Join(os.Getenv("ZETDIR"))
	Pictures    = filepath.Join(os.Getenv("HOME"), "Pictures", "zet")
	Screenshots = filepath.Join(os.Getenv("HOME"), "Pictures", "zet")
	Downloads   = filepath.Join(os.Getenv("HOME"), "Downloads")
)

var CreateCmd = &Z.Cmd{
	Name:    `create`,
	Aliases: []string{"new", "c"},
	Summary: `Create a new zet`,
	Dynamic: template.FuncMap{"editor": func() string { return Editor }},
	Description: `
			Create a new zettelkasten entry by passing in the title as the only argument.

			This will then create a new README.md in a directory with
			an *isosec* timestamp as its name and use the system editor 
			({{ editor }}) to drop you into the file.

			**Multi-word titles must be encapsulated within quotations.** 
`,
	Other: []Z.Section{
		{`Examples`, `
					zet create "A New Thing"

					zet c title
`},
	},
	Commands: []*Z.Cmd{help.Cmd},
	Usage:    `must provide a title for each zet`,
	MinArgs:  1,
	MaxArgs:  1,
	Call: func(caller *Z.Cmd, args ...string) error {
		z := Zet{Title: args[0]}

		dir, err := z.CreateDir()
		if err != nil {
			return err
		}
		z.Path = dir

		err = z.CreateReadme(z, dir)
		if err != nil {
			return err
		}

		// Drop into vim and write Zet contents
		zet := z.GetReadme(z.Path)
		err = Z.Exec(Editor, zet)
		if err != nil {
			return err
		}
		err = z.scanAndCommit(z.Path)
		if err != nil {
			return err
		}
		return nil
	},
}

var GetCmd = &Z.Cmd{
	Name:    `get`,
	Aliases: []string{"g"},
	Summary: `Retrieve a zet for editing`,
	MinArgs: 1,
	Dynamic: template.FuncMap{"isosec": func() string { return Isosec() }},
	Other: []Z.Section{
		{`Examples`, `zet get {{ isosec }}`},
	},
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

var LastCmd = &Z.Cmd{
	Name:    `last`,
	Aliases: []string{"l", "latest"},
	Summary: `Get the most recent zet isosec and print it screen`,
	Description: `
			Prints the last modified zet entry's isosec value
`,
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

var ViewCmd = &Z.Cmd{
	Name:    `view`,
	Aliases: []string{"v"},
	Summary: `view command for zet entries.`,
	Description: `
			View supports both direct 'isosec' lookup's and keyword searches. 

			If a valid entry is found the markdown will be rendered in the terminal.
`,
	Commands: []*Z.Cmd{help.Cmd, viewAll},
	Call: func(caller *Z.Cmd, args ...string) error {
		z := new(Zet)
		r := regexp.MustCompile(zetRegex)

		if r.MatchString(args[0]) {
			// Allow render of specific isosec
			zet, err := z.GetZet(args[0])
			if err != nil {
				return err
			}
			file := filepath.Join(ZetRepo, zet, "README.md")
			c, err := os.ReadFile(file)
			if err != nil {
				return err
			}
			r, err := glamour.NewTermRenderer(
				glamour.WithAutoStyle(), glamour.WithWordWrap(zetWordWrap),
			)
			if err != nil {
				return err
			}
			_, err = r.Render(string(c))
			if err != nil {
				return err
			}
			return nil
		}
		err := z.render(args[0])
		if err != nil {
			return err
		}
		return nil
	},
}

func (z *Zet) render(arg string) error {
	zet, err := z.searchScanner(arg)
	if err != nil {
		return err
	}
	p := z.GetReadme(filepath.Join(ZetRepo, zet))
	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(), glamour.WithWordWrap(zetWordWrap),
	)
	if err != nil {
		return err
	}
	c, err := os.ReadFile(p)
	out, err := r.Render(string(c))
	fmt.Print(out)
	return nil
}

var viewAll = &Z.Cmd{
	Name:    `all`,
	Summary: `view all zet entries`,
	Description: `
			Output all zet entries from the local git repo.
			
			To view the response in the system pager ({{ pager }}) you must
			pipe the response using the terminal.
				

			**zet view all | {{ pager }}**
`,
	Dynamic: template.FuncMap{
		"pager": func() string { return Pager },
	},
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
		for _, v := range titles {
			fmt.Println(v.Id, v.Title)
		}
		return nil
	},
}

var QueryCmd = &Z.Cmd{
	Name:    `query`,
	Aliases: []string{"q"},
	Summary: `create a searchable URL with a query string`,
	Description: `
			Create a URL with a search term for the remote git hosting provider.
			
			Must place multi-word search terms inside quotations.
`,
	Other: []Z.Section{
		{`Examples`, `
			zet query "Multi-word must be quoted"

			*Outputs:* https://github.com/danielmichaels/zet/search?q=multi-word+must+be+quoted`,
		},
	},
	MinArgs:  1,
	MaxArgs:  1,
	Usage:    `must provide a search term`,
	Commands: []*Z.Cmd{help.Cmd},
	Call: func(caller *Z.Cmd, args ...string) error {
		term := url2.QueryEscape(args[0])
		url := fmt.Sprintf("https://github.com/%s/%s/search?q=%s", GitUser, RepoName, strings.ToLower(term))
		fmt.Println(url)
		return nil
	},
}

var FindCmd = &Z.Cmd{
	Name:    `find`,
	Aliases: []string{"f"},
	Summary: `Find a zet title by search term`,
	Description: `
			Search for a zet by title and retrieve any entries with that term
			in the title.
			Also captures partial matches, so "go" will find "golang".
			
			Only prints entries to the terminal.
`,
	MinArgs:  1,
	MaxArgs:  1,
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

var TagsCmd = &Z.Cmd{
	Name:    `tags`,
	Aliases: []string{"t"},
	Summary: `Find zet(s) by tag'`,
	Description: `
			Search for a zet by tag and retrieve any entries with that tag.
			Also captures partial matches, so "go" will find "golang".
			
			Only prints entries to the terminal.
`,
	MinArgs:  1,
	MaxArgs:  1,
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

var LinkCmd = &Z.Cmd{
	Name:     `link`,
	Summary:  `create a browsable link to a zet`,
	MinArgs:  1,
	MaxArgs:  1,
	Usage:    `zet link <search-term>`,
	Commands: []*Z.Cmd{help.Cmd},
	Call: func(caller *Z.Cmd, args ...string) error {
		z := new(Zet)
		zet, err := z.linkScanner(args[0])
		if err != nil {
			return err
		}

		title := zet.Title
		url := fmt.Sprintf(zetUrl, zet.Id)
		result := fmt.Sprintf("%s - %s\n", title, url)
		Z.PrintMark(result)
		return nil
	},
}

var CheckCmd = &Z.Cmd{
	Name:     `check`,
	Summary:  `check environment variables and configuration`,
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

// Zet is the struct to hang methods from which are used to create, edit, find
// and delete Zet's.
type Zet struct {
	Title  string
	Path   string
	Latest string
}

// GetRepo returns the GitRepo and RepoName as a filepath.
func (z *Zet) GetRepo() string { return ZetRepo }

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
	path := filepath.Join(ZetRepo, Isosec())
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
	z.Path = filepath.Join(ZetRepo, last)
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
	fmt.Println(term.Blue + "ZetRepo: " + term.Reset + ZetRepo)
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

// mkdir is the functional equivalent of 'mkdir -p' and is used to create new
// folders recursively.
func mkdir(path string) error {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return err
	}
	return nil
}
