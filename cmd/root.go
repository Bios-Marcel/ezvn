package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func Execute(out io.Writer) error {
	// This seems to provide no value whatsoever, it seemingly doesn't even do
	// what's documented. All it does, is take time.
	cobra.MousetrapHelpText = ""

	root := &cobra.Command{
		Use:   "ezvn",
		Short: "Wrapper around SVN, that offers the same functionality, but better.",
		// By default, subcommand aliases aren't autocompleted.
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			var aliases []string
			for _, subCmd := range cmd.Commands() {
				aliases = append(aliases, subCmd.Aliases...)
			}
			return aliases, cobra.ShellCompDirectiveNoFileComp
		},
		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: true,
		},
	}
	root.AddCommand(NewChangelist(), NewUndo(), NewPurge())
	err := root.Execute()
	if err != nil && strings.HasPrefix(err.Error(), "unknown command") {
		fmt.Println("Delegating to svn ...")
		// If a subcommand is unknown, we show redirect to svn instead, as
		// it could be an original command.
		redirectCmd := createCommand("svn", os.Args[1:]...)
		return redirectCmd.Run()
	}
	return err
}
