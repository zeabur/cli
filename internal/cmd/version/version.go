package version

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
)

func NewCmdVersion(f *cmdutil.Factory, version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "version",
		Short:   "Print the version number of Zeabur CLI",
		Aliases: []string{"v", "ver"},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version: %s\n", version)
		},
	}

	// no authentication required
	cmdutil.DisableAuthCheck(cmd)

	return cmd
}
