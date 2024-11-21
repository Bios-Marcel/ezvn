package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/Bios-Marcel/ezvn/svn"
	"github.com/spf13/cobra"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func NewChangelist() *cobra.Command {
	addCmd := &cobra.Command{
		Use:     "add",
		Aliases: []string{"track"},
		// Name of the changelist to add to and at least one file to be added.
		Args: cobra.MinimumNArgs(2),
		// Disabled til we have a proper impl
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Delegating to svn ...")
			redirectCmd := createCommand("svn", append([]string{"changelist"}, args...)...)
			return redirectCmd.Run()
		},
	}
	removeCmd := &cobra.Command{
		Use:     "remove",
		Aliases: []string{"untrack"},
		// At least one file should be removed
		Args: cobra.MinimumNArgs(1),
		// Disabled til we have a proper impl
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Delegating to svn ...")
			redirectCmd := createCommand("svn", append([]string{"changelist", "--remove"}, args...)...)
			return redirectCmd.Run()
		},
	}
	clearCmd := &cobra.Command{
		Use: "clear",
		// Exactly one changelist can be cleared
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			status, err := svn.GetStatus()
			if err != nil {
				return fmt.Errorf("error retrieving SVN status: %w", err)
			}

			toRemove := args[0]
			for _, changelist := range status.Changelists {
				// FIXME Add proper fuzzy matching
				if !strings.HasPrefix(strings.ToLower(changelist.Name), toRemove) {
					continue
				}

				var filesAsArgs []string
				for _, file := range changelist.Files {
					filesAsArgs = append(filesAsArgs, file.Path)
				}
				argFile, err := createArgFile(filesAsArgs)
				if err != nil {
					fmt.Printf("error creating argfile: %s", err)
					os.Exit(1)
				}
				defer os.Remove(argFile)

				svnChangelistRemove := exec.Command("svn", "changelist", "--remove", "--targets", argFile)
				svnChangelistRemove.Stdout = os.Stdout
				svnChangelistRemove.Stderr = os.Stderr
				svnChangelistRemove.Stdin = os.Stdin
				if err := svnChangelistRemove.Run(); err != nil {
					return fmt.Errorf("error removing changelist: %w", err)
				}
				break
			}
			return nil
		},
	}
	changelistCmd := &cobra.Command{
		Use:     "changelist",
		Aliases: []string{"changelists", "changesets", "changeset", "cl"},
		RunE: func(cmd *cobra.Command, args []string) error {
			status, err := svn.GetStatus()
			if err != nil {
				return fmt.Errorf("error retrieving SVN status: %w", err)
			}

			for _, changelist := range status.Changelists {
				fmt.Println(changelist.Name)
				for _, file := range changelist.Files {
					fmt.Printf("\t%s\n", file.Path)
				}
			}

			return nil
		},
	}
	changelistCmd.AddCommand(addCmd, removeCmd, clearCmd)

	return changelistCmd
}

func createArgFile(args []string) (string, error) {
	tempDir := os.TempDir()
	argFile, err := os.CreateTemp(tempDir, "svn_changelist_*.args")
	if err != nil {
		return "", fmt.Errorf("Error creating tempfile: %w", err)
	}
	defer argFile.Close()

	shittyEncodingWriter := transform.NewWriter(argFile, charmap.Windows1250.NewEncoder())
	defer shittyEncodingWriter.Close()
	for _, arg := range args {
		if _, err := io.WriteString(shittyEncodingWriter, arg+"\n"); err != nil {
			return "", fmt.Errorf("error writing to arg file: %w", err)
		}
	}

	return argFile.Name(), nil
}
