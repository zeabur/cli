package verify

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

func NewCmdVerify(f *cmdutil.Factory) *cobra.Command {
	opts := Options{}

	cmd := &cobra.Command{
		Use:   "verify",
		Short: "Verify an email webhook",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runVerify(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.id, "id", "", "Webhook ID")

	return cmd
}

func runVerify(f *cmdutil.Factory, opts Options) error {
	if f.Interactive {
		return runVerifyInteractive(f, opts)
	}
	return runVerifyNonInteractive(f, opts)
}

func runVerifyInteractive(f *cmdutil.Factory, opts Options) error {
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
	return verifyWebhook(f, opts)
}

func runVerifyNonInteractive(f *cmdutil.Factory, opts Options) error {
	if err := paramCheck(opts); err != nil {
		return err
	}
	return verifyWebhook(f, opts)
}

func verifyWebhook(f *cmdutil.Factory, opts Options) error {
	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Verifying webhook..."),
	)
	s.Start()
	reply, err := f.ApiClient.VerifyZSendWebhook(context.Background(), opts.id)
	s.Stop()
	if err != nil {
		return err
	}

	if f.JSON {
		if err := f.Printer.JSON(reply); err != nil {
			return err
		}
	} else if reply.Success {
		f.Log.Infof("Webhook verification successful: %s", reply.Message)
	}

	if !reply.Success {
		return fmt.Errorf("webhook verification failed: %s", reply.Message)
	}

	return nil
}

func paramCheck(opts Options) error {
	if opts.id == "" {
		return fmt.Errorf("webhook ID is required")
	}
	return nil
}
