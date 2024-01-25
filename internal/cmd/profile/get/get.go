package get

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
)

// profileOptions contains the input to the profile command.
type profileOptions struct {
	verbose bool
}

func NewCmdProfile(f *cmdutil.Factory) *cobra.Command {
	opts := &profileOptions{}
	cmd := &cobra.Command{
		Use:   "info",
		Short: "Show User Profile Info",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStatus(f, opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.verbose, "verbose", "v", false, "Show more details")

	return cmd
}

func runStatus(f *cmdutil.Factory, opts *profileOptions) error {
	user, err := f.ApiClient.GetUserInfo(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get user info: %w", err)
	}

	f.Printer.Table(user.Header(), user.Rows())

	return nil
}
