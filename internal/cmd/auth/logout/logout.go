package logout

import (
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
)

type logoutOptions struct {
}

func NewCmdLogout(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Logout from Zeabur",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLogout(f)
		},
	}

	return cmd
}

// Note: don't import other packages directly in this function, or it will be hard to mock and test
// If you want to add new dependencies, please add them in the LogoutOptions struct

func runLogout(f *cmdutil.Factory) error {
	if !f.LoggedIn() {
		f.Log.Warnf("Not logged in, nothing to do")
		return nil
	}

	f.Config.SetTokenString("")
	f.Config.SetUser("")

	// reset token detail if exists
	if f.Config.GetToken() != nil {
		f.Config.SetToken(nil)
	}

	f.Log.Info("Logout successful!")
	return nil
}
