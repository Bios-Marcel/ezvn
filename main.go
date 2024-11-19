package main

import (
	"os"

	"github.com/Bios-Marcel/ezvn/cmd"
)

// Uses spaces for indentation, just like svn, since tabs have variable size.
// Just like svn, we list the name of each command as the first word, followed
// by, if available, the aliases in curved braces separated by commas.
const mainHelpPageExtension = `

ezvn extension commands:
    undo (uncommit) - removes changes made in a revision or range of revisions. Expects comma separated numbers or ranges in format of FROM:TO
    purge - removes all local changes including untracked files`

func main() {
	if err := cmd.Execute(os.Stdout); err != nil {
		os.Exit(1)
	}
}
