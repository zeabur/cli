package get

import (
	"context"
	"fmt"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct {
	id string
}

func NewCmdGet(f *cmdutil.Factory) *cobra.Command {
	opts := Options{}

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get details of an email record",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGet(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.id, "id", "", "Email ID")

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
		id, err := f.Prompter.Input("Email ID: ", "")
		if err != nil {
			return err
		}
		opts.id = id
	}
	if err := paramCheck(opts); err != nil {
		return err
	}
	return getEmail(f, opts)
}

func runGetNonInteractive(f *cmdutil.Factory, opts Options) error {
	if err := paramCheck(opts); err != nil {
		return err
	}
	return getEmail(f, opts)
}

func getEmail(f *cmdutil.Factory, opts Options) error {
	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Fetching email..."),
	)
	s.Start()
	email, err := f.ApiClient.GetZSendEmail(context.Background(), opts.id)
	s.Stop()
	if err != nil {
		return err
	}

	if f.JSON {
		return f.Printer.JSON(email)
	}

	scheduledAt := ""
	if email.ScheduledAt != nil {
		scheduledAt = email.ScheduledAt.String()
	}

	f.Printer.Table(
		[]string{"Field", "Value"},
		[][]string{
			{"ID", email.ID},
			{"From", email.From},
			{"To", fmt.Sprintf("%v", email.To)},
			{"Subject", email.Subject},
			{"Status", email.Status},
			{"Job Type", email.JobType},
			{"Job ID", email.JobID},
			{"Scheduled At", scheduledAt},
			{"Created At", email.CreatedAt.String()},
		},
	)
	return nil
}

func paramCheck(opts Options) error {
	if opts.id == "" {
		return fmt.Errorf("email ID is required")
	}
	return nil
}
