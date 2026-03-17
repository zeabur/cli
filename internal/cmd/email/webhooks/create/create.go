package create

import (
	"context"
	"fmt"
	"strings"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/pkg/model"
)

var eventChoices = []string{
	"send",
	"delivery",
	"bounce",
	"complaint",
	"reject",
	"open",
	"click",
}

type Options struct {
	name     string
	endpoint string
	events   string
}

func NewCmdCreate(f *cmdutil.Factory) *cobra.Command {
	opts := Options{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an email webhook",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreate(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.name, "name", "", "Webhook name")
	cmd.Flags().StringVar(&opts.endpoint, "endpoint", "", "Webhook endpoint URL")
	cmd.Flags().StringVar(&opts.events, "events", "", "Comma-separated events (send,delivery,bounce,complaint,reject,open,click)")

	return cmd
}

func runCreate(f *cmdutil.Factory, opts Options) error {
	if f.Interactive {
		return runCreateInteractive(f, opts)
	}
	return runCreateNonInteractive(f, opts)
}

func runCreateInteractive(f *cmdutil.Factory, opts Options) error {
	if opts.name == "" {
		name, err := f.Prompter.Input("Name: ", "")
		if err != nil {
			return err
		}
		opts.name = name
	}

	if opts.endpoint == "" {
		endpoint, err := f.Prompter.Input("Endpoint URL: ", "")
		if err != nil {
			return err
		}
		opts.endpoint = endpoint
	}

	if opts.events == "" {
		events, err := f.Prompter.Input("Events (comma-separated, e.g. send,delivery,bounce): ", "")
		if err != nil {
			return err
		}
		opts.events = events
	}

	return createWebhook(f, opts)
}

func runCreateNonInteractive(f *cmdutil.Factory, opts Options) error {
	if err := paramCheck(opts); err != nil {
		return err
	}
	return createWebhook(f, opts)
}

func createWebhook(f *cmdutil.Factory, opts Options) error {
	var events []string
	if opts.events != "" {
		for _, e := range strings.Split(opts.events, ",") {
			e = strings.TrimSpace(e)
			if e != "" {
				events = append(events, e)
			}
		}
	}

	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Creating webhook..."),
	)
	s.Start()
	reply, err := f.ApiClient.CreateZSendWebhook(context.Background(), model.CreateZSendWebhookInput{
		Name:     opts.name,
		Endpoint: opts.endpoint,
		Events:   events,
	})
	s.Stop()
	if err != nil {
		return err
	}

	if f.JSON {
		return f.Printer.JSON(reply)
	}

	f.Log.Infof("Webhook %q created successfully (ID: %s)", reply.Webhook.Name, reply.Webhook.ID)
	f.Log.Infof("Secret: %s", reply.Secret)
	f.Log.Infof("WARNING: This secret will only be shown once. Please save it now.")

	return nil
}

func paramCheck(opts Options) error {
	if opts.name == "" {
		return fmt.Errorf("name is required")
	}
	if opts.endpoint == "" {
		return fmt.Errorf("endpoint is required")
	}
	return nil
}
