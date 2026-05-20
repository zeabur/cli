package cmdutil

import (
	"context"
	"strings"
	"testing"

	"github.com/zeabur/cli/pkg/api"
	"github.com/zeabur/cli/pkg/model"
)

// fakeListTeamsClient stubs api.Client.ListTeams. Other Client methods are
// inherited from the embedded nil interface and will panic if exercised —
// ResolveWorkspaceArg must only call ListTeams.
type fakeListTeamsClient struct {
	api.Client
	teams []model.Team
}

func (c *fakeListTeamsClient) ListTeams(_ context.Context) ([]model.Team, error) {
	return c.teams, nil
}

func ptrRole(r model.TeamMemberRole) *model.TeamMemberRole { return &r }

func TestResolveWorkspaceArg_EmptyArg(t *testing.T) {
	_, err := ResolveWorkspaceArg(context.Background(), &fakeListTeamsClient{}, "  ")
	if err == nil || !strings.Contains(err.Error(), "required") {
		t.Fatalf("want 'required' error, got %v", err)
	}
}

func TestResolveWorkspaceArg_ByID_Match(t *testing.T) {
	id := "65aa1234567890abcdef1234"
	c := &fakeListTeamsClient{teams: []model.Team{{ID: id, Name: "acme", MyRole: ptrRole(model.TeamMemberRoleAdministrator)}}}
	team, err := ResolveWorkspaceArg(context.Background(), c, id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if team.ID != id {
		t.Fatalf("ID = %s, want %s", team.ID, id)
	}
}

func TestResolveWorkspaceArg_ByID_NotAMember(t *testing.T) {
	c := &fakeListTeamsClient{teams: []model.Team{{ID: "65aa1234567890abcdef1234", Name: "acme"}}}
	_, err := ResolveWorkspaceArg(context.Background(), c, "65bbffffffffffffffffffff")
	if err == nil || !strings.Contains(err.Error(), "not a team") && !strings.Contains(err.Error(), "no team") {
		t.Fatalf("want membership error, got %v", err)
	}
}

func TestResolveWorkspaceArg_ByName_Unique(t *testing.T) {
	c := &fakeListTeamsClient{teams: []model.Team{
		{ID: "65aa1234567890abcdef1234", Name: "acme", MyRole: ptrRole(model.TeamMemberRoleEditor)},
		{ID: "65bb5678901234abcdef5678", Name: "beta", MyRole: ptrRole(model.TeamMemberRoleViewer)},
	}}
	team, err := ResolveWorkspaceArg(context.Background(), c, "acme")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if team.Name != "acme" {
		t.Fatalf("Name = %q, want acme", team.Name)
	}
}

func TestResolveWorkspaceArg_ByName_NotFound(t *testing.T) {
	c := &fakeListTeamsClient{teams: []model.Team{{ID: "65aa1234567890abcdef1234", Name: "acme"}}}
	_, err := ResolveWorkspaceArg(context.Background(), c, "zeta")
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
	c := &fakeListTeamsClient{teams: []model.Team{
		{ID: id1, Name: "acme", MyRole: ptrRole(model.TeamMemberRoleAdministrator)},
		{ID: id2, Name: "acme", MyRole: ptrRole(model.TeamMemberRoleViewer)},
	}}
	_, err := ResolveWorkspaceArg(context.Background(), c, "acme")
	if err == nil {
		t.Fatal("want ambiguous error, got nil")
	}
	msg := err.Error()
	for _, want := range []string{"ambiguous", "2 workspaces named", id1, id2, "Administrator", "Viewer", "zeabur workspace switch"} {
		if !strings.Contains(msg, want) {
			t.Errorf("error missing %q\nfull: %s", want, msg)
		}
	}
}

// TestIsObjectIDHex sanity-checks the predicate that decides whether to take
// the ID path vs the name path in ResolveWorkspaceArg.
func TestIsObjectIDHex(t *testing.T) {
	cases := []struct {
		s    string
		want bool
	}{
		{"65aa1234567890abcdef1234", true},
		{"65AA1234567890ABCDEF1234", true},
		{"65aa1234567890abcdef123", false},  // 23 chars
		{"65aa1234567890abcdef12345", false}, // 25 chars
		{"65aa1234567890abcdef123g", false},  // non-hex
		{"", false},
		{"acme", false},
	}
	for _, tc := range cases {
		if got := isObjectIDHex(tc.s); got != tc.want {
			t.Errorf("isObjectIDHex(%q) = %v, want %v", tc.s, got, tc.want)
		}
	}
}
