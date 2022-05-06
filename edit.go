package zet

import (
	"fmt"
	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/help"
	"os"
	"path/filepath"
	"text/template"
)

var EditCmd = &Z.Cmd{
	Name:    `edit`,
	Aliases: []string{"e"},
	Summary: `edit a zet`,
	MinArgs: 1,
	Usage:   `must provide a zet isosec value`,
	Dynamic: template.FuncMap{"editor": func() string { return Editor }},
	Description: `
			Enter a valid isosec value (e.g. 20220424000235) and it will be opened using your system editor ({{ editor}}).
`,
	Commands: []*Z.Cmd{help.Cmd, editLast, findEdit},
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

var editLast = &Z.Cmd{
	Name:    `last`,
	Summary: `edit the last modified zet entry from the git repo`,
	Dynamic: template.FuncMap{"editor": func() string { return Editor }},
	Description: `
			Open the last modified zet using your system editor ({{ editor}}).
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

type Found struct {
	Index int
	Id    string
	Title string
}

var findEdit = &Z.Cmd{
	Name:    `find`,
	Summary: `search titles for a zet to edit`,
	Dynamic: template.FuncMap{"editor": func() string { return Editor }},
	Description: `
			Search for a zet title and retrieve any matching entry.
			
			To edit it, enter the index into the terminal and press enter. 
			This will open the file in the system editor ({{ editor }}).
`,
	Commands: []*Z.Cmd{help.Cmd},
	Call: func(caller *Z.Cmd, args ...string) error {
		z := new(Zet)

		err := z.ChangeDir(ZetRepo)
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
		var ff []Found
		for idx, v := range results {
			var f Found
			f.Index = idx
			f.Id = v.Id
			f.Title = v.Title
			ff = append(ff, f)
		}
		if len(ff) == 0 {
			fmt.Printf("No entries found for %q\n", args[0])
			return nil
		}
		for _, k := range ff {
			fmt.Printf("%d) %s %s\n", k.Index, k.Id, k.Title)
		}
		var s int
		fmt.Printf("#> ")
		_, err = fmt.Scanln(&s)
		if err != nil {
			switch err.Error() {
			case "unexpected newline":
				fmt.Println("Did not enter a value. Exiting.")
				return nil
			case "expected integer":
				fmt.Println("Must enter an integer")
				return nil
			default:
				return err
			}
		}

		var zet string
		for _, k := range ff {
			if s == k.Index {
				zet = k.Id
			}
		}
		if zet == "" {
			fmt.Println("Key entered does not match, or zet could not be found")
			return nil
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
