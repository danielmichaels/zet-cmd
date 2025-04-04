package zet

import (
	"fmt"
	"github.com/danielmichaels/zet-cmd/internal/term"
	"os"
	"path/filepath"
	"strconv"
)

// Found represents a search result for a zet note, containing its index, ID, and title.
type Found struct {
	Index int
	Id    string
	Title string
}

// openZetForEdit opens the README.md file of a specified zet note for editing using the configured editor.
func (z *Zet) openZetForEdit(zet string) error {
	file := filepath.Join(zet, "README.md")
	err := term.Exec(Editor, file)
	if err != nil {
		return err
	}
	return nil
}
func (z *Zet) edit(args ...string) error {
	err := z.searchScanner(args[0])
	if err != nil {
		return err
	}
	err = z.openZetForEdit(z.Path)
	if err != nil {
		return err
	}
	err = z.scanAndCommit(z.Path)
	if err != nil {
		return err
	}
	return nil
}

// searchScanner searches for zet notes matching the provided search term, displays matching results,
// and prompts the user to select a specific zet note. It updates the Zet struct with the selected
// note's path and title. If no matching notes are found or an invalid selection is made, the program exits.
func (z *Zet) searchScanner(args ...string) error {
	err := z.ChangeDir(Repo)
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
		os.Exit(0)
	}
	for _, k := range ff {
		fmt.Printf("%d) %s %s\n", k.Index, k.Id, k.Title)
	}
	prompt := term.Prompt("#> ")
	if prompt == "" {
		fmt.Println("exiting. did not provide valid entry.")
		os.Exit(0)
	}

	s, _ := strconv.Atoi(prompt)
	var zet string
	for _, k := range ff {
		if s == k.Index {
			zet = k.Id
			z.Path = k.Id
			z.Title = k.Title
		}
	}
	if zet == "" {
		fmt.Println("Key entered does not match, or zet could not be found")
		os.Exit(0)
	}
	return nil
}
