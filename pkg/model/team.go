package model

// TeamMemberRole mirrors the backend GraphQL TeamMemberRole enum. Values are
// the raw enum names so they print identically to the backend / dashboard.
type TeamMemberRole string

const (
	TeamMemberRoleAdministrator TeamMemberRole = "ADMINISTRATOR"
	TeamMemberRoleEditor        TeamMemberRole = "EDITOR"
	TeamMemberRoleViewer        TeamMemberRole = "VIEWER"
)

// Display returns the role spelled the way the dashboard shows it
// ("Administrator" / "Editor" / "Viewer") rather than the SCREAMING form.
func (r TeamMemberRole) Display() string {
	switch r {
	case TeamMemberRoleAdministrator:
		return "Administrator"
	case TeamMemberRoleEditor:
		return "Editor"
	case TeamMemberRoleViewer:
		return "Viewer"
	default:
		return string(r)
	}
}

// Team is the slim shape returned by the `teams` query — enough to drive
// `zeabur workspace list / switch`. `MyRole` is the caller's own role and
// comes from backend PLA-1589 (Team.myRole field). It's a pointer because the
// backend marks it nullable; for `teams` results it is always set.
type Team struct {
	ID     string          `graphql:"_id"`
	Name   string          `graphql:"name"`
	MyRole *TeamMemberRole `graphql:"myRole"`
}
