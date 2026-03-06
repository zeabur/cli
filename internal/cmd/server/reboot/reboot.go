package reboot

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct {
	id          string
	skipConfirm bool
}

func NewCmdReboot(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "reboot [server-id]",
		Short: "Reboot a dedicated server",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 && opts.id == "" {
				opts.id = args[0]
			}
			return runReboot(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.id, "id", "", "Server ID")
	cmd.Flags().BoolVarP(&opts.skipConfirm, "yes", "y", false, "Skip confirmation")

	return cmd
}

func runReboot(f *cmdutil.Factory, opts *Options) error {
	if f.Interactive {
		return runRebootInteractive(f, opts)
	}
	return runRebootNonInteractive(f, opts)
}

func runRebootInteractive(f *cmdutil.Factory, opts *Options) error {
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

		idx, err := f.Prompter.Select("Select a server to reboot", "", options)
		if err != nil {
			return err
		}
		opts.id = servers[idx].ID
	}

	return runRebootNonInteractive(f, opts)
}

func runRebootNonInteractive(f *cmdutil.Factory, opts *Options) error {
	if opts.id == "" {
		return fmt.Errorf("--id is required")
	}

	if f.Interactive && !opts.skipConfirm {
		confirm, err := f.Prompter.Confirm(fmt.Sprintf("Are you sure to reboot server <%s>?", opts.id), false)
		if err != nil {
			return err
		}
		if !confirm {
			return nil
		}
	}

	err := f.ApiClient.RebootServer(context.Background(), opts.id)
	if err != nil {
		return fmt.Errorf("reboot server failed: %w", err)
	}

	f.Log.Infof("Server <%s> rebooted successfully", opts.id)

	return nil
}
