package completion

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
)

const (
	ShellBash       = "bash"
	ShellZsh        = "zsh"
	ShellFish       = "fish"
	ShellPowerShell = "powershell"
)

func NewCmdCompletion(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "completion <shell>",
		Short:                 "Generate completion script",
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{ShellBash, ShellZsh, ShellFish, ShellPowerShell},
		Args:                  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCompletion(f, cmd, args[0])
		},
	}

	return cmd
}

// todo: test if this works
func runCompletion(f *cmdutil.Factory, cmd *cobra.Command, cmdType string) error {
	switch cmdType {
	case ShellBash:
		return cmd.GenBashCompletion(os.Stdout)
	case ShellZsh:
		return cmd.GenZshCompletion(os.Stdout)
	case ShellFish:
		return cmd.GenFishCompletion(os.Stdout, true)
	case ShellPowerShell:
		return cmd.GenPowerShellCompletion(os.Stdout)
	default:
		return fmt.Errorf("unsupported shell type %q", cmdType)
	}
}
