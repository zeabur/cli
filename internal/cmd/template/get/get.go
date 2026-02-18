package get

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct {
	code string
	raw  bool
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
	cmd.Flags().BoolVar(&opts.raw, "raw", false, "Output raw YAML spec")

	return cmd
}

func runGet(f *cmdutil.Factory, opts Options) error {
	if f.Interactive {
		return runGetInteractive(f, opts)
	}
	return runGetNonInteractive(f, opts)
}

func runGetInteractive(f *cmdutil.Factory, opts Options) error {
	if opts.code == "" {
		code, err := f.Prompter.Input("Template Code: ", "")
		if err != nil {
			return err
		}
		opts.code = code
	}

	return getTemplate(f, opts)
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
	if opts.raw {
		return getTemplateRaw(opts.code)
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
	s.Stop()

	if template == nil || template.Code == "" {
		fmt.Println("Template not found")
	} else {
		f.Printer.Table([]string{"Code", "Name", "Description"}, [][]string{{template.Code, template.Name, template.Description}})
	}

	return nil
}

func getTemplateRaw(code string) error {
	resp, err := http.Get("https://zeabur.com/templates/" + code + ".yaml")
	if err != nil {
		return fmt.Errorf("failed to fetch template YAML: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("template not found (HTTP %d)", resp.StatusCode)
	}

	_, err = io.Copy(os.Stdout, resp.Body)
	return err
}

func paramCheck(opts Options) error {
	if opts.code == "" {
		return fmt.Errorf("template code is required")
	}

	return nil
}
