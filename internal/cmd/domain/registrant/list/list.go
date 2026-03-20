package list

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
)

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List registrant profiles",
		Aliases: []string{"ls"},
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(f)
		},
	}
	return cmd
}

func runList(f *cmdutil.Factory) error {
	profiles, err := f.ApiClient.ListRegistrantProfiles(context.Background())
	if err != nil {
		return fmt.Errorf("list registrant profiles failed: %w", err)
	}

	if len(profiles) == 0 {
		if f.JSON {
			return f.Printer.JSON([]any{})
		}
		f.Log.Infof("No registrant profiles found")
		return nil
	}

	if f.JSON {
		return f.Printer.JSON(profiles)
	}

	f.Printer.Table(profiles.Header(), profiles.Rows())
	return nil
}
