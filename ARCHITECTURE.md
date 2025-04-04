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

## Notes
Todo:

- Git Pull the repo (incase it gets out of sync)
- Delete zet (saves needing to do it in browser or using `git`)
- Search tags
