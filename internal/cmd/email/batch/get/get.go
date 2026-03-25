package get

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
)

func NewCmdGet(f *cmdutil.Factory) *cobra.Command {
	var apiKey string

	cmd := &cobra.Command{
		Use:   "get <job-id>",
		Short: "Get details of a batch email job",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGet(f, apiKey, args[0])
		},
	}

	cmd.Flags().StringVar(&apiKey, "api-key", "", "Z-Send API key (or set ZSEND_API_KEY)")

	return cmd
}

func runGet(f *cmdutil.Factory, apiKey, id string) error {
	if apiKey == "" {
		apiKey = os.Getenv("ZSEND_API_KEY")
	}
	if apiKey == "" {
		return fmt.Errorf("Z-Send API key is required (--api-key or ZSEND_API_KEY)")
	}
	if !strings.HasPrefix(apiKey, "zs_") {
		return fmt.Errorf("invalid API key format: must start with zs_")
	}

	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Fetching batch job..."),
	)
	s.Start()
	job, err := f.ApiClient.GetZSendBatchEmailJob(context.Background(), apiKey, id)
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
