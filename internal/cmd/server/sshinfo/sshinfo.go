package sshinfo

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct {
	id string
}

type SSHInfo struct {
	IP       string `json:"ip"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewCmdSSHInfo(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "ssh-info [server-id]",
		Short: "Get SSH connection info for a dedicated server",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				if opts.id != "" && opts.id != args[0] {
					return fmt.Errorf("conflicting server IDs: arg=%q, --id=%q", args[0], opts.id)
				}
				opts.id = args[0]
			}
			return runSSHInfo(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.id, "id", "", "Server ID")

	return cmd
}

func runSSHInfo(f *cmdutil.Factory, opts *Options) error {
	if f.Interactive && opts.id == "" {
		return runSSHInfoInteractive(f, opts)
	}
	return runSSHInfoNonInteractive(f, opts)
}

func runSSHInfoInteractive(f *cmdutil.Factory, opts *Options) error {
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

	return runSSHInfoNonInteractive(f, opts)
}

func runSSHInfoNonInteractive(f *cmdutil.Factory, opts *Options) error {
	if opts.id == "" {
		return fmt.Errorf("server ID is required (use --id or positional server-id)")
	}

	ctx := context.Background()

	server, err := f.ApiClient.GetServer(ctx, opts.id)
	if err != nil {
		return fmt.Errorf("get server failed: %w", err)
	}

	username := "root"
	if server.SSHUsername != nil && *server.SSHUsername != "" {
		username = *server.SSHUsername
	}

	var password string
	if server.IsManaged {
		pw, err := f.ApiClient.RevealServerPassword(ctx, opts.id)
		if err != nil {
			return fmt.Errorf("reveal server password failed: %w", err)
		}
		password = pw
	}

	info := SSHInfo{
		IP:       server.IP,
		Port:     server.SSHPort,
		Username: username,
		Password: password,
	}

	data, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("marshal failed: %w", err)
	}

	if _, err := fmt.Println(string(data)); err != nil { //nolint:gosec // Intentionally outputting credentials for agent consumption
		return fmt.Errorf("write output failed: %w", err)
	}
	return nil
}
