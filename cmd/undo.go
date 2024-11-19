package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

func NewUndo() *cobra.Command {
	return &cobra.Command{
		Use:     "undo",
		Short:   "Reverts commits (revisions).",
		Aliases: []string{"uncommit"},
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var builder bytes.Buffer
			for _, arg := range args {
				builder.WriteString(strings.TrimSpace(arg))
			}

			commaSeparatedRevisions, err := parseRevisionsArgument(strings.Split(builder.String(), ","))
			if err != nil {
				return fmt.Errorf("invalid revision format: %w", err)
			}

			mergeCommand := createCommand("svn", "merge", "-c", commaSeparatedRevisions, ".")
			if err := mergeCommand.Run(); err != nil {
				return fmt.Errorf("error running svn merge: %w", err)
			}
			return nil
		},
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
