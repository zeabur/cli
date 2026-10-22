package list

import (
	"context"
	"fmt"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/pkg/model"
)

type Options struct {
	apiKey  string
	status  string
	page    int
	perPage int
}

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	opts := Options{page: 1, perPage: 20}

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List scheduled emails",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.apiKey, "api-key", "", cmdutil.ZSendAPIKeyFlagUsage)
	cmd.Flags().StringVar(&opts.status, "status", "", "Filter by status (scheduled/sent/cancelled/failed)")
	cmd.Flags().IntVar(&opts.page, "page", 1, "Page number")
	cmd.Flags().IntVar(&opts.perPage, "per-page", 20, "Items per page (max 100)")

	return cmd
}

func runList(f *cmdutil.Factory, opts Options) error {
	if opts.page < 1 {
		return fmt.Errorf("page must be >= 1")
	}
	if opts.perPage < 1 || opts.perPage > 100 {
		return fmt.Errorf("per-page must be between 1 and 100")
	}

	apiKey, err := f.ResolveZSendAPIKey(opts.apiKey)
	if err != nil {
		return err
	}
	opts.apiKey = apiKey

	var status *string
	if opts.status != "" {
		status = &opts.status
	}

	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Fetching scheduled emails..."),
	)
	s.Start()
	reply, err := f.ApiClient.ListZSendScheduledEmails(context.Background(), opts.apiKey, &opts.page, &opts.perPage, status)
	s.Stop()
	if err != nil {
		return err
	}

	items := make(model.ZSendScheduledEmails, 0, len(reply.ScheduledEmails))
	for i := range reply.ScheduledEmails {
		items = append(items, &reply.ScheduledEmails[i])
	}

	if f.JSON {
		if len(items) == 0 {
			return f.Printer.JSON([]any{})
		}
		return f.Printer.JSON(items)
	}

	f.Printer.Table(items.Header(), items.Rows())
	return nil
}
