package update

import (
	"github.com/spf13/cobra"

	tagUpdateCmd "github.com/zeabur/cli/internal/cmd/service/update/tag"
	"github.com/zeabur/cli/internal/cmdutil"
)

// NewCmdUpdate creates the update command
func NewCmdUpdate(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <command>",
		Short:   "Update service",
		Aliases: []string{"up"},
	}

	cmd.AddCommand(tagUpdateCmd.NewCmdTag(f))

	return cmd
}
