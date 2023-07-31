// Package root provides the root command
package root

import (
	"errors"
	"fmt"
	"github.com/zeabur/cli/pkg/fill"
	"github.com/zeabur/cli/pkg/selector"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	completionCmd "github.com/zeabur/cli/internal/cmd/completion"

	authCmd "github.com/zeabur/cli/internal/cmd/auth"
	contextCmd "github.com/zeabur/cli/internal/cmd/context"
	deploymentCmd "github.com/zeabur/cli/internal/cmd/deployment"
	projectCmd "github.com/zeabur/cli/internal/cmd/project"
	serviceCmd "github.com/zeabur/cli/internal/cmd/service"
	versionCmd "github.com/zeabur/cli/internal/cmd/version"
	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/pkg/api"
	"github.com/zeabur/cli/pkg/config"
	"github.com/zeabur/cli/pkg/log"
)

// NewCmdRoot creates the root command
func NewCmdRoot(f *cmdutil.Factory, version, commit, date string) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "zeabur <command> <subcommand> [flags]",
		Short: "Zeabur CLI",
		Long:  `Zeabur CLI is the official command line tool for Zeabur.`,
		Example: heredoc.Doc(`
			$ zeabur auth login
			$ zeabur project list
			$ zeabur service create
		`),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// set up logging
			if f.Debug {
				f.Log = log.NewDebugLevel()
			} else {
				f.Log = log.NewInfoLevel()
			}

			// require that the user is authenticated before running most commands
			if cmdutil.IsAuthCheckEnabled(cmd) {
				f.Log.Debug("Checking authentication")
				if !f.LoggedIn() {
					return errors.New("not authenticated")
				}
				// set up the client
				f.ApiClient = api.New(f.Config.GetTokenString())
				f.Selector = selector.New(f.ApiClient, f.Log, f.Prompter)
				f.ParamFiller = fill.NewParamFiller(f.Selector)
			}

			return nil
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			err := f.Config.Write()
			if err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}

			// refresh the token if the token is from OAuth2 and it's expired
			if f.AutoRefreshToken && f.LoggedIn() && f.Config.GetToken() != nil {
				if f.Config.GetToken().Expiry.Before(time.Now()) {
					f.Log.Info("Token is from OAuth2 and it's expired, refreshing it")

					token := f.Config.GetToken()
					newToken, err := f.AuthClient.RefreshToken(token)
					if err != nil {
						return fmt.Errorf("failed to refresh token, it is recommended to logout and login again: %w", err)
					}
					f.Config.SetToken(newToken)
					f.Config.SetTokenString(newToken.AccessToken)
					if err := f.Config.Write(); err != nil {
						return fmt.Errorf("failed to save config: %w", err)
					}

					f.Log.Info("Token refreshed successfully")
				}
			}

			return nil
		},
	}

	// don't print usage when error happens
	cmd.SilenceUsage = true
	// don't print error when error happens(we will print it ourselves)
	cmd.SilenceErrors = true

	// Persistent flags
	cmd.PersistentFlags().BoolVar(&f.Debug, "debug", false, "Enable debug logging")
	cmd.PersistentFlags().BoolVarP(&f.Interactive, config.KeyInteractive, "i", true, "use interactive mode")
	cmd.PersistentFlags().BoolVar(&f.AutoRefreshToken, config.KeyAutoRefreshToken, true,
		"automatically refresh token when it's expired, only works when the token is from browser(OAuth2)")

	// Child commands
	cmd.AddCommand(versionCmd.NewCmdVersion(f, version, commit, date))
	cmd.AddCommand(authCmd.NewCmdAuth(f))
	cmd.AddCommand(projectCmd.NewCmdProject(f))
	cmd.AddCommand(serviceCmd.NewCmdService(f))
	cmd.AddCommand(deploymentCmd.NewCmdDeployment(f))

	cmd.AddCommand(contextCmd.NewCmdContext(f))

	cmd.AddCommand(completionCmd.NewCmdCompletion(f))

	return cmd, nil
}
