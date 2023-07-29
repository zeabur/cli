package status

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/pkg/api"
)

// statusOptions contains the input to the status command.
type statusOptions struct {
	verbose bool
}

func NewCmdStatus(f *cmdutil.Factory) *cobra.Command {
	opts := &statusOptions{}
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show Zeabur login status",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStatus(f, opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.verbose, "verbose", "v", false, "Show more details")

	return cmd
}

// Note: don't import other packages directly in this function, or it will be hard to mock and test
// If you want to add new dependencies, please add them in the statusOptions struct

func runStatus(f *cmdutil.Factory, opts *statusOptions) error {
	if !f.LoggedIn() {
		f.Log.Infof("Not logged in.")
		return nil
	}

	f.ApiClient = api.New(f.Config.GetTokenString())

	user, err := f.ApiClient.GetUserInfo(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get user info: %w", err)
	}

	f.Log.Infof("Logged in as %s, email: %s", user.Name, user.Email)

	if opts.verbose {
		f.Printer.Table(user.Header(), user.Rows())
	}

	return nil
}
