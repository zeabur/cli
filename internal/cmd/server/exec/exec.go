// Package exec implements `zeabur server exec`: run a command on a dedicated
// server over SSH and stream its output.
//
// Unlike the two-step `server ssh-info` + hand-rolled ssh2/sshpass flow (which
// embeds the password into a multi-layer-escaped shell script and corrupts
// passwords containing special characters), this fetches the credentials and
// opens the connection entirely in-process via golang.org/x/crypto/ssh. The
// password is never printed nor passed through a shell, so special characters
// are safe and no external ssh/sshpass client is required (DES-802).
package exec

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"

	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct {
	id string
}

func NewCmdExec(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "exec [server-id] -- <command> [args...]",
		Short: "Run a command on a dedicated server over SSH",
		Long: `Run a command on a dedicated server over SSH and stream its output.

Credentials are fetched and used inside the CLI — never printed and never passed
through a shell — so passwords with special characters work and no external ssh
or sshpass client is required.

The command after '--' is joined with spaces and run by the remote login shell,
matching 'ssh host <command>' semantics. Quote a compound command as a single
argument to keep it intact.`,
		Example: `  zeabur server exec --id <server-id> -- uname -a
  zeabur server exec <server-id> -- sudo kubectl get pods -A
  zeabur server exec <server-id> -- 'echo start && sudo systemctl status k3s'`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Everything after '--' is the remote command; anything before it is
			// the optional positional server id. cobra reports the '--' position.
			dash := cmd.ArgsLenAtDash()
			if dash < 0 {
				return fmt.Errorf("missing command; use: zeabur server exec [server-id] -- <command>")
			}

			positional := args[:dash]
			remoteArgs := args[dash:]
			if len(remoteArgs) == 0 {
				return fmt.Errorf("no command given after '--'")
			}
			if len(positional) > 1 {
				return fmt.Errorf("unexpected arguments before '--': %v", positional[1:])
			}
			if len(positional) == 1 {
				if opts.id != "" && opts.id != positional[0] {
					return fmt.Errorf("conflicting server IDs: arg=%q, --id=%q", positional[0], opts.id)
				}
				opts.id = positional[0]
			}

			return runExec(f, opts, strings.Join(remoteArgs, " "))
		},
	}

	cmd.Flags().StringVar(&opts.id, "id", "", "Server ID")

	return cmd
}

func runExec(f *cmdutil.Factory, opts *Options, remoteCmd string) error {
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

	var auth []ssh.AuthMethod
	if server.IsManaged {
		pw, err := f.ApiClient.RevealServerPassword(ctx, opts.id)
		if err == nil && pw != "" {
			auth = append(auth, ssh.Password(pw))
		}
	}
	if len(auth) == 0 {
		return fmt.Errorf(
			"no password available for server %s; `server exec` uses password auth for managed servers — for key-based access use `zeabur server ssh`",
			opts.id,
		)
	}

	config := &ssh.ClientConfig{
		User: username,
		Auth: auth,
		// Equivalent to `-o StrictHostKeyChecking=no`: dedicated servers are
		// reached by IP with rotating host keys, so pinning would only wedge.
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         15 * time.Second,
	}
	addr := fmt.Sprintf("%s:%d", server.IP, server.SSHPort)

	f.Log.Infof("Connecting to %s@%s ...", username, addr)

	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return fmt.Errorf("ssh connection failed: %w", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("ssh session failed: %w", err)
	}
	defer session.Close()

	// Stream output live (long output like `kubectl logs` shouldn't buffer).
	// Stdin is intentionally left unset: exec is non-interactive, so the remote
	// command sees EOF immediately — attaching a TTY stdin here could hang.
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	runErr := session.Run(remoteCmd)
	if runErr != nil {
		// Propagate the remote command's exit code as the CLI's exit code.
		var exitErr *ssh.ExitError
		if errors.As(runErr, &exitErr) {
			session.Close()
			client.Close()
			os.Exit(exitErr.ExitStatus())
		}
		return fmt.Errorf("command execution failed: %w", runErr)
	}

	return nil
}
