package get

import (
	"context"
	"fmt"
	"strings"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
)

func NewCmdGet(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get details of a webhook",
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
		spinner.WithSuffix(" Fetching webhook..."),
	)
	s.Start()
	wh, err := f.ApiClient.GetZSendWebhook(context.Background(), id)
	s.Stop()
	if err != nil {
		return err
	}

	if f.JSON {
		return f.Printer.JSON(wh)
	}

	enabled := "No"
	if wh.Enabled {
		enabled = "Yes"
	}
	statusMsg := ""
	if wh.StatusMsg != nil {
		statusMsg = *wh.StatusMsg
	}

	f.Printer.Table(
		[]string{"Field", "Value"},
		[][]string{
			{"ID", wh.ID},
			{"Name", wh.Name},
			{"Endpoint", wh.Endpoint},
			{"Events", strings.Join(wh.Events, ", ")},
			{"Status", wh.Status},
			{"Status Msg", statusMsg},
			{"Enabled", enabled},
			{"Total Sent", fmt.Sprintf("%d", wh.TotalSent)},
			{"Success", fmt.Sprintf("%d", wh.SuccessCount)},
			{"Failed", fmt.Sprintf("%d", wh.FailureCount)},
			{"Created At", wh.CreatedAt.String()},
		},
	)
	return nil
}
