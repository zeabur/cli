package usage

import (
	"context"
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/pkg/model"
)

type Options struct {
	month string
}

func NewCmdUsage(f *cmdutil.Factory) *cobra.Command {
	opts := Options{}

	cmd := &cobra.Command{
		Use:   "usage",
		Short: "Show AI Hub monthly usage",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUsage(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.month, "month", "", "Month in YYYY-MM format")

	return cmd
}

func runUsage(f *cmdutil.Factory, opts Options) error {
	var monthPtr *string
	if opts.month != "" {
		if _, err := time.Parse("2006-01", opts.month); err != nil {
			return fmt.Errorf("--month must be in YYYY-MM format: %w", err)
		}
		monthPtr = &opts.month
	}

	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Fetching AI Hub usage..."),
	)
	s.Start()
	usage, err := f.ApiClient.GetAIHubMonthlyUsage(context.Background(), monthPtr)
	s.Stop()
	if err != nil {
		return err
	}

	if f.JSON {
		return f.Printer.JSON(usage)
	}

	f.Log.Infof("Total Spend: $%.6f", usage.TotalSpend)

	if len(usage.ModelsCost) > 0 {
		f.Log.Infof("")
		f.Log.Infof("Per-Model Breakdown:")
		costs := model.AIHubModelCosts(usage.ModelsCost)
		f.Printer.Table(costs.Header(), costs.Rows())
	}

	return nil
}
