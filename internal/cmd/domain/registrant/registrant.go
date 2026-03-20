package registrant

import (
	"github.com/spf13/cobra"

	registrantCreateCmd "github.com/zeabur/cli/internal/cmd/domain/registrant/create"
	registrantDeleteCmd "github.com/zeabur/cli/internal/cmd/domain/registrant/delete"
	registrantListCmd "github.com/zeabur/cli/internal/cmd/domain/registrant/list"
	registrantUpdateCmd "github.com/zeabur/cli/internal/cmd/domain/registrant/update"
	"github.com/zeabur/cli/internal/cmdutil"
)

func NewCmdRegistrant(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "registrant <command>",
		Short: "Manage registrant profiles",
	}

	cmd.AddCommand(registrantListCmd.NewCmdList(f))
	cmd.AddCommand(registrantCreateCmd.NewCmdCreate(f))
	cmd.AddCommand(registrantUpdateCmd.NewCmdUpdate(f))
	cmd.AddCommand(registrantDeleteCmd.NewCmdDelete(f))

	return cmd
}
