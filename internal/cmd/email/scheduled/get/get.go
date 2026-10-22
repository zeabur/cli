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
		Short: "Get details of a scheduled email",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGet(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.apiKey, "api-key", "", cmdutil.ZSendAPIKeyFlagUsage)
	cmd.Flags().StringVar(&opts.id, "id", "", "Scheduled email ID")

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
		id, err := f.Prompter.Input("Scheduled Email ID: ", "")
		if err != nil {
			return err
		}
		opts.id = id
	}
	if err := paramCheck(opts); err != nil {
		return err
	}
	return getScheduledEmail(f, opts)
}

func runGetNonInteractive(f *cmdutil.Factory, opts Options) error {
	if err := paramCheck(opts); err != nil {
		return err
	}
	return getScheduledEmail(f, opts)
}

func getScheduledEmail(f *cmdutil.Factory, opts Options) error {
	apiKey, err := f.ResolveZSendAPIKey(opts.apiKey)
	if err != nil {
		return err
	}
	opts.apiKey = apiKey

	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Fetching scheduled email..."),
	)
	s.Start()
	email, err := f.ApiClient.GetZSendScheduledEmail(context.Background(), opts.apiKey, opts.id)
	s.Stop()
	if err != nil {
		return err
	}

	if f.JSON {
		return f.Printer.JSON(email)
	}

	f.Printer.Table(
		[]string{"Field", "Value"},
		[][]string{
			{"ID", email.ID},
			{"From", email.From},
			{"To", fmt.Sprintf("%v", email.To)},
			{"Subject", email.Subject},
			{"Status", email.Status},
			{"Scheduled At", email.ScheduledAt},
			{"Sent At", email.SentAt},
			{"Attempts", fmt.Sprintf("%d", email.Attempts)},
			{"Last Error", email.LastError},
			{"Created At", email.CreatedAt},
		},
	)
	return nil
}

func paramCheck(opts Options) error {
	if strings.TrimSpace(opts.id) == "" {
		return fmt.Errorf("scheduled email ID is required")
	}
	return nil
}
