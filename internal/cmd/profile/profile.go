// Package profile contains the cmd for managing profile
package profile

import (
	"github.com/spf13/cobra"

	profileCmd "github.com/zeabur/cli/internal/cmd/profile/get"
	"github.com/zeabur/cli/internal/cmdutil"
)

// NewCmdProject creates the profile command
func NewCmdProfile(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profile",
		Short: "Get profile",
	}

	cmd.AddCommand(profileCmd.NewCmdProfile(f))

	return cmd
}
