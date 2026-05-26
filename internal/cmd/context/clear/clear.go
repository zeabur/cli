package clear

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct{}

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
	// `context clear` modifies the persisted inner context. Under a
	// `--workspace` override the persisted state belongs to a (potentially)
	// different workspace than the user thinks they're in, so silently
	// clearing it would surprise them. Reject up front and tell them to
	// switch first if they really mean to wipe the persisted context
	// (PLA-1590 B+).
	if f.HasWorkspaceOverride() {
		return fmt.Errorf(
			"`context clear` cannot be combined with `--workspace`; the override does not modify persisted context",
		)
	}

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
