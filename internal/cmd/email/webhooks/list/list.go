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
		Short:   "List email webhooks",
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
		spinner.WithSuffix(" Fetching webhooks..."),
	)
	s.Start()
	reply, err := f.ApiClient.ListZSendWebhooks(context.Background(), nil, nil)
	s.Stop()
	if err != nil {
		return err
	}

	webhooks := make(model.ZSendWebhooks, 0, len(reply.Webhooks))
	for i := range reply.Webhooks {
		webhooks = append(webhooks, &reply.Webhooks[i])
	}

	if f.JSON {
		if len(webhooks) == 0 {
			return f.Printer.JSON([]any{})
		}
		return f.Printer.JSON(webhooks)
	}

	f.Printer.Table(webhooks.Header(), webhooks.Rows())

	return nil
}
