package rename

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct {
	id   string
	name string
}

func NewCmdRename(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "rename [server-id]",
		Short: "Rename a dedicated server",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				if opts.id != "" && opts.id != args[0] {
					return fmt.Errorf("conflicting server IDs: arg=%q, --id=%q", args[0], opts.id)
				}
				opts.id = args[0]
			}
			return runRename(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.id, "id", "", "Server ID")
	cmd.Flags().StringVar(&opts.name, "name", "", "New server name")

	return cmd
}

func runRename(f *cmdutil.Factory, opts *Options) error {
	if f.Interactive {
		return runRenameInteractive(f, opts)
	}
	return runRenameNonInteractive(f, opts)
}

func runRenameInteractive(f *cmdutil.Factory, opts *Options) error {
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

		idx, err := f.Prompter.Select("Select a server to rename", "", options)
		if err != nil {
			return err
		}
		opts.id = servers[idx].ID
	}

	if strings.TrimSpace(opts.name) == "" {
		name, err := f.Prompter.Input("New server name", "")
		if err != nil {
			return err
		}
		opts.name = name
	}

	return runRenameNonInteractive(f, opts)
}

func runRenameNonInteractive(f *cmdutil.Factory, opts *Options) error {
	if opts.id == "" {
		return fmt.Errorf("--id is required")
	}
	opts.name = strings.TrimSpace(opts.name)
	if opts.name == "" {
		return fmt.Errorf("--name is required")
	}

	if err := f.ApiClient.RenameServer(context.Background(), opts.id, opts.name); err != nil {
		return fmt.Errorf("rename server failed: %w", err)
	}

	f.Log.Infof("Server <%s> renamed to %q", opts.id, opts.name)
	return nil
}
