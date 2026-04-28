package search

import (
	"context"
	"fmt"
	"strings"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct {
	query string
}

func NewCmdSearch(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "search [query]",
		Short: "Search for available domains",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.query = args[0]
			}
			return runSearch(f, opts)
		},
	}

	return cmd
}

func runSearch(f *cmdutil.Factory, opts *Options) error {
	if opts.query == "" {
		if !f.Interactive {
			return fmt.Errorf("query argument is required")
		}
		query, err := f.Prompter.Input("Search for a domain: ", "")
		if err != nil {
			return err
		}
		opts.query = query
	}

	if !strings.ContainsRune(opts.query, '.') {
		return fmt.Errorf("please provide a full domain name with TLD (e.g., example.com, example.io)")
	}

	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Checking availability..."),
	)
	s.Start()
	result, err := f.ApiClient.CheckDomainRegistrationAvailability(context.Background(), opts.query)
	s.Stop()
	if err != nil {
		return fmt.Errorf("check domain availability failed: %w", err)
	}

	if f.JSON {
		return f.Printer.JSON(result)
	}

	f.Printer.Table(result.Header(), result.Rows())
	return nil
}
