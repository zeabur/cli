package get

import (
	"context"
	"fmt"
	"strings"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct {
	id string
}

func NewCmdGet(f *cmdutil.Factory) *cobra.Command {
	opts := Options{}

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get details of an API key",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGet(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.id, "id", "", "API key ID")

	return cmd
}

func runGet(f *cmdutil.Factory, opts Options) error {
	if f.Interactive {
		return runGetInteractive(f, opts)
	}
	return runGetNonInteractive(f, opts)
}

func runGetInteractive(f *cmdutil.Factory, opts Options) error {
	if opts.id == "" {
		id, err := f.Prompter.Input("API Key ID: ", "")
		if err != nil {
			return err
		}
		opts.id = id
	}
	if err := paramCheck(opts); err != nil {
		return err
	}
	return getKey(f, opts)
}

func runGetNonInteractive(f *cmdutil.Factory, opts Options) error {
	if err := paramCheck(opts); err != nil {
		return err
	}
	return getKey(f, opts)
}

func getKey(f *cmdutil.Factory, opts Options) error {
	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Fetching API key..."),
	)
	s.Start()
	key, err := f.ApiClient.GetZSendAPIKey(context.Background(), opts.id)
	s.Stop()
	if err != nil {
		return err
	}

	if f.JSON {
		return f.Printer.JSON(key)
	}

	domains := strings.Join(key.Domains, ", ")

	f.Printer.Table(
		[]string{"Field", "Value"},
		[][]string{
			{"ID", key.ID},
			{"Name", key.Name},
			{"Permission", key.Permission},
			{"Domains", domains},
			{"Created At", key.CreatedAt.String()},
		},
	)
	return nil
}

func paramCheck(opts Options) error {
	if strings.TrimSpace(opts.id) == "" {
		return fmt.Errorf("API key ID is required")
	}
	return nil
}
