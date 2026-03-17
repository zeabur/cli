package status

import (
	"context"
	"fmt"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct{}

func NewCmdStatus(f *cmdutil.Factory) *cobra.Command {
	opts := Options{}

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show Zeabur Email user status",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStatus(f, opts)
		},
	}

	return cmd
}

func runStatus(f *cmdutil.Factory, opts Options) error {
	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Fetching email status..."),
	)
	s.Start()
	status, err := f.ApiClient.GetZSendUserStatus(context.Background())
	s.Stop()
	if err != nil {
		return err
	}

	if f.JSON {
		return f.Printer.JSON(status)
	}

	statusMsg := ""
	if status.StatusMsg != nil {
		statusMsg = *status.StatusMsg
	}

	f.Printer.Table(
		[]string{"Field", "Value"},
		[][]string{
			{"Status", status.Status},
			{"Status Message", statusMsg},
			{"Daily Quota", fmt.Sprintf("%d", status.DailyQuota)},
			{"Daily Sent", fmt.Sprintf("%d", status.DailySent)},
			{"Monthly Quota", fmt.Sprintf("%d", status.MonthlyQuota)},
			{"Monthly Sent", fmt.Sprintf("%d", status.MonthlySent)},
			{"Max Domains", fmt.Sprintf("%d", status.MaxDomains)},
			{"Max API Keys", fmt.Sprintf("%d", status.MaxAPIKeys)},
			{"Max Webhooks", fmt.Sprintf("%d", status.MaxWebhooks)},
			{"Sent (24h)", fmt.Sprintf("%d", status.SentCount24h)},
			{"Bounces (24h)", fmt.Sprintf("%d", status.BounceCount24h)},
			{"Complaints (24h)", fmt.Sprintf("%d", status.ComplaintCount24h)},
		},
	)

	return nil
}
