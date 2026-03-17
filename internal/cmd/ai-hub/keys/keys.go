package keys

import (
	"github.com/spf13/cobra"

	keysCreateCmd "github.com/zeabur/cli/internal/cmd/ai-hub/keys/create"
	keysDeleteCmd "github.com/zeabur/cli/internal/cmd/ai-hub/keys/delete"
	keysListCmd "github.com/zeabur/cli/internal/cmd/ai-hub/keys/list"
	"github.com/zeabur/cli/internal/cmdutil"
)

func NewCmdKeys(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "keys <command>",
		Short: "Manage AI Hub API keys",
	}

	cmd.AddCommand(keysListCmd.NewCmdList(f))
	cmd.AddCommand(keysCreateCmd.NewCmdCreate(f))
	cmd.AddCommand(keysDeleteCmd.NewCmdDelete(f))

	return cmd
}
