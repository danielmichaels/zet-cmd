Zet

> zettelkasten is a place to store small knowledge 'slips' for recall later

## Design

Zet will use github and git as its storage medium.  Each new zet will get a directory containing a README.md. this ensures a nice viewing experience when accessing that directory from github.

 The created directory will use an iso seconds format.
 
  /zet/2022010112000/README.md 
  
  
  
  $contents
  
  Title
  
  Body
  
  Tags
  
  Example (in codeblocks but must be markdown)
  
 

  #
 This is a Zet title
  
  The body goes here.
  
  It can contain any valid markdown.
  
  There must be a space between body end and the Zet tags
  
  > #tag1 #tag2 #tag3
 

Tags are used for quicker searching. Terminal tools and github do well searching for hashtag prepended names.

### Implementation

Many facets will be stolen from rwxrob/Zet-cmd
https://github.com/rwxrob/cmd-zet

This will be a standalone bonzai branch which can be imported into ds



## Commands

Zet commands required to create, read and edit zet's

### Create

Creating a new Zet is the most common command. 

create flow
  
  1. Zet create $1
  2. Create $dir
  3. Vim $dir/README.md
  4. Write contents (see $contents)
  5. Git commit -s -a -m $1
  6. Git pull -q (prevent conflicts)
  7. Git push -s
  8. Echo zet pushed as $isosec ($dir name)
 
### Edit

Editing any previous Zet is achieved by appending the isosec after Zet edit. 

### Last

This should open the last Zet in vim.

Finding the last Zet could probably be achieved by identifying the latest dir from the list of dirs

### View

This could be a search option for reading or finding a Zet to edit

## Notes
Todo:

- Git Pull the repo (incase it gets out of sync)
- Delete zet (saves needing to do it in browser or using `git`)
- Search tags
