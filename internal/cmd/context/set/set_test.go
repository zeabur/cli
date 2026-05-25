package set

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/spf13/viper"

	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/pkg/api"
	"github.com/zeabur/cli/pkg/config"
	"github.com/zeabur/cli/pkg/model"
	"github.com/zeabur/cli/pkg/zcontext"
	"go.uber.org/zap"
)

// stubProjectClient stubs the two ApiClient methods setProject's ID path
// actually reaches for. Other Client methods inherit the embedded nil
// interface and panic if accidentally exercised — a safety net so a future
// caller can't sneak in an unverified backend call.
type stubProjectClient struct {
	api.Client

	getProjectRet *model.Project
	getProjectErr error

	listAllOwner       string
	listAllOwnerCalled bool
	listAllRet         model.Projects
	listAllErr         error
}

func (c *stubProjectClient) GetProject(_ context.Context, _, _, _ string) (*model.Project, error) {
	return c.getProjectRet, c.getProjectErr
}

func (c *stubProjectClient) ListAllProjects(_ context.Context, ownerID string) (model.Projects, error) {
	c.listAllOwnerCalled = true
	c.listAllOwner = ownerID
	return c.listAllRet, c.listAllErr
}

// stubConfig is a minimal config.Config that wraps a viper instance for the
// inner context only. setProject writes through `f.Config.GetContext()` so
// the persisted side has to be a real-ish thing.
type stubConfig struct {
	v *viper.Viper
}

func newStubConfig() *stubConfig {
	return &stubConfig{v: viper.New()}
}

func (s *stubConfig) GetTokenString() string         { return "" }
func (s *stubConfig) SetTokenString(string)          {}
func (s *stubConfig) GetUser() string                { return "" }
func (s *stubConfig) SetUser(string)                 {}
func (s *stubConfig) GetUsername() string            { return "alice" }
func (s *stubConfig) SetUsername(string)             {}
func (s *stubConfig) GetContext() zcontext.Context   { return zcontext.NewViperContext(s.v) }
func (s *stubConfig) Write() error                   { return nil }

var _ config.Config = (*stubConfig)(nil)

func newFactory(t *testing.T, apiClient api.Client, persistedWorkspace *zcontext.Workspace) (*cmdutil.Factory, *stubConfig) {
	t.Helper()
	cfg := newStubConfig()
	if persistedWorkspace != nil {
		cfg.GetContext().SetWorkspace(persistedWorkspace)
	}
	f := &cmdutil.Factory{
		Config:    cfg,
		ApiClient: apiClient,
		Log:       zap.NewNop().Sugar(),
	}
	return f, cfg
}

// TestSetProject_ID_TeamWorkspace_AllowsOwnProject — the legitimate path.
// In a team workspace, pinning a project that *does* belong to the team
// succeeds and writes the pin to context. Critically, ListAllProjects is
// called with the current team's ID — the membership check is real.
func TestSetProject_ID_TeamWorkspace_AllowsOwnProject(t *testing.T) {
	target := &model.Project{ID: "65aa1234567890abcdef1234", Name: "team-A-foo"}
	stub := &stubProjectClient{
		getProjectRet: target,
		listAllRet:    model.Projects{target},
	}
	teamA := &zcontext.Workspace{ID: "65cc1230000000000000000a", Name: "team-A", Kind: zcontext.WorkspaceKindTeam}
	f, cfg := newFactory(t, stub, teamA)

	if err := setProject(f, target.ID, "", true); err != nil {
		t.Fatalf("setProject: %v", err)
	}
	if !stub.listAllOwnerCalled || stub.listAllOwner != teamA.ID {
		t.Fatalf("ListAllProjects called with owner=%q (called=%v), want %q", stub.listAllOwner, stub.listAllOwnerCalled, teamA.ID)
	}
	if got := cfg.GetContext().GetProject().GetID(); got != target.ID {
		t.Errorf("project context = %q, want %q (legitimate pin should succeed)", got, target.ID)
	}
}

