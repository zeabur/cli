// Package auth contains the cmd for managing authentication
package auth

import (
	"github.com/spf13/cobra"

	authLoginCmd "github.com/zeabur/cli/internal/cmd/auth/login"
	authLogoutCmd "github.com/zeabur/cli/internal/cmd/auth/logout"
	authStatusCmd "github.com/zeabur/cli/internal/cmd/auth/status"
	"github.com/zeabur/cli/internal/cmdutil"
)

func NewCmdAuth(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth <command>",
		Short: "Authenticate Zeabur with browser or token",
	}

	// all sub-commands do not require authentication
	cmdutil.DisableAuthCheck(cmd)

	cmd.AddCommand(authLoginCmd.NewCmdLogin(f))
	cmd.AddCommand(authLogoutCmd.NewCmdLogout(f))
	cmd.AddCommand(authStatusCmd.NewCmdStatus(f))
	// cmd.AddCommand(authTokenCmd.NewCmdToken(f, nil))

	return cmd
}
