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
		Use:   "get <id>",
		Short: "Get details of a scheduled email",
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
		spinner.WithSuffix(" Fetching scheduled email..."),
	)
	s.Start()
	email, err := f.ApiClient.GetZSendScheduledEmail(context.Background(), apiKey, id)
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
