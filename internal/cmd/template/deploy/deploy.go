package deploy

import (
	"context"
	"fmt"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct {
	code string
}

func NewCmdDeploy(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy a template",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDeploy(f, opts)
		},
	}

	cmd.Flags().StringVarP(&opts.code, "code", "c", "", "Template code")

	return cmd
}

func runDeploy(f *cmdutil.Factory, opts *Options) error {
	var err error

	err = paramCheck(opts)
	if err != nil {
		return err
	}

	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Fetching template..."),
	)
	s.Start()

	template, err := f.ApiClient.GetTemplate(context.Background(), opts.code)
	if err != nil {
		return err
	}

	fmt.Printf("Template: %s\n", template.Name)

	s.Stop()

	return nil
}

func paramCheck(opts *Options) error {
	if opts.code == "" {
		return fmt.Errorf("code is required")
	}

	return nil
}
