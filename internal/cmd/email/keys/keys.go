package keys

import (
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"

	keysCreateCmd "github.com/zeabur/cli/internal/cmd/email/keys/create"
	keysDeleteCmd "github.com/zeabur/cli/internal/cmd/email/keys/delete"
	keysGetCmd "github.com/zeabur/cli/internal/cmd/email/keys/get"
	keysListCmd "github.com/zeabur/cli/internal/cmd/email/keys/list"
)

func NewCmdKeys(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "keys",
		Short:   "Manage email API keys",
		Aliases: []string{"key"},
	}

	cmd.AddCommand(keysListCmd.NewCmdList(f))
	cmd.AddCommand(keysGetCmd.NewCmdGet(f))
	cmd.AddCommand(keysCreateCmd.NewCmdCreate(f))
	cmd.AddCommand(keysDeleteCmd.NewCmdDelete(f))

	return cmd
}
