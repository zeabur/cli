package list

import (
	"context"
	"fmt"
	"os"
	"strings"

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
		Short:   "List batch email jobs",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.apiKey, "api-key", "", "Z-Send API key (or set ZSEND_API_KEY)")
	cmd.Flags().StringVar(&opts.status, "status", "", "Filter by status (pending/processing/completed/failed)")
	cmd.Flags().IntVar(&opts.page, "page", 1, "Page number")
	cmd.Flags().IntVar(&opts.perPage, "per-page", 20, "Items per page (max 100)")

	return cmd
}

func runList(f *cmdutil.Factory, opts Options) error {
	if opts.apiKey == "" {
		opts.apiKey = os.Getenv("ZSEND_API_KEY")
	}
	if opts.apiKey == "" {
		return fmt.Errorf("Z-Send API key is required (--api-key or ZSEND_API_KEY)")
	}
	if !strings.HasPrefix(opts.apiKey, "zs_") {
		return fmt.Errorf("invalid API key format: must start with zs_")
	}

	var status *string
	if opts.status != "" {
		status = &opts.status
	}

	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Fetching batch jobs..."),
	)
	s.Start()
	reply, err := f.ApiClient.ListZSendBatchEmailJobs(context.Background(), opts.apiKey, &opts.page, &opts.perPage, status)
	s.Stop()
	if err != nil {
		return err
	}

	jobs := make(model.ZSendBatchJobs, 0, len(reply.Jobs))
	for i := range reply.Jobs {
		jobs = append(jobs, &reply.Jobs[i])
	}

	if f.JSON {
		if len(jobs) == 0 {
			return f.Printer.JSON([]any{})
		}
		return f.Printer.JSON(jobs)
	}

	f.Printer.Table(jobs.Header(), jobs.Rows())
	return nil
}
