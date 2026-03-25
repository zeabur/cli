package get

import (
	"context"
	"fmt"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
)

func NewCmdGet(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get details of an email record",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGet(f, args[0])
		},
	}
	return cmd
}

func runGet(f *cmdutil.Factory, id string) error {
	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Fetching email..."),
	)
	s.Start()
	email, err := f.ApiClient.GetZSendEmail(context.Background(), id)
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
