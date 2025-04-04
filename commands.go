package zet

import (
	"fmt"
	"github.com/charmbracelet/glamour"
	"github.com/danielmichaels/zet-cmd/internal/term"
	"os"
	"path/filepath"
	"regexp"
)

type Globals struct {
	Verbose bool `help:"Enable verbose mode" short:"v"`
}

type CreateCmd struct {
	Title string `arg:"" help:"Title of the zet to create"`
}

func (c *CreateCmd) Run() error {
	z := Zet{Title: c.Title}

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
	err = term.Exec(Editor, zet)
	if err != nil {
		return err
	}
	err = z.scanAndCommit(z.Path)
	if err != nil {
		return err
	}
	return nil

}

type LastCmd struct{}

func (c *LastCmd) Run() error {
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
}

type EditCmd struct {
	Isosec string `help:"Enter a valid isosec value (e.g. 20220424000235) and it will be opened using your system editor"`
	Last   bool   `help:"Open most recent zet using system editor"`
}

func (c *EditCmd) Run() error {
	z := new(Zet)

	if c.Last {
		err := z.ChangeDir(z.GetRepo())
		if err != nil {
			return err
		}

		last, err := z.Last()
		if err != nil {
			return err
		}

		err = z.openZetForEdit(last)
		if err != nil {
			return err
		}

		err = z.scanAndCommit(last)
		if err != nil {
			return err
		}

		return nil
	}

	r := regexp.MustCompile(zetRegex)

	if r.MatchString(c.Isosec) {
		zet, err := z.GetZet(c.Isosec)
		if err != nil {
			return err
		}
		err = z.openZetForEdit(zet)
		if err != nil {
			return err
		}

		err = z.scanAndCommit(zet)
		if err != nil {
			return err
		}
		return nil
	}
	err := z.edit(c.Isosec)
	if err != nil {
		return err
	}
	return nil
}

type FindCmd struct {
	Query string `arg:"" help:"Query to search"`
}

func (c *FindCmd) Run() error {
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
	results, err := z.SearchTitles(c.Query, titles)
	if err != nil {
		return err
	}
	for _, v := range results {
		fmt.Println(v.Id, v.Title)
	}
	return nil
}

type TagsCmd struct {
	Query string `arg:"" help:"Query to search"`
}

func (c *TagsCmd) Run() error {
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
	results, err := z.FindTags(c.Query, files)
	if err != nil {
		return err
	}
	for _, v := range results {
		fmt.Println(v.Id, v.Title)
	}

	return nil
}

type ViewCmd struct {
	Query string `arg:"" help:"Query to search"`
}

func (c *ViewCmd) Run() error {
	z := new(Zet)
	r := regexp.MustCompile(zetRegex)

	if r.MatchString(c.Query) {
		// Allow render of specific isosec
		zet, err := z.GetZet(c.Query)
		if err != nil {
			return err
		}
		file := filepath.Join(Repo, zet, "README.md")
		r, err := glamour.NewTermRenderer(
			glamour.WithAutoStyle(), glamour.WithWordWrap(zetWordWrap),
		)
		if err != nil {
			return err
		}
		c, err := os.ReadFile(file)
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
	err := z.render(c.Query)
	if err != nil {
		return err
	}
	return nil
}

type ViewAllCmd struct{}

func (c *ViewAllCmd) Run() error {
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
}

type CheckCmd struct{}

func (c *CheckCmd) Run() error {
	z := new(Zet)
	err := z.CheckZetConfig()
	if err != nil {
		return err
	}
	return nil
}
