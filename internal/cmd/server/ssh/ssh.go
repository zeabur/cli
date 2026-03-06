package ssh

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct {
	id string
}

func NewCmdSSH(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "ssh [server-id]",
		Short: "SSH into a dedicated server",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				if opts.id != "" && opts.id != args[0] {
					return fmt.Errorf("conflicting server IDs: arg=%q, --id=%q", args[0], opts.id)
				}
				opts.id = args[0]
			}
			return runSSH(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.id, "id", "", "Server ID")

	return cmd
}

func runSSH(f *cmdutil.Factory, opts *Options) error {
	if f.Interactive {
		return runSSHInteractive(f, opts)
	}
	return runSSHNonInteractive(f, opts)
}

func runSSHInteractive(f *cmdutil.Factory, opts *Options) error {
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

	return runSSHNonInteractive(f, opts)
}

func runSSHNonInteractive(f *cmdutil.Factory, opts *Options) error {
	if opts.id == "" {
		return fmt.Errorf("--id is required")
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

	// Try to get password for managed servers
	var password string
	if server.IsManaged {
		pw, err := f.ApiClient.RevealServerPassword(ctx, opts.id)
		if err == nil && pw != "" {
			password = pw
		}
	}

	sshArgs := []string{
		fmt.Sprintf("%s@%s", username, server.IP),
		"-p", strconv.Itoa(server.SSHPort),
		"-o", "StrictHostKeyChecking=no",
	}

	f.Log.Infof("Connecting to %s@%s:%d ...", username, server.IP, server.SSHPort)

	if password != "" {
		sshpassPath, err := exec.LookPath("sshpass")
		if err == nil {
			// Use sshpass for automatic password authentication
			sshCmd := exec.Command(sshpassPath, append([]string{"-p", password, "ssh"}, sshArgs...)...)
			sshCmd.Stdin = os.Stdin
			sshCmd.Stdout = os.Stdout
			sshCmd.Stderr = os.Stderr
			return sshCmd.Run()
		}
		// sshpass not available, show password to user
		f.Log.Infof("Password: %s", password)
		f.Log.Infof("(Install sshpass for automatic login)")
	}

	sshCmd := exec.Command("ssh", sshArgs...)
	sshCmd.Stdin = os.Stdin
	sshCmd.Stdout = os.Stdout
	sshCmd.Stderr = os.Stderr

	return sshCmd.Run()
}
