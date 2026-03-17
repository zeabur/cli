package aihub

import (
	"github.com/spf13/cobra"

	addBalanceCmd "github.com/zeabur/cli/internal/cmd/ai-hub/add-balance"
	autoRechargeCmd "github.com/zeabur/cli/internal/cmd/ai-hub/auto-recharge"
	keysCmd "github.com/zeabur/cli/internal/cmd/ai-hub/keys"
	statusCmd "github.com/zeabur/cli/internal/cmd/ai-hub/status"
	usageCmd "github.com/zeabur/cli/internal/cmd/ai-hub/usage"
	"github.com/zeabur/cli/internal/cmdutil"
)

func NewCmdAIHub(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ai-hub <command>",
		Short: "Manage AI Hub",
	}

	cmd.AddCommand(statusCmd.NewCmdStatus(f))
	cmd.AddCommand(keysCmd.NewCmdKeys(f))
	cmd.AddCommand(addBalanceCmd.NewCmdAddBalance(f))
	cmd.AddCommand(usageCmd.NewCmdUsage(f))
	cmd.AddCommand(autoRechargeCmd.NewCmdAutoRecharge(f))

	return cmd
}
