// Package login provides the login command
package login

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hasura/go-graphql-client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"

	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/pkg/api"
	"github.com/zeabur/cli/pkg/config"
)

// Options is the struct for the login command
type Options struct {
	NewClient func(string) api.Client // to mock in tests
}

// NewCmdLogin creates the login command
func NewCmdLogin(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{
		NewClient: api.New,
	}
	cmd := &cobra.Command{
		Use:   "login",
		Short: "LoggedIn to Zeabur",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunLogin(f, opts)
		},
	}

	cmd.Flags().String(config.KeyTokenString, "", "Zeabur token to use for authentication")
	err := viper.BindPFlag(config.KeyTokenString, cmd.Flags().Lookup(config.KeyTokenString))
	if err != nil {
		return nil
	}

	return cmd
}

// Note: don't import other packages directly in this function, or it will be hard to mock and test
// If you want to add new dependencies, please add them in the Options struct, like `NewClient func(string) client.Client`

// RunLogin is the actual execution of the login command
// Note: it needn't auth, so f.ApiClient is nil
func RunLogin(f *cmdutil.Factory, opts *Options) error {
	if f.Interactive {
		f.Log.Debug("Running login in interactive mode")
	} else {
		f.Log.Debug("Running login in non-interactive mode")
	}

	if f.LoggedIn() {
		f.ApiClient = opts.NewClient(f.Config.GetTokenString())
		user, err := f.ApiClient.GetUserInfo(context.Background())
		if err != nil {
			var graphqlErrors graphql.Errors
			if errors.As(err, &graphqlErrors) &&
				len(graphqlErrors) > 0 &&
				strings.HasPrefix(graphqlErrors[0].Message, "401 Unauthorized") {
				f.Log.Debug("Token is expired or invalid, need to login again")
			} else {
				return fmt.Errorf("failed to get user info: %w", err)
			}
		} else {
			f.Log.Debugw("Already logged in", "token string", f.Config.GetTokenString(),
				"token detail", f.Config.GetToken(), "user", user)
			f.Log.Infof("Already logged in as %s, "+
				"if you want to use a different account, please logout first", user.Name)
			return nil
		}
	}

	var (
		tokenString string
		token       *oauth2.Token
		err         error
	)

	if f.Interactive {
		f.Log.Info("A browser window will be opened for you to login, please confirm")
		// get token from web
		token, err := f.AuthClient.GenerateToken(context.Background())
		if err != nil {
			return fmt.Errorf("failed to login: %w", err)
		}
		tokenString = token

		f.Config.SetTokenString(tokenString)
	} else {
		// get token from flag, env or config
		if tokenString = f.Config.GetTokenString(); tokenString == "" {
			return fmt.Errorf("please set ZEABUR_TOKEN environment variable or use --token flag to set token")
		}
	}

	f.Log.Debugw("Token", "token", tokenString)

	// because we just logged in, we need to create a new client
	f.ApiClient = opts.NewClient(tokenString)

	user, err := f.ApiClient.GetUserInfo(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get user info: %w", err)
	}

	f.Config.SetUser(user.Name)         // nickname
	f.Config.SetUsername(user.Username) // username(same as GitHub id)

	f.Log.Infow("Logged in as", "user", user.Name, "email", user.Email)

	return nil
}
