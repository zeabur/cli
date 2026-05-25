package util

import (
	"github.com/spf13/cobra"
	"github.com/zeabur/cli/pkg/zcontext"
)

// DefaultIDNameByContext returns a Cobra PreRunE that auto-fills `id` and
// `name` from the resource the supplied closure returns, provided both are
// empty. Useful for commands like `project get/delete` that should default
// to "whatever the user has pinned" when invoked with no flags.
//
// The argument is a closure rather than the BasicInfo directly so the
// resolution happens at PreRunE time (after global flags are parsed),
// not at Cobra command construction time. This matters for PLA-1590:
// `--workspace` only takes effect after PersistentPreRunE runs, so a
// caller passing `f.EffectiveContext().GetProject()` directly would
// capture the persisted project, not the override-aware empty.
func DefaultIDNameByContext(basicInfoFn func() zcontext.BasicInfo, id, name *string) CobraRunE {
	return func(cmd *cobra.Command, args []string) error {
		defaultByContext(basicInfoFn(), id, name)
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
