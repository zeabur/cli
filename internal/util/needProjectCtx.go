package util

import (
	"errors"
	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
)

// NeedProjectContext checks if the project context is set in the non-interactive mode
func NeedProjectContext(f *cmdutil.Factory) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if !f.Interactive && f.Config.GetContext().GetProject().Empty() {
			return errors.New("please run <zeabur context set project> first")
		}
		return nil
	}
}
