package zcontext

// EphemeralContext is an in-memory Context used under `--workspace` override
// (PLA-1590 B+). It exists so the CLI can run an interactive command that
// transiently selects a project / service / environment without ever writing
// those choices back to the persisted config — that would silently pin
// resources from the override workspace under the persisted workspace's name
// and cause cross-workspace operations on later commands.
//
// Reads start empty (no implicit fallback to persisted state) and writes go
// to in-memory fields only. ParamFiller's "Set then re-read within the same
// call" pattern still works, because reads see whatever was Set during this
// process; the values just don't survive to the next command.
//
// `GetWorkspace()` returns the workspace the caller passed at construction —
// the override workspace, not personal — so any consumer that asks "what
// workspace is this context for?" gets the right answer. That avoids a
// second-order trap where some future helper reads `ctx.GetWorkspace()`
// expecting it to match `Factory.CurrentWorkspace()` and finds a personal
// reading instead.
type ephemeralContext struct {
	workspace   *Workspace
	project     BasicInfo
	environment BasicInfo
	service     BasicInfo
}

// NewEphemeralContext returns a new in-memory Context whose GetWorkspace()
// reports the supplied workspace. Pass `nil` for personal (or when no
// override is active and the caller really wants a blank scratch context —
// uncommon).
func NewEphemeralContext(workspace *Workspace) Context {
	if workspace == nil {
		workspace = &Workspace{}
	}
	return &ephemeralContext{
		workspace:   workspace,
		project:     &basicInfo{},
		environment: &basicInfo{},
		service:     &basicInfo{},
	}
}

func (c *ephemeralContext) GetWorkspace() *Workspace {
	return c.workspace
}

func (c *ephemeralContext) SetWorkspace(workspace *Workspace) {
	if workspace == nil {
		c.workspace = &Workspace{}
		return
	}
	c.workspace = workspace
}

func (c *ephemeralContext) ClearWorkspace() {
	c.workspace = &Workspace{}
}

func (c *ephemeralContext) GetProject() BasicInfo { return c.project }

func (c *ephemeralContext) SetProject(project BasicInfo) {
	if project == nil {
		c.project = &basicInfo{}
		return
	}
	c.project = project
}

func (c *ephemeralContext) ClearProject() { c.project = &basicInfo{} }

func (c *ephemeralContext) GetEnvironment() BasicInfo { return c.environment }

func (c *ephemeralContext) SetEnvironment(environment BasicInfo) {
	if environment == nil {
		c.environment = &basicInfo{}
		return
	}
	c.environment = environment
}

func (c *ephemeralContext) ClearEnvironment() { c.environment = &basicInfo{} }

func (c *ephemeralContext) GetService() BasicInfo { return c.service }

func (c *ephemeralContext) SetService(service BasicInfo) {
	if service == nil {
		c.service = &basicInfo{}
		return
	}
	c.service = service
}

func (c *ephemeralContext) ClearService() { c.service = &basicInfo{} }

func (c *ephemeralContext) ClearAll() {
	c.ClearProject()
	c.ClearEnvironment()
	c.ClearService()
}

var _ Context = &ephemeralContext{}
