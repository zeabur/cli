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
	id string
}

func NewCmdGet(f *cmdutil.Factory) *cobra.Command {
	opts := Options{}

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get details of a webhook",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGet(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.id, "id", "", "Webhook ID")

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
		id, err := f.Prompter.Input("Webhook ID: ", "")
		if err != nil {
			return err
		}
		opts.id = id
	}
	if err := paramCheck(opts); err != nil {
		return err
	}
	return getWebhook(f, opts)
}

func runGetNonInteractive(f *cmdutil.Factory, opts Options) error {
	if err := paramCheck(opts); err != nil {
		return err
	}
	return getWebhook(f, opts)
}

func getWebhook(f *cmdutil.Factory, opts Options) error {
	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Fetching webhook..."),
	)
	s.Start()
	wh, err := f.ApiClient.GetZSendWebhook(context.Background(), opts.id)
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

func paramCheck(opts Options) error {
	if opts.id == "" {
		return fmt.Errorf("webhook ID is required")
	}
	return nil
}
