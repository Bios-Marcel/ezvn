package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewPurge() *cobra.Command {
	return &cobra.Command{
		Use:   "purge",
		Short: "Removes all changes, including unversioned files",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Revert all changes excluding untracked files.
			if err := createCommand("svn", "revert", "--recursive", ".").Run(); err != nil {
				return fmt.Errorf("error reverting changes: %w", err)
			}
			// Remove untracked files
			if err := createCommand("svn", "cleanup", ".", "--remove-unversioned").Run(); err != nil {
				return fmt.Errorf("error cleaning up: %w", err)
			}
			return nil
		},
	}
}
