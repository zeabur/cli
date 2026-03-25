package list

import (
	"context"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/pkg/model"
)

type Options struct {
	status  string
	jobType string
	jobID   string
	page    int
	perPage int
}

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	opts := Options{page: 1, perPage: 20}

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List email sending records",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.status, "status", "", "Filter by status (pending/sent/delivered/bounced/complained)")
	cmd.Flags().StringVar(&opts.jobType, "job-type", "", "Filter by job type (direct/scheduled/batch)")
	cmd.Flags().StringVar(&opts.jobID, "job-id", "", "Filter by scheduled email ID or batch job ID")
	cmd.Flags().IntVar(&opts.page, "page", 1, "Page number")
	cmd.Flags().IntVar(&opts.perPage, "per-page", 20, "Items per page (max 100)")

	return cmd
}

func runList(f *cmdutil.Factory, opts Options) error {
	var status, jobType, jobID *string
	if opts.status != "" {
		status = &opts.status
	}
	if opts.jobType != "" {
		jobType = &opts.jobType
	}
	if opts.jobID != "" {
		jobID = &opts.jobID
	}

	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Fetching email records..."),
	)
	s.Start()
	reply, err := f.ApiClient.ListZSendEmails(context.Background(), &opts.page, &opts.perPage, status, jobType, jobID)
	s.Stop()
	if err != nil {
		return err
	}

	emails := make(model.ZSendEmails, 0, len(reply.Emails))
	for i := range reply.Emails {
		emails = append(emails, &reply.Emails[i])
	}

	if f.JSON {
		if len(emails) == 0 {
			return f.Printer.JSON([]any{})
		}
		return f.Printer.JSON(emails)
	}

	f.Printer.Table(emails.Header(), emails.Rows())
	return nil
}
