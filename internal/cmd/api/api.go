// Package api provides the `zeabur api` command for raw Zeabur API requests.
package api

import (
	"github.com/spf13/cobra"

	graphqlCmd "github.com/zeabur/cli/internal/cmd/api/graphql"
	"github.com/zeabur/cli/internal/cmdutil"
)

func NewCmdAPI(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "api <command>",
		Short: "Make authenticated requests to the Zeabur API",
	}

	cmd.AddCommand(graphqlCmd.NewCmdGraphQL(f))

	return cmd
}
