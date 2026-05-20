package cmdutil_test

import (
	"strings"
	"testing"

	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/pkg/model"
)

func ptrRole(r model.TeamMemberRole) *model.TeamMemberRole { return &r }

func TestResolveWorkspaceArg_EmptyArg(t *testing.T) {
	_, err := cmdutil.ResolveWorkspaceArg(nil, "  ")
	if err == nil || !strings.Contains(err.Error(), "required") {
		t.Fatalf("want 'required' error, got %v", err)
	}
}

func TestResolveWorkspaceArg_ByID_Match(t *testing.T) {
	id := "65aa1234567890abcdef1234"
	teams := []model.Team{{ID: id, Name: "acme", MyRole: ptrRole(model.TeamMemberRoleAdministrator)}}
	team, err := cmdutil.ResolveWorkspaceArg(teams, id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if team.ID != id {
		t.Fatalf("ID = %s, want %s", team.ID, id)
	}
}

// TestResolveWorkspaceArg_ByID_CaseInsensitive: isObjectIDHex accepts upper
// and lower hex, so the comparison against teams[i].ID must do the same.
// Mongo ObjectIDs are conventionally lowercase, but a user pasting a UUID
// from a browser tab might end up with uppercase — that must still resolve.
func TestResolveWorkspaceArg_ByID_CaseInsensitive(t *testing.T) {
	teams := []model.Team{{ID: "65aa1234567890abcdef1234", Name: "acme"}}
	team, err := cmdutil.ResolveWorkspaceArg(teams, "65AA1234567890ABCDEF1234")
	if err != nil {
		t.Fatalf("uppercase ID should resolve, got %v", err)
	}
	if team.ID != "65aa1234567890abcdef1234" {
		t.Fatalf("returned team ID = %s, want lowercase canonical", team.ID)
	}
}

func TestResolveWorkspaceArg_ByID_NotAMember(t *testing.T) {
	teams := []model.Team{{ID: "65aa1234567890abcdef1234", Name: "acme"}}
	_, err := cmdutil.ResolveWorkspaceArg(teams, "65bbffffffffffffffffffff")
	if err == nil || !(strings.Contains(err.Error(), "not a team") || strings.Contains(err.Error(), "no team")) {
		t.Fatalf("want membership error, got %v", err)
	}
}

func TestResolveWorkspaceArg_ByName_Unique(t *testing.T) {
	teams := []model.Team{
		{ID: "65aa1234567890abcdef1234", Name: "acme", MyRole: ptrRole(model.TeamMemberRoleEditor)},
		{ID: "65bb5678901234abcdef5678", Name: "beta", MyRole: ptrRole(model.TeamMemberRoleViewer)},
	}
	team, err := cmdutil.ResolveWorkspaceArg(teams, "acme")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if team.Name != "acme" {
		t.Fatalf("Name = %q, want acme", team.Name)
	}
}

func TestResolveWorkspaceArg_ByName_NotFound(t *testing.T) {
	teams := []model.Team{{ID: "65aa1234567890abcdef1234", Name: "acme"}}
	_, err := cmdutil.ResolveWorkspaceArg(teams, "zeta")
	if err == nil || !strings.Contains(err.Error(), `no workspace named`) {
		t.Fatalf("want 'no workspace named' error, got %v", err)
	}
}

// TestResolveWorkspaceArg_ByName_Ambiguous covers the duplicate-team-name
// case Bruce called out explicitly: two teams sharing a name must not be
// resolvable by name alone, and the error must spell out the disambiguating
// commands so the user can pick by ID.
func TestResolveWorkspaceArg_ByName_Ambiguous(t *testing.T) {
	id1 := "65aa1234567890abcdef1234"
	id2 := "65bb5678901234abcdef5678"
	teams := []model.Team{
		{ID: id1, Name: "acme", MyRole: ptrRole(model.TeamMemberRoleAdministrator)},
		{ID: id2, Name: "acme", MyRole: ptrRole(model.TeamMemberRoleViewer)},
	}
	_, err := cmdutil.ResolveWorkspaceArg(teams, "acme")
	if err == nil {
		t.Fatal("want ambiguous error, got nil")
	}
	msg := err.Error()
	for _, want := range []string{"ambiguous", "2 workspaces named", id1, id2, "Administrator", "Viewer", "workspace switch"} {
		if !strings.Contains(msg, want) {
			t.Errorf("error missing %q\nfull: %s", want, msg)
		}
	}
}

// TestResolveWorkspaceArg_NonHex_NotMistakenForID guards the hex-vs-name
// branch: a 24-char string that isn't valid hex falls into the name path,
// not the ID path. Otherwise a user typing a name that happens to be 24
// characters long would get the misleading "no team with id" error.
func TestResolveWorkspaceArg_NonHex_NotMistakenForID(t *testing.T) {
	teams := []model.Team{{ID: "65aa1234567890abcdef1234", Name: "non-hex-but-twentyfour"}}
	team, err := cmdutil.ResolveWorkspaceArg(teams, "non-hex-but-twentyfour")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if team.ID != "65aa1234567890abcdef1234" {
		t.Fatalf("team ID = %s, want 65aa...", team.ID)
	}
}
