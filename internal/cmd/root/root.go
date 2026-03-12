// Package root provides the root command
package root

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	authCmd "github.com/zeabur/cli/internal/cmd/auth"
	completionCmd "github.com/zeabur/cli/internal/cmd/completion"
	helpCmd "github.com/zeabur/cli/internal/cmd/help"
	contextCmd "github.com/zeabur/cli/internal/cmd/context"
	deployCmd "github.com/zeabur/cli/internal/cmd/deploy"
	deploymentCmd "github.com/zeabur/cli/internal/cmd/deployment"
	domainCmd "github.com/zeabur/cli/internal/cmd/domain"
	profileCmd "github.com/zeabur/cli/internal/cmd/profile"
	projectCmd "github.com/zeabur/cli/internal/cmd/project"
	serverCmd "github.com/zeabur/cli/internal/cmd/server"
	serviceCmd "github.com/zeabur/cli/internal/cmd/service"
	templateCmd "github.com/zeabur/cli/internal/cmd/template"
	uploadCmd "github.com/zeabur/cli/internal/cmd/upload"
	variableCmd "github.com/zeabur/cli/internal/cmd/variable"
	versionCmd "github.com/zeabur/cli/internal/cmd/version"
	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/pkg/api"
	"github.com/zeabur/cli/pkg/config"
	"github.com/zeabur/cli/pkg/fill"
	"github.com/zeabur/cli/pkg/log"
	"github.com/zeabur/cli/pkg/selector"
)

// NewCmdRoot creates the root command
func NewCmdRoot(f *cmdutil.Factory, version, commit, date string) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "zeabur",
		Short: "Zeabur CLI",
		Long:  `Zeabur CLI is the official command line tool for Zeabur.`,
		Example: heredoc.Doc(`
			$ zeabur auth login
			$ zeabur project list
			$ zeabur service create
		`),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// set up logging
			if f.JSON {
				f.Log = log.NewSilent()
			} else if f.Debug {
				f.Log = log.NewDebugLevel()
			} else {
				f.Log = log.NewInfoLevel()
			}

			// normalize ID flags: strip prefix from prefixed ObjectIDs
			// e.g. "service-662e24fca7d5..." → "662e24fca7d5..."
			var normalizeErr error
			cmd.Flags().VisitAll(func(flag *pflag.Flag) {
				if !flag.Changed || normalizeErr != nil {
					return
				}
				name := flag.Name
				if name == "id" || strings.HasSuffix(name, "-id") {
					normalizeErr = normalizeIDFlag(flag)
				}
			})
			if normalizeErr != nil {
				return normalizeErr
			}

			// require that the user is authenticated before running most commands
			if cmdutil.IsAuthCheckEnabled(cmd) {
				// in JSON mode, fail fast if not authenticated instead of opening a browser
				if f.JSON && !f.LoggedIn() {
					return fmt.Errorf("not authenticated: run `zeabur auth login` before using --json")
				}

				// do not return error, guide user to login instead
				if !f.LoggedIn() {
					f.Log.Info("A browser window will be opened for you to login, please confirm")

					var (
						tokenString string
						err         error
					)

					tokenString, err = f.AuthClient.GenerateToken(context.Background())
					if err != nil {
						return fmt.Errorf("failed to login: %w", err)
					}
					f.Config.SetTokenString(tokenString)
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

			if f.AutoCheckUpdate && !f.Debug && version != "dev" {
				currentVersion := TrimPrefixV(version)

				upstreamVersionInfo, err := GetLatestRelease("zeabur/cli")
				if err != nil {
					return nil
				}

				upstreamVersion := TrimPrefixV(upstreamVersionInfo.TagName)

				needUpdate, err := IsVersionNewerSemver(upstreamVersion, currentVersion)
				if err != nil {
					return nil
				}

				if needUpdate {
					f.Log.Infof("A new version of Zeabur CLI is available: %s, you are using %s", upstreamVersion, currentVersion)
					f.Log.Infof("Please visit %s to download the latest version", upstreamVersionInfo.URL)
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
	cmd.PersistentFlags().BoolVar(&f.AutoCheckUpdate, config.KeyAutoCheckUpdate, true, "automatically check update")
	cmd.PersistentFlags().BoolVar(&f.JSON, "json", false, "output in JSON format")

	// Child commands
	cmd.AddCommand(deployCmd.NewCmdDeploy(f))
	cmd.AddCommand(uploadCmd.NewCmdUpload(f))
	cmd.AddCommand(versionCmd.NewCmdVersion(f, version, commit, date))
	cmd.AddCommand(authCmd.NewCmdAuth(f))
	cmd.AddCommand(projectCmd.NewCmdProject(f))
	cmd.AddCommand(serverCmd.NewCmdServer(f))
	cmd.AddCommand(serviceCmd.NewCmdService(f))
	cmd.AddCommand(deploymentCmd.NewCmdDeployment(f))
	cmd.AddCommand(templateCmd.NewCmdTemplate(f))
	cmd.AddCommand(domainCmd.NewCmdDomain(f))
	cmd.AddCommand(profileCmd.NewCmdProfile(f))
	cmd.AddCommand(contextCmd.NewCmdContext(f))
	cmd.AddCommand(completionCmd.NewCmdCompletion(f))
	cmd.AddCommand(variableCmd.NewCmdVariable(f))

	// replace default help command with our custom one that supports --all
	cmd.SetHelpCommand(helpCmd.NewCmdHelp(cmd))

	return cmd, nil
}

// normalizeIDFlag strips a known prefix from a prefixed MongoDB ObjectID flag value.
// e.g. "service-662e24fca7d5abcdef123456" → "662e24fca7d5abcdef123456"
func normalizeIDFlag(flag *pflag.Flag) error {
	val := flag.Value.String()
	if idx := strings.LastIndex(val, "-"); idx != -1 {
		suffix := val[idx+1:]
		if len(suffix) != 24 {
			return nil
		}
		if _, err := hex.DecodeString(suffix); err != nil {
			return nil
		}
		if err := flag.Value.Set(suffix); err != nil {
			return fmt.Errorf("normalize %s: %w", flag.Name, err)
		}
	}
	return nil
}
