package get

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct {
	id string
}

func NewCmdGet(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "get [server-id]",
		Short: "Get a dedicated server",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				if opts.id != "" && opts.id != args[0] {
					return fmt.Errorf("conflicting server IDs: arg=%q, --id=%q", args[0], opts.id)
				}
				opts.id = args[0]
			}
			return runGet(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.id, "id", "", "Server ID")

	return cmd
}

func runGet(f *cmdutil.Factory, opts *Options) error {
	if f.Interactive {
		return runGetInteractive(f, opts)
	}
	return runGetNonInteractive(f, opts)
}

func runGetInteractive(f *cmdutil.Factory, opts *Options) error {
	if opts.id == "" {
		servers, err := f.ApiClient.ListServers(context.Background())
		if err != nil {
			return fmt.Errorf("list servers failed: %w", err)
		}
		if len(servers) == 0 {
			return fmt.Errorf("no servers found")
		}

		options := make([]string, len(servers))
		for i, s := range servers {
			location := s.IP
			if s.City != nil && s.Country != nil {
				location = fmt.Sprintf("%s, %s", *s.City, *s.Country)
			} else if s.Country != nil {
				location = *s.Country
			}
			options[i] = fmt.Sprintf("%s (%s)", s.Name, location)
		}

		idx, err := f.Prompter.Select("Select a server", "", options)
		if err != nil {
			return err
		}
		opts.id = servers[idx].ID
	}

	return runGetNonInteractive(f, opts)
}

func runGetNonInteractive(f *cmdutil.Factory, opts *Options) error {
	if opts.id == "" {
		return fmt.Errorf("--id is required")
	}

	server, err := f.ApiClient.GetServer(context.Background(), opts.id)
	if err != nil {
		return fmt.Errorf("get server failed: %w", err)
	}

	f.Printer.Table(server.Header(), server.Rows())

	return nil
}
