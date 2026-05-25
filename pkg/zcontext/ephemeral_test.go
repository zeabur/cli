package zcontext_test

import (
	"testing"

	"github.com/zeabur/cli/pkg/zcontext"
)

// TestEphemeralContext_WorkspaceFromConstructor: the override's workspace must
// flow through GetWorkspace(). Otherwise consumers reading
// `ctx.GetWorkspace()` get "personal" under a team override — a second-order
// trap.
func TestEphemeralContext_WorkspaceFromConstructor(t *testing.T) {
	ws := &zcontext.Workspace{ID: "65aa1234567890abcdef1234", Name: "acme", Kind: zcontext.WorkspaceKindTeam}
	ctx := zcontext.NewEphemeralContext(ws)
	got := ctx.GetWorkspace()
	if got == nil || got.ID != ws.ID || got.Name != ws.Name || got.Kind != ws.Kind {
		t.Fatalf("got %+v, want %+v", got, ws)
	}
	if !got.IsTeam() {
		t.Errorf("override workspace must report IsTeam(), got %+v", got)
	}
}

// TestEphemeralContext_NilWorkspaceIsPersonal: passing nil at construction
// yields a personal-shaped workspace (zero value). Used when the caller wants
// a scratch context without a known override (uncommon).
func TestEphemeralContext_NilWorkspaceIsPersonal(t *testing.T) {
	ctx := zcontext.NewEphemeralContext(nil)
	ws := ctx.GetWorkspace()
	if ws == nil {
		t.Fatal("GetWorkspace must never return nil")
	}
	if !ws.IsPersonal() {
		t.Errorf("nil constructor must produce personal workspace, got %+v", ws)
	}
}

// TestEphemeralContext_ReadEmptyByDefault: every inner-context getter starts
// empty. This is the core "no implicit fallback to persisted state" property:
// an interactive command running under override sees no leftovers from the
// persisted config, so it always prompts the user fresh.
func TestEphemeralContext_ReadEmptyByDefault(t *testing.T) {
	ctx := zcontext.NewEphemeralContext(&zcontext.Workspace{ID: "x"})
	for _, tc := range []struct {
		name string
		info zcontext.BasicInfo
	}{
		{"project", ctx.GetProject()},
		{"environment", ctx.GetEnvironment()},
		{"service", ctx.GetService()},
	} {
		if !tc.info.Empty() {
			t.Errorf("ephemeral %s must start empty, got id=%q name=%q", tc.name, tc.info.GetID(), tc.info.GetName())
		}
	}
}

// TestEphemeralContext_SetReadCycleWorksInMemory: ParamFiller's flow depends
// on `Set then later Get` returning what was just Set. This must keep
// working under override — the values just shouldn't leak to disk.
func TestEphemeralContext_SetReadCycleWorksInMemory(t *testing.T) {
	ctx := zcontext.NewEphemeralContext(&zcontext.Workspace{ID: "x"})
	want := zcontext.NewBasicInfo("65aa1234567890abcdef1234", "my-project")
	ctx.SetProject(want)
	got := ctx.GetProject()
	if got.GetID() != "65aa1234567890abcdef1234" || got.GetName() != "my-project" {
		t.Fatalf("got id=%q name=%q, want SetProject value through", got.GetID(), got.GetName())
	}
	// Environment / Service same shape.
	ctx.SetEnvironment(zcontext.NewBasicInfo("env-id", "env-name"))
	if e := ctx.GetEnvironment(); e.GetID() != "env-id" {
		t.Errorf("environment round-trip lost value: %+v", e)
	}
	ctx.SetService(zcontext.NewBasicInfo("svc-id", "svc-name"))
	if s := ctx.GetService(); s.GetID() != "svc-id" {
		t.Errorf("service round-trip lost value: %+v", s)
	}
}

// TestEphemeralContext_ClearAll: ClearAll wipes the inner fields but leaves
// the workspace alone — the override workspace is the whole reason we're
// running ephemerally.
func TestEphemeralContext_ClearAll(t *testing.T) {
	ws := &zcontext.Workspace{ID: "65aa1234567890abcdef1234", Name: "acme", Kind: zcontext.WorkspaceKindTeam}
	ctx := zcontext.NewEphemeralContext(ws)
	ctx.SetProject(zcontext.NewBasicInfo("p", "P"))
	ctx.SetEnvironment(zcontext.NewBasicInfo("e", "E"))
	ctx.SetService(zcontext.NewBasicInfo("s", "S"))
	ctx.ClearAll()
	for _, info := range []zcontext.BasicInfo{ctx.GetProject(), ctx.GetEnvironment(), ctx.GetService()} {
		if !info.Empty() {
			t.Errorf("after ClearAll inner context must be empty, got %+v", info)
		}
	}
	if got := ctx.GetWorkspace(); got.ID != ws.ID {
		t.Errorf("ClearAll must not touch workspace: got %+v want %+v", got, ws)
	}
}
