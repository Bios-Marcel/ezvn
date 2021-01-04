package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) == 1 {
		createCommand("svn", "help").Run()
		fmt.Print("\n\n")
		fmt.Println("ezvn extension commands:")
		fmt.Println("\tundo (uncommit) - removes changes made in a revision or range of revisions. Expects comma separated numbers or ranges in format of FROM:TO")
		fmt.Println("\tpurge - removes all local changes including untracked files")
		return
	}

	firstArg := os.Args[1]
	if strings.EqualFold("undo", firstArg) || strings.EqualFold("uncommit", firstArg) {
		if len(os.Args) <= 2 {
			panic("not enough arguments")
		}

		var builder bytes.Buffer
		for _, arg := range os.Args[2:] {
			builder.WriteString(strings.TrimSpace(arg))
		}

		commaSeparatedRevisions, parseError := parseRevisionsArgument(strings.Split(builder.String(), ","))
		if parseError != nil {
			fmt.Printf("Invalider input:\n\t%s\n", parseError)
			os.Exit(0)
		}

		mergeCommand := createCommand("svn", "merge", "-c", commaSeparatedRevisions, ".")
		executeError := mergeCommand.Run()
		if executeError != nil {
			panic(executeError)
		}
	} else if strings.EqualFold("purge", firstArg) {
		//Revert all changes excluding untracked files.
		createCommand("svn", "revert", "--recursive", ".").Run()
		//Remove untracked files
		createCommand("svn", "cleanup", ".", "--remove-unversioned").Run()
	} else {
		svnRedirectCommand := createCommand("svn", os.Args[1:]...)
		svnRedirectCommand.Run()
	}
}

func parseRevisionsArgument(revisions []string) (string, error) {
	var changeSetsToRemove bytes.Buffer
	for index, revision := range revisions {
		if index != 0 && index != len(revision)-1 {
			changeSetsToRemove.WriteRune(',')
		}

		doubleColonCount := strings.Count(revision, ":")
		if doubleColonCount == 0 {
			_, parseError := strconv.ParseInt(revision, 10, 64)
			if parseError != nil {
				return "", errors.New("revision is not numeric")
			}
			changeSetsToRemove.WriteRune('-')
			changeSetsToRemove.WriteString(revision)
			continue
		}

		if doubleColonCount > 2 {
			return "", errors.New("too many double colons. Revisions must be of format XX:YY and be separate by commas: XX:YY,AA:BB")
		}

		doubleColonIndex := strings.IndexRune(revision, ':')
		revFrom := revision[:doubleColonIndex]
		revTo := revision[doubleColonIndex+1:]

		parsedFrom, parseErrFrom := strconv.ParseInt(revFrom, 10, 64)
		if parseErrFrom != nil {
			return "", parseErrFrom
		}

		parsedTo, parseErrTo := strconv.ParseInt(revTo, 10, 64)
		if parseErrTo != nil {
			return "", parseErrTo
		}

		for from := parsedFrom; from <= parsedTo; from++ {
			if from != parsedFrom {
				changeSetsToRemove.WriteRune(',')
			}
			changeSetsToRemove.WriteRune('-')
			changeSetsToRemove.WriteString(strconv.FormatInt(from, 10))
		}
	}

	return changeSetsToRemove.String(), nil
}

func createCommand(executable string, args ...string) *exec.Cmd {
	command := exec.Command(executable, args...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin
	return command
}
