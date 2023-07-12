package clear

import (
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct {
}

func NewCmdClear(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := cobra.Command{
		Use:   "clear",
		Short: "Clear Contexts",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runClear(f, opts)
		},
	}

	return &cmd
}

func runClear(f *cmdutil.Factory, opts *Options) error {
	confirm := true

	if f.Interactive {
		var err error
		confirm, err = f.Prompter.Confirm("Are you sure you want to clear all contexts?", true)
		if err != nil {
			return err
		}
	}
	if confirm {
		f.Config.GetContext().ClearAll()
		f.Log.Info("all contexts cleared")
	}

	return nil
}
