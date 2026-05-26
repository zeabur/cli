package zcontext

// WorkspaceKindTeam marks a workspace as a Zeabur team. The empty string is
// "personal" — the user's own resources, the default when the CLI has never
// been switched to a team.
const WorkspaceKindTeam = "team"

// Workspace identifies which owner the CLI is currently acting against for
// directory-level commands (project list, project create, deploy with no
// linked project). Other commands address resources by their own ID and stay
// workspace-independent.
type Workspace struct {
	// ID is the team's MongoDB ObjectID hex string. Empty for personal.
	ID string
	// Name is the team's display name. Empty for personal.
	Name string
	// Kind is "team" for team workspaces, empty for personal.
	Kind string
}

// IsPersonal reports whether the workspace addresses the caller's own
// resources (no team scope). The zero value is personal.
func (w *Workspace) IsPersonal() bool {
	return w == nil || w.ID == ""
}

// IsTeam reports whether the workspace addresses a team.
func (w *Workspace) IsTeam() bool {
	return w != nil && w.Kind == WorkspaceKindTeam && w.ID != ""
}
