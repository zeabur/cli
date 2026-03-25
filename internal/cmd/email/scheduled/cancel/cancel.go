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

func NewCmdCancel(f *cmdutil.Factory) *cobra.Command {
	var apiKey string

	cmd := &cobra.Command{
		Use:   "cancel <id>",
		Short: "Cancel a scheduled email",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCancel(f, apiKey, args[0])
		},
	}

	cmd.Flags().StringVar(&apiKey, "api-key", "", "Z-Send API key (or set ZSEND_API_KEY)")

	return cmd
}

func runCancel(f *cmdutil.Factory, apiKey, id string) error {
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
		spinner.WithSuffix(" Cancelling scheduled email..."),
	)
	s.Start()
	err := f.ApiClient.CancelZSendScheduledEmail(context.Background(), apiKey, id)
	s.Stop()
	if err != nil {
		return err
	}

	fmt.Printf("Scheduled email %s cancelled successfully.\n", id)
	return nil
}
