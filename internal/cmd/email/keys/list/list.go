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
		Short:   "List email API keys",
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
		spinner.WithSuffix(" Fetching API keys..."),
	)
	s.Start()
	reply, err := f.ApiClient.ListZSendAPIKeys(context.Background(), nil, nil)
	s.Stop()
	if err != nil {
		return err
	}

	keys := make(model.ZSendAPIKeys, 0, len(reply.APIKeys))
	for i := range reply.APIKeys {
		keys = append(keys, &reply.APIKeys[i])
	}

	if f.JSON {
		if len(keys) == 0 {
			return f.Printer.JSON([]any{})
		}
		return f.Printer.JSON(keys)
	}

	f.Printer.Table(keys.Header(), keys.Rows())

	return nil
}
