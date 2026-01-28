package get

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

func NewCmdGet(f *cmdutil.Factory) *cobra.Command {
	opts := Options{}

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get template by code",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGet(f, opts)
		},
	}

	cmd.Flags().StringVarP(&opts.code, "code", "c", "", "Template code")

	return cmd
}

func runGet(f *cmdutil.Factory, opts Options) error {
	if f.Interactive {
		return runGetInteractive(f, opts)
	}
	return runGetNonInteractive(f, opts)
}

func runGetInteractive(f *cmdutil.Factory, opts Options) error {
	code, err := f.Prompter.Input("Template Code: ", "")
	if err != nil {
		return err
	}

	opts.code = code

	err = getTemplate(f, opts)
	if err != nil {
		return err
	}

	return nil
}

func runGetNonInteractive(f *cmdutil.Factory, opts Options) error {
	err := paramCheck(opts)
	if err != nil {
		return err
	}

	err = getTemplate(f, opts)
	if err != nil {
		return err
	}

	return nil
}

func getTemplate(f *cmdutil.Factory, opts Options) error {
	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Fetching template..."),
	)
	s.Start()
	template, err := f.ApiClient.GetTemplate(context.Background(), opts.code)
	if err != nil {
		return err
	}
	s.Stop()

	if template == nil || template.Code == "" {
		fmt.Println("Template not found")
	} else {
		f.Printer.Table([]string{"Code", "Name", "Description"}, [][]string{{template.Code, template.Name, template.Description}})
	}

	return nil
}

func paramCheck(opts Options) error {
	if opts.code == "" {
		return fmt.Errorf("template code is required")
	}

	return nil
}
