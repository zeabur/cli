package cancel

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct {
	apiKey string
	id     string
}

func NewCmdCancel(f *cmdutil.Factory) *cobra.Command {
	opts := Options{}

	cmd := &cobra.Command{
		Use:   "cancel",
		Short: "Cancel a scheduled email",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCancel(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.apiKey, "api-key", "", "Z-Send API key (or set ZSEND_API_KEY)")
	cmd.Flags().StringVar(&opts.id, "id", "", "Scheduled email ID")

	return cmd
}

func runCancel(f *cmdutil.Factory, opts Options) error {
	if f.Interactive {
		return runCancelInteractive(f, opts)
	}
	return runCancelNonInteractive(f, opts)
}

func runCancelInteractive(f *cmdutil.Factory, opts Options) error {
	if opts.id == "" {
		id, err := f.Prompter.Input("Scheduled Email ID: ", "")
		if err != nil {
			return err
		}
		opts.id = id
	}
	if err := paramCheck(opts); err != nil {
		return err
	}

	confirmed, err := f.Prompter.Confirm("Are you sure you want to cancel this scheduled email?", false)
	if err != nil {
		return err
	}
	if !confirmed {
		return nil
	}

	return doCancel(f, opts)
}

func runCancelNonInteractive(f *cmdutil.Factory, opts Options) error {
	if err := paramCheck(opts); err != nil {
		return err
	}
	return doCancel(f, opts)
}

func doCancel(f *cmdutil.Factory, opts Options) error {
	if opts.apiKey == "" {
		opts.apiKey = os.Getenv("ZSEND_API_KEY")
	}
	if opts.apiKey == "" {
		return fmt.Errorf("Z-Send API key is required (--api-key or ZSEND_API_KEY)")
	}
	if !strings.HasPrefix(opts.apiKey, "zs_") {
		return fmt.Errorf("invalid API key format: must start with zs_")
	}

	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Cancelling scheduled email..."),
	)
	s.Start()
	err := f.ApiClient.CancelZSendScheduledEmail(context.Background(), opts.apiKey, opts.id)
	s.Stop()
	if err != nil {
		return err
	}

	f.Log.Infof("Scheduled email %s cancelled successfully.", opts.id)
	return nil
}

func paramCheck(opts Options) error {
	if strings.TrimSpace(opts.id) == "" {
		return fmt.Errorf("scheduled email ID is required")
	}
	return nil
}
