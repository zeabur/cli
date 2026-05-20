package zcontext

// Context represents the current context of the CLI, including the current
// workspace (personal or team), and the pinned project / environment / service
// inside that workspace.
type Context interface {
	// GetWorkspace returns the persisted workspace. Always non-nil; check
	// IsPersonal() / IsTeam() to branch. Personal is the zero value.
	GetWorkspace() *Workspace
	// SetWorkspace persists the workspace; passing nil is equivalent to
	// ClearWorkspace().
	SetWorkspace(workspace *Workspace)
	// ClearWorkspace returns to personal — the default zero state.
	ClearWorkspace()

	GetProject() BasicInfo
	SetProject(project BasicInfo)
	ClearProject()

	GetEnvironment() BasicInfo
	SetEnvironment(environment BasicInfo)
	ClearEnvironment()

	GetService() BasicInfo
	SetService(service BasicInfo)
	ClearService()

	// ClearAll clears the inner project / environment / service context but
	// leaves the workspace intact. Use ClearWorkspace() to go back to personal,
	// which on its own does not touch the inner context — switching commands
	// orchestrate that explicitly so they can report what was cleared.
	ClearAll()
}

// BasicInfo represents the basic information of a resource.
type BasicInfo interface {
	GetID() string
	GetName() string
	Empty() bool
}
