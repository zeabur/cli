package util_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/zeabur/cli/internal/util"
	"github.com/zeabur/cli/pkg/api"
	"github.com/zeabur/cli/pkg/model"
)

// fakeProjectClient stubs the two ProjectAPI methods util.GetProjectByName
// reaches for. Everything else inherits the nil embedded interface and panics
// if accidentally exercised.
type fakeProjectClient struct {
	api.Client

	// personal-path arguments captured for assertions
	getProjectCalls []getProjectCall
	getProjectRet   *model.Project
	getProjectErr   error

	// team-path arguments captured for assertions
	listAllOwner    string
	listAllProjects []*model.Project
	listAllErr      error
}

type getProjectCall struct{ id, owner, name string }

func (c *fakeProjectClient) GetProject(_ context.Context, id, owner, name string) (*model.Project, error) {
	c.getProjectCalls = append(c.getProjectCalls, getProjectCall{id, owner, name})
	return c.getProjectRet, c.getProjectErr
}

func (c *fakeProjectClient) ListAllProjects(_ context.Context, ownerID string) (model.Projects, error) {
	c.listAllOwner = ownerID
	return c.listAllProjects, c.listAllErr
}

// TestGetProjectByName_Personal: ownerID == "" → backend `project(owner, name)`
// query against the personal username, unchanged from before workspace support.
func TestGetProjectByName_Personal(t *testing.T) {
	want := &model.Project{ID: "65aa1234567890abcdef1234", Name: "api"}
	c := &fakeProjectClient{getProjectRet: want}

	got, err := util.GetProjectByName(c, "", "alice", "api")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
	if len(c.getProjectCalls) != 1 {
		t.Fatalf("expected 1 GetProject call, got %d", len(c.getProjectCalls))
	}
	call := c.getProjectCalls[0]
	if call.id != "" || call.owner != "alice" || call.name != "api" {
		t.Fatalf("call args = %+v, want id=\"\" owner=alice name=api", call)
	}
	if c.listAllOwner != "" {
		t.Fatalf("ListAllProjects should not have been called on personal path (owner=%q)", c.listAllOwner)
	}
}

// TestGetProjectByName_TeamFound: ownerID set → ListAllProjects(ownerID) +
// match by name. Critically, MUST NOT touch the personal-username path.
func TestGetProjectByName_TeamFound(t *testing.T) {
	teamID := "65cc1234567890abcdef0000"
	team1 := &model.Project{ID: "65aa1234567890abcdef1234", Name: "api"}
	team2 := &model.Project{ID: "65bb5678901234abcdef5678", Name: "web"}
	c := &fakeProjectClient{listAllProjects: []*model.Project{team1, team2}}

	got, err := util.GetProjectByName(c, teamID, "alice", "web")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != team2.ID {
		t.Fatalf("got ID %q, want %q", got.ID, team2.ID)
	}
	if c.listAllOwner != teamID {
		t.Fatalf("ListAllProjects called with owner %q, want %q", c.listAllOwner, teamID)
	}
	if len(c.getProjectCalls) != 0 {
		t.Fatalf("personal GetProject path must not run in team workspace; got %d calls", len(c.getProjectCalls))
	}
}

// TestGetProjectByName_TeamNotFound: ownerID set, name missing in the
// team's list → error names the workspace, NOT a 404 from a personal lookup.
func TestGetProjectByName_TeamNotFound(t *testing.T) {
	c := &fakeProjectClient{listAllProjects: []*model.Project{{ID: "x", Name: "api"}}}

	_, err := util.GetProjectByName(c, "65cc1234567890abcdef0000", "alice", "missing")
	if err == nil {
		t.Fatal("want error, got nil")
	}
	if !strings.Contains(err.Error(), "no project named") {
		t.Fatalf("error = %v, want 'no project named ...' message", err)
	}
}

// TestGetProjectByName_TeamListErr: a backend failure on the team path must
// propagate the underlying error, not silently fall through to personal.
func TestGetProjectByName_TeamListErr(t *testing.T) {
	c := &fakeProjectClient{listAllErr: errors.New("boom")}

	_, err := util.GetProjectByName(c, "65cc1234567890abcdef0000", "alice", "api")
	if err == nil {
		t.Fatal("want error, got nil")
	}
	if !strings.Contains(err.Error(), "list projects in workspace") || !strings.Contains(err.Error(), "boom") {
		t.Fatalf("error = %v, want wrapped boom", err)
	}
	if len(c.getProjectCalls) != 0 {
		t.Fatalf("personal fallback must not run when team list fails; got %d calls", len(c.getProjectCalls))
	}
}
