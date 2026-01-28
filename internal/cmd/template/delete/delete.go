package delete

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

func NewCmdDelete(f *cmdutil.Factory) *cobra.Command {
	opts := Options{}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete template by code",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDelete(f, opts)
		},
	}

	cmd.Flags().StringVarP(&opts.code, "code", "c", "", "Template code")

	return cmd
}

func runDelete(f *cmdutil.Factory, opts Options) error {
	if f.Interactive {
		return runDeleteInteractive(f, opts)
	}
	return runDeleteNonInteractive(f, opts)
}

func runDeleteInteractive(f *cmdutil.Factory, opts Options) error {
	code, err := f.Prompter.Input("Template Code: ", "")
	if err != nil {
		return err
	}

	opts.code = code

	err = deleteTemplate(f, opts)
	if err != nil {
		return err
	}

	return nil
}

func runDeleteNonInteractive(f *cmdutil.Factory, opts Options) error {
	err := paramCheck(opts)
	if err != nil {
		return err
	}

	err = deleteTemplate(f, opts)
	if err != nil {
		return err
	}

	return nil
}

func deleteTemplate(f *cmdutil.Factory, opts Options) error {
	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Deleting template..."),
	)
	s.Start()
	err := f.ApiClient.DeleteTemplate(context.Background(), opts.code)
	if err != nil {
		return err
	}
	s.Stop()

	f.Log.Info("Delete template successfully")

	return nil
}

func paramCheck(opts Options) error {
	if opts.code == "" {
		return fmt.Errorf("template code is required")
	}

	return nil
}
