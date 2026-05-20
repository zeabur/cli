package cmdutil_test

import (
	"context"
	"errors"
	"testing"

	"github.com/spf13/viper"
	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/pkg/api"
	"github.com/zeabur/cli/pkg/model"
	"github.com/zeabur/cli/pkg/zcontext"
)

// stubConfig is a minimal config.Config — only GetContext is exercised by
// the workspace plumbing. Other methods are no-ops so any test that needs
// them would notice the missing behavior.
type stubConfig struct {
	ctx zcontext.Context
}

func (s stubConfig) GetTokenString() string       { return "" }
func (s stubConfig) SetTokenString(string)        {}
func (s stubConfig) GetUser() string              { return "" }
func (s stubConfig) SetUser(string)               {}
func (s stubConfig) GetUsername() string          { return "" }
func (s stubConfig) SetUsername(string)           {}
func (s stubConfig) GetContext() zcontext.Context { return s.ctx }
func (s stubConfig) Write() error                 { return nil }

// TestFactory_PersonalUserInvariant guards the single most important
// backward-compat rule of PLA-1590: a brand-new caller — one who has never
// run `workspace switch` and never set `--workspace` — must report
// CurrentOwnerID() == "". That empty string is what every owner-aware util
// helper checks before deciding between the legacy personal query path and
// the new team-aware one. If this ever returns non-empty for a vanilla
// Factory we'd be silently routing personal users through the team branch.
func TestFactory_PersonalUserInvariant(t *testing.T) {
	f := &cmdutil.Factory{} // zero Factory = brand-new user
	if got := f.CurrentOwnerID(); got != "" {
		t.Fatalf("brand-new Factory must report personal (empty ownerID), got %q", got)
	}
	ws := f.CurrentWorkspace()
	if !ws.IsPersonal() {
		t.Fatalf("brand-new Factory must report personal workspace, got %+v", ws)
	}
}

// TestFactory_CurrentOwnerID_PersistedPersonal mirrors the run-time shape:
// a Factory with Config but the workspace field unset (the explicit
// "personal" persisted state, or a fresh config file). Must still report
// personal.
func TestFactory_CurrentOwnerID_PersistedPersonal(t *testing.T) {
	cfg := stubConfig{ctx: zcontext.NewViperContext(viper.New())}
	f := &cmdutil.Factory{Config: cfg}
	if got := f.CurrentOwnerID(); got != "" {
		t.Fatalf("persisted empty workspace must report personal, got %q", got)
	}
	if !f.CurrentWorkspace().IsPersonal() {
		t.Fatalf("CurrentWorkspace().IsPersonal() must be true for empty persisted")
	}
}

// TestFactory_CurrentOwnerID_PersistedTeam: the persisted-team path returns
// the team's ID and a workspace marked IsTeam().
func TestFactory_CurrentOwnerID_PersistedTeam(t *testing.T) {
	v := viper.New()
	v.Set("workspace.id", "65aa1234567890abcdef1234")
	v.Set("workspace.name", "acme")
	v.Set("workspace.kind", zcontext.WorkspaceKindTeam)
	cfg := stubConfig{ctx: zcontext.NewViperContext(v)}
	f := &cmdutil.Factory{Config: cfg}

	if got := f.CurrentOwnerID(); got != "65aa1234567890abcdef1234" {
		t.Fatalf("got %q, want persisted team ID", got)
	}
	ws := f.CurrentWorkspace()
	if !ws.IsTeam() || ws.Name != "acme" {
		t.Fatalf("got %+v, want team workspace acme", ws)
	}
}

// TestFactory_CurrentOwnerID_OverrideBeatsPersisted: --workspace flag
// resolution sets an override on the Factory that must take precedence
// over the persisted workspace for the lifetime of that invocation. The
// persisted file is left alone.
func TestFactory_CurrentOwnerID_OverrideBeatsPersisted(t *testing.T) {
	v := viper.New()
	v.Set("workspace.id", "persisted-id")
	v.Set("workspace.name", "persisted-team")
	v.Set("workspace.kind", zcontext.WorkspaceKindTeam)
	cfg := stubConfig{ctx: zcontext.NewViperContext(v)}
	f := &cmdutil.Factory{Config: cfg}

	f.SetWorkspaceOverride(&zcontext.Workspace{
		ID: "override-id", Name: "override-team", Kind: zcontext.WorkspaceKindTeam,
	})
	if got := f.CurrentOwnerID(); got != "override-id" {
		t.Fatalf("got %q, want override-id", got)
	}
	if ws := f.CurrentWorkspace(); ws.Name != "override-team" {
		t.Fatalf("got %q, want override name", ws.Name)
	}
	// Persisted file untouched: the override is process-local.
	if v.GetString("workspace.id") != "persisted-id" {
		t.Fatalf("override leaked into persisted config: %q", v.GetString("workspace.id"))
	}
}

// TestFactory_CurrentOwnerID_OverrideNilClears: passing nil to
// SetWorkspaceOverride drops the override and the persisted workspace
// becomes effective again.
func TestFactory_CurrentOwnerID_OverrideNilClears(t *testing.T) {
	cfg := stubConfig{ctx: zcontext.NewViperContext(viper.New())}
	f := &cmdutil.Factory{Config: cfg}
	f.SetWorkspaceOverride(&zcontext.Workspace{ID: "abc"})
	if got := f.CurrentOwnerID(); got != "abc" {
		t.Fatalf("got %q, want abc after set", got)
	}
	f.SetWorkspaceOverride(nil)
	if got := f.CurrentOwnerID(); got != "" {
		t.Fatalf("got %q, want empty after clear", got)
	}
}

// fakeListTeamsAPI counts ListTeams invocations so the cache tests can
// assert how many backend round-trips the Factory makes.
type fakeListTeamsAPI struct {
	api.Client
	calls    int
	teamsRet []model.Team
	errRet   error
}

func (f *fakeListTeamsAPI) ListTeams(_ context.Context) ([]model.Team, error) {
	f.calls++
	return f.teamsRet, f.errRet
}

// TestFactory_ListTeams_Memoizes guards F2 from the review: a single CLI
// invocation can ask for the teams list from up to three sites (flag
// resolution, lazy verify, the command itself). They must all share one
// backend call.
func TestFactory_ListTeams_Memoizes(t *testing.T) {
	stub := &fakeListTeamsAPI{teamsRet: []model.Team{{ID: "x"}}}
	f := &cmdutil.Factory{ApiClient: stub}

	for i := 0; i < 3; i++ {
		teams, err := f.ListTeams(context.Background())
		if err != nil {
			t.Fatalf("call %d: %v", i, err)
		}
		if len(teams) != 1 || teams[0].ID != "x" {
			t.Fatalf("call %d: got %+v", i, teams)
		}
	}
	if stub.calls != 1 {
		t.Errorf("ListTeams hit backend %d times, want exactly 1 (memoized)", stub.calls)
	}
}

// TestFactory_ListTeams_StickyError: a failed fetch is cached as a failure
// so subsequent callers within the same process don't retry against an
// already-known-broken backend.
func TestFactory_ListTeams_StickyError(t *testing.T) {
	stub := &fakeListTeamsAPI{errRet: errors.New("boom")}
	f := &cmdutil.Factory{ApiClient: stub}

	for i := 0; i < 3; i++ {
		if _, err := f.ListTeams(context.Background()); err == nil {
			t.Fatalf("call %d: want error, got nil", i)
		}
	}
	if stub.calls != 1 {
		t.Errorf("ListTeams hit backend %d times, want exactly 1 (sticky cache)", stub.calls)
	}
}
