package get

import (
	"context"
	"fmt"
	"strings"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct {
	apiKey string
	id     string
}

func NewCmdGet(f *cmdutil.Factory) *cobra.Command {
	opts := Options{}

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get details of a batch email job",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGet(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.apiKey, "api-key", "", cmdutil.ZSendAPIKeyFlagUsage)
	cmd.Flags().StringVar(&opts.id, "id", "", "Batch job ID")

	return cmd
}

func runGet(f *cmdutil.Factory, opts Options) error {
	if f.Interactive {
		return runGetInteractive(f, opts)
	}
	return runGetNonInteractive(f, opts)
}

func runGetInteractive(f *cmdutil.Factory, opts Options) error {
	if opts.id == "" {
		id, err := f.Prompter.Input("Batch Job ID: ", "")
		if err != nil {
			return err
		}
		opts.id = id
	}
	if err := paramCheck(opts); err != nil {
		return err
	}
	return getBatchJob(f, opts)
}

func runGetNonInteractive(f *cmdutil.Factory, opts Options) error {
	if err := paramCheck(opts); err != nil {
		return err
	}
	return getBatchJob(f, opts)
}

func getBatchJob(f *cmdutil.Factory, opts Options) error {
	apiKey, err := f.ResolveZSendAPIKey(opts.apiKey)
	if err != nil {
		return err
	}
	opts.apiKey = apiKey

	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Fetching batch job..."),
	)
	s.Start()
	job, err := f.ApiClient.GetZSendBatchEmailJob(context.Background(), opts.apiKey, opts.id)
	s.Stop()
	if err != nil {
		return err
	}

	if f.JSON {
		return f.Printer.JSON(job)
	}

	f.Printer.Table(
		[]string{"Field", "Value"},
		[][]string{
			{"Job ID", job.JobID},
			{"Status", job.Status},
			{"Total", fmt.Sprintf("%d", job.TotalCount)},
			{"Sent", fmt.Sprintf("%d", job.SentCount)},
			{"Failed", fmt.Sprintf("%d", job.FailedCount)},
			{"Created At", job.CreatedAt},
			{"Started At", job.StartedAt},
			{"Completed At", job.CompletedAt},
			{"Last Error", job.LastError},
		},
	)
	return nil
}

func paramCheck(opts Options) error {
	if strings.TrimSpace(opts.id) == "" {
		return fmt.Errorf("batch job ID is required")
	}
	return nil
}
