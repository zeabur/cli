package create

import (
	"context"
	"fmt"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/pkg/model"
)

var permissionChoices = []string{
	"all",
	"send_only",
	"read_only",
}

type Options struct {
	name       string
	permission string
}

func NewCmdCreate(f *cmdutil.Factory) *cobra.Command {
	opts := Options{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an email API key",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreate(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.name, "name", "", "API key name")
	cmd.Flags().StringVar(&opts.permission, "permission", "", "Permission (all, send_only, read_only)")

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

	if opts.permission == "" {
		idx, err := f.Prompter.Select("Permission: ", "", permissionChoices)
		if err != nil {
			return err
		}
		opts.permission = permissionChoices[idx]
	}

	return createKey(f, opts)
}

func runCreateNonInteractive(f *cmdutil.Factory, opts Options) error {
	if err := paramCheck(opts); err != nil {
		return err
	}
	return createKey(f, opts)
}

func createKey(f *cmdutil.Factory, opts Options) error {
	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Creating API key..."),
	)
	s.Start()
	reply, err := f.ApiClient.CreateZSendAPIKey(context.Background(), model.CreateZSendAPIKeyInput{
		Name:       opts.name,
		Permission: opts.permission,
	})
	s.Stop()
	if err != nil {
		return err
	}

	if f.JSON {
		return f.Printer.JSON(reply)
	}

	f.Log.Infof("API key %q created successfully (ID: %s)", reply.APIKey.Name, reply.APIKey.ID)
	if reply.APIKey.Token != nil {
		f.Log.Infof("Token: %s", *reply.APIKey.Token)
		f.Log.Infof("WARNING: This token will only be shown once. Please save it now.")
	}

	return nil
}

func paramCheck(opts Options) error {
	if opts.name == "" {
		return fmt.Errorf("name is required")
	}
	if opts.permission == "" {
		return fmt.Errorf("permission is required")
	}
	return nil
}
