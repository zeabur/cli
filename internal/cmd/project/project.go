// Package project contains the cmd for managing projects
package project

import (
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"

	projectCreateCmd "github.com/zeabur/cli/internal/cmd/project/create"
	projectGetCmd "github.com/zeabur/cli/internal/cmd/project/get"
	projectListCmd "github.com/zeabur/cli/internal/cmd/project/list"
)

// NewCmdProject creates the project command
func NewCmdProject(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project <command>",
		Short: "Manage projects",
		//todo: is this alias too short?
		Aliases: []string{"p"},
	}

	cmd.AddCommand(projectGetCmd.NewCmdGet(f))
	cmd.AddCommand(projectListCmd.NewCmdList(f))
	cmd.AddCommand(projectCreateCmd.NewCmdCreate(f))

	return cmd
}
