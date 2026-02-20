package util

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/pkg/zcontext"
)

// NeedProjectContextWhenNonInteractive checks if the project context is set in the non-interactive mode.
// If overrideID is provided and non-empty, the check is skipped (the caller already has a project ID from a flag).
func NeedProjectContextWhenNonInteractive(f *cmdutil.Factory, overrideID ...*string) CobraRunE {
	return func(cmd *cobra.Command, args []string) error {
		if len(overrideID) > 0 && overrideID[0] != nil && *overrideID[0] != "" {
			return nil
		}
		if !f.Interactive && f.Config.GetContext().GetProject().Empty() {
			return errors.New("please run <zeabur context set project> first")
		}
		return nil
	}
}

func DefaultIDNameByContext(basicInfo zcontext.BasicInfo, id, name *string) CobraRunE {
	return func(cmd *cobra.Command, args []string) error {
		defaultByContext(basicInfo, id, name)
		return nil
	}
}

func DefaultIDByContext(basicInfo zcontext.BasicInfo, id *string) CobraRunE {
	var unused string

	return func(cmd *cobra.Command, args []string) error {
		defaultByContext(basicInfo, id, &unused)
		return nil
	}
}

// defaultByContext if id and name both are empty, then use the context to fill them, (param should not be nil)
func defaultByContext(basicInfo zcontext.BasicInfo, id, name *string) {
	if id == nil || name == nil {
		return
	}
	if *id == "" && *name == "" && !basicInfo.Empty() {
		*id = basicInfo.GetID()
		*name = basicInfo.GetName()
	}
}
