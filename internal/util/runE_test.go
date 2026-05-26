package util_test

import (
	"testing"

	"github.com/zeabur/cli/internal/util"
	"github.com/zeabur/cli/pkg/zcontext"
)

// TestDefaultIDNameByContext_LazyEvaluation guards the PLA-1590 invariant
// that this helper resolves its source closure at PreRunE time, not at the
// caller's call site. The two project / delete callers depend on this so
// that `f.EffectiveContext().GetProject()` reflects whatever workspace
// override `PersistentPreRunE` has resolved by the time the helper runs.
//
// The test threads through a mutable BasicInfo: we swap it AFTER calling
// DefaultIDNameByContext and BEFORE invoking the returned PreRunE. If the
// helper captured the BasicInfo eagerly (the pre-fix behaviour) the swap
// would be ignored and id/name would be filled from the original value.
func TestDefaultIDNameByContext_LazyEvaluation(t *testing.T) {
	var source = zcontext.NewBasicInfo("old-id", "old-name")
	getter := func() zcontext.BasicInfo { return source }

	var id, name string
	preRun := util.DefaultIDNameByContext(getter, &id, &name)

	// Swap the source AFTER constructing PreRunE; eager capture would miss this.
	source = zcontext.NewBasicInfo("new-id", "new-name")

	if err := preRun(nil, nil); err != nil {
		t.Fatalf("preRun: %v", err)
	}
	if id != "new-id" || name != "new-name" {
		t.Fatalf("got id=%q name=%q, want new-id/new-name (lazy resolution)", id, name)
	}
}

// TestDefaultIDNameByContext_EmptyBasicInfoSkipsFill is the override case:
// when EffectiveContext returns an empty BasicInfo (because --workspace is
// active and the ephemeral context starts empty), no auto-fill happens. The
// caller's runE then has to handle missing id/name itself, which is what
// produces the "please specify project by --name or --id" actionable error
// we observed in dev-2 E2E C3.
func TestDefaultIDNameByContext_EmptyBasicInfoSkipsFill(t *testing.T) {
	empty := zcontext.NewBasicInfo("", "")
	getter := func() zcontext.BasicInfo { return empty }

	id, name := "", ""
	preRun := util.DefaultIDNameByContext(getter, &id, &name)
	if err := preRun(nil, nil); err != nil {
		t.Fatalf("preRun: %v", err)
	}
	if id != "" || name != "" {
		t.Fatalf("got id=%q name=%q, want both empty (override path)", id, name)
	}
}

// TestDefaultIDNameByContext_RespectsUserFlags: when the user has already
// passed --id or --name explicitly, the helper must not overwrite their
// value with the context default. This is unchanged from pre-PLA-1590 and
// guards the back-compat path.
func TestDefaultIDNameByContext_RespectsUserFlags(t *testing.T) {
	source := zcontext.NewBasicInfo("ctx-id", "ctx-name")
	getter := func() zcontext.BasicInfo { return source }

	cases := []struct {
		name string
		id   string
		nm   string
	}{
		{"id-only", "user-id", ""},
		{"name-only", "", "user-name"},
		{"both", "user-id", "user-name"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			id, name := tc.id, tc.nm
			preRun := util.DefaultIDNameByContext(getter, &id, &name)
			if err := preRun(nil, nil); err != nil {
				t.Fatalf("preRun: %v", err)
			}
			if id != tc.id || name != tc.nm {
				t.Errorf("user flag overwritten: got id=%q name=%q, want id=%q name=%q",
					id, name, tc.id, tc.nm)
			}
		})
	}
}
