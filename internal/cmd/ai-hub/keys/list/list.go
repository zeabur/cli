package list

import (
	"context"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/pkg/model"
)

type Options struct{}

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	opts := Options{}

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List AI Hub API keys",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(f, opts)
		},
	}

	return cmd
}

func runList(f *cmdutil.Factory, opts Options) error {
	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Fetching AI Hub keys..."),
	)
	s.Start()
	tenant, err := f.ApiClient.GetAIHubTenant(context.Background())
	s.Stop()
	if err != nil {
		return err
	}

	keys := model.AIHubKeys(tenant.Keys)

	if f.JSON {
		if len(keys) == 0 {
			return f.Printer.JSON([]any{})
		}
		return f.Printer.JSON(keys)
	}

	f.Printer.Table(keys.Header(), keys.Rows())

	return nil
}
