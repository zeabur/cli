// Package version contains the cmd for managing the
package version

import (
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
)

func NewCmdVersion(f *cmdutil.Factory, version, commit, date string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "version",
		Short:   "Print the version number of Zeabur CLI",
		Aliases: []string{"v", "ver"},
		Run: func(cmd *cobra.Command, args []string) {
			f.Printer.Table([]string{"Version", "Commit", "Date"}, [][]string{{version, commit, date}})
		},
	}

	// no authentication required
	cmdutil.DisableAuthCheck(cmd)

	return cmd
}
