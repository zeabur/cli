package get

import (
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
	err := paramCheck(opts)
	if err != nil {
		return err
	}

	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Fetching template..."),
	)
	s.Start()
  template, err := f.ApiClient.GetTemplate(ctx context.Context, code string)

	return nil
}

func paramCheck(opts Options) error {
	if opts.code == "" {
		return fmt.Errorf("template code is required")
	}

	return nil
}
