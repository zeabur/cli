package logout

import (
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
)

type logoutOptions struct{}

func NewCmdLogout(f *cmdutil.Factory) *cobra.Command {
	opts := &logoutOptions{}
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Logout from Zeabur",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLogout(f, opts)
		},
	}

	return cmd
}

// Note: don't import other packages directly in this function, or it will be hard to mock and test
// If you want to add new dependencies, please add them in the LogoutOptions struct

func runLogout(f *cmdutil.Factory, opts *logoutOptions) error {
	if !f.LoggedIn() {
		f.Log.Warnf("Not logged in, nothing to do")
		return nil
	}

	// reset token string and user
	f.Config.SetTokenString("")
	f.Config.SetUser("")
	f.Config.SetUsername("")

	// reset context
	f.Config.GetContext().ClearAll()

	f.Log.Info("Logout successful!")
	return nil
}