// TestSetProject_ID_TeamWorkspace_RejectsForeignProject is the Codex
// finding's exact attack: persisted workspace team-A, --id of a team-B
// project. The CLI must refuse and leave persisted context untouched —
// otherwise subsequent name-based service / variable / etc. commands
// would silently operate on team-B.
func TestSetProject_ID_TeamWorkspace_RejectsForeignProject(t *testing.T) {
	teamBProject := &model.Project{ID: "65bb5678901234abcdef5678", Name: "team-B-foo"}
	// ListAllProjects for team-A returns only team-A projects — the team-B
	// project ID is NOT in the list.
	teamAProject := &model.Project{ID: "65aa1234567890abcdef1234", Name: "team-A-only"}
	stub := &stubProjectClient{
		getProjectRet: teamBProject, // GetProject by ID still returns the cross-team project (backend doesn't gate)
		listAllRet:    model.Projects{teamAProject},
	}
	teamA := &zcontext.Workspace{ID: "65cc1230000000000000000a", Name: "team-A", Kind: zcontext.WorkspaceKindTeam}
	f, cfg := newFactory(t, stub, teamA)

	err := setProject(f, teamBProject.ID, "", true)
	if err == nil {
		t.Fatal("setProject must refuse cross-workspace --id, got nil error")
	}
	if !strings.Contains(err.Error(), "does not belong to workspace") {
		t.Errorf("error should explain cross-workspace mismatch, got: %v", err)
	}
	// Project context must be untouched.
	if got := cfg.GetContext().GetProject().GetID(); got != "" {
		t.Errorf("context contaminated on rejection: got project.id=%q", got)
	}
}

// TestSetProject_ID_PersonalWorkspace_BypassesCheck is the back-compat
// guard. Personal workspace must NOT call ListAllProjects — that would
// break collaborator workflows where a user pins by-ID a project owned by
// someone else (Codex's explicit non-regression).
func TestSetProject_ID_PersonalWorkspace_BypassesCheck(t *testing.T) {
	collaboratorProject := &model.Project{ID: "65aa1234567890abcdef1234", Name: "shared-with-me"}
	stub := &stubProjectClient{
		getProjectRet: collaboratorProject,
		// listAllRet deliberately unset — if the code paths into
		// ListAllProjects, the empty list would reject and the test fails.
	}
	// Personal workspace: persistedWorkspace nil.
	f, cfg := newFactory(t, stub, nil)

	if err := setProject(f, collaboratorProject.ID, "", true); err != nil {
		t.Fatalf("personal --id must not be gated, got: %v", err)
	}
	if stub.listAllOwnerCalled {
		t.Errorf("personal --id must NOT call ListAllProjects (collaborator workflow), but it was called with owner=%q", stub.listAllOwner)
	}
	if got := cfg.GetContext().GetProject().GetID(); got != collaboratorProject.ID {
		t.Errorf("personal pin failed: context.project.id=%q, want %q", got, collaboratorProject.ID)
	}
}

// TestSetProject_ID_TeamWorkspace_ListErr propagates: if the membership
// list call itself fails, we must NOT fall back to "trust the user" —
// that would re-open the cross-workspace gap whenever the backend is
// flaky. Surface the error.
func TestSetProject_ID_TeamWorkspace_ListErr(t *testing.T) {
	target := &model.Project{ID: "65aa1234567890abcdef1234", Name: "x"}
	stub := &stubProjectClient{
		getProjectRet: target,
		listAllErr:    errors.New("boom"),
	}
	teamA := &zcontext.Workspace{ID: "65cc1230000000000000000a", Name: "team-A", Kind: zcontext.WorkspaceKindTeam}
	f, cfg := newFactory(t, stub, teamA)

	err := setProject(f, target.ID, "", true)
	if err == nil {
		t.Fatal("ListAllProjects failure must propagate, not silently pass")
	}
	if !strings.Contains(err.Error(), "verify project workspace membership") || !strings.Contains(err.Error(), "boom") {
		t.Errorf("error should wrap the list call failure, got: %v", err)
	}
	if got := cfg.GetContext().GetProject().GetID(); got != "" {
		t.Errorf("context contaminated despite list error: %q", got)
	}
}
