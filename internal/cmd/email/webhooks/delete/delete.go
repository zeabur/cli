package delete

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

func NewCmdDelete(f *cmdutil.Factory) *cobra.Command {
	opts := Options{}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete an email webhook",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDelete(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.id, "id", "", "Webhook ID")

	return cmd
}

func runDelete(f *cmdutil.Factory, opts Options) error {
	if f.Interactive {
		return runDeleteInteractive(f, opts)
	}
	return runDeleteNonInteractive(f, opts)
}

func runDeleteInteractive(f *cmdutil.Factory, opts Options) error {
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

	confirmed, err := f.Prompter.Confirm("Are you sure you want to delete this webhook?", false)
	if err != nil {
		return err
	}
	if !confirmed {
		return nil
	}

	return deleteWebhook(f, opts)
}

func runDeleteNonInteractive(f *cmdutil.Factory, opts Options) error {
	if err := paramCheck(opts); err != nil {
		return err
	}
	return deleteWebhook(f, opts)
}

func deleteWebhook(f *cmdutil.Factory, opts Options) error {
	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Deleting webhook..."),
	)
	s.Start()
	err := f.ApiClient.DeleteZSendWebhook(context.Background(), opts.id)
	s.Stop()
	if err != nil {
		return err
	}

	f.Log.Infof("Webhook deleted successfully")

	return nil
}

func paramCheck(opts Options) error {
	if opts.id == "" {
		return fmt.Errorf("webhook ID is required")
	}
	return nil
}
