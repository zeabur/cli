package cmdutil

import (
	"context"

	"github.com/zeabur/cli/pkg/fill"
	"github.com/zeabur/cli/pkg/printer"
	"github.com/zeabur/cli/pkg/selector"
	"go.uber.org/zap"

	"github.com/zeabur/cli/pkg/api"
	"github.com/zeabur/cli/pkg/auth"
	"github.com/zeabur/cli/pkg/config"
	"github.com/zeabur/cli/pkg/model"
	"github.com/zeabur/cli/pkg/prompt"
	"github.com/zeabur/cli/pkg/zcontext"
)

type (
	// Factory is a factory for command runners
	// It is used to pass common dependencies to commands.
	// It is kind of like a "context" for commands.
	Factory struct {
		Log         *zap.SugaredLogger // logger
		Printer     printer.Printer    // printer
		Config      config.Config      // config(flag, env, file)
		ApiClient   api.Client         // query api
		AuthClient  auth.Client        // login, refresh token
		Prompter    prompt.Prompter    // interactive prompter
		Selector    selector.Selector  // interactive selector
		ParamFiller fill.ParamFiller   // fill params
		PersistentFlags

		// workspaceOverride is the team resolved from the --workspace flag
		// during PersistentPreRunE. Nil when the flag is not set;
		// CurrentOwnerID / CurrentWorkspace then fall back to the persisted
		// workspace. Stored as the full Workspace (not just an ID) so
		// downstream code that wants to show a name — e.g. the
		// "creating new project in team workspace X" hint in deploy — gets
		// the same effective workspace as CurrentOwnerID.
		workspaceOverride *zcontext.Workspace

		// teamsCache memoizes the per-process ListTeams reply. A single CLI
		// invocation can otherwise hit `teams` up to three times (flag
		// resolution, persisted-workspace verify, the command itself); see
		// PLA-1590 review feedback.
		teamsCache    []model.Team
		teamsCacheErr error
		teamsCacheHit bool
	}
	// PersistentFlags are flags that are common to all commands
	PersistentFlags struct {
		Debug            bool   // debug mode, default false
		Interactive      bool   // interactive mode, default true
		AutoRefreshToken bool   // auto refresh token, default true, only when token is from browser(OAuth2)
		AutoCheckUpdate  bool   // auto check update, default true
		JSON             bool   // output in JSON format, default false
		Workspace        string // --workspace <name|id> one-shot override
	}
)

// CurrentOwnerID returns the team ObjectID hex that directory-level commands
// (project list / create / deploy-no-link) should act under. Empty string ==
// the caller's personal account.
//
// Resolution order:
//  1. --workspace flag (resolved to a Workspace during PersistentPreRunE)
//  2. Persisted workspace in the config file
//
// Returning an empty string is the canonical "personal" signal that the
// project API uses to fall back to the un-owner-scoped GraphQL query.
func (f *Factory) CurrentOwnerID() string {
	return f.CurrentWorkspace().ID
}

// CurrentWorkspace returns the effective workspace under the same resolution
// rules as CurrentOwnerID, including the name and kind. Callers that want to
// display the active workspace (the "creating new project in team workspace
// X" hint in deploy) should read this rather than the persisted workspace,
// so a --workspace override shows up correctly.
func (f *Factory) CurrentWorkspace() *zcontext.Workspace {
	if f.workspaceOverride != nil {
		return f.workspaceOverride
	}
	if f.Config == nil {
		return &zcontext.Workspace{}
	}
	return f.Config.GetContext().GetWorkspace()
}

// SetWorkspaceOverride records the resolved workspace for a --workspace flag
// value. Called from PersistentPreRunE after the flag string has been
// disambiguated against the list of teams. Passing nil clears any prior
// override.
func (f *Factory) SetWorkspaceOverride(ws *zcontext.Workspace) {
	f.workspaceOverride = ws
}

// ListTeams returns the caller's teams via api.Client.ListTeams, memoized for
// the lifetime of this Factory. The same Factory is shared across every
// PersistentPreRunE / Run / PersistentPostRunE callback within a single CLI
// invocation, so one fetch covers --workspace flag resolution, the lazy
// startup verify, and downstream commands that surface roles.
//
// Errors are sticky: a failed fetch is cached and returned on every
// subsequent call, so callers don't accidentally retry against a broken
// backend within the same process.
func (f *Factory) ListTeams(ctx context.Context) ([]model.Team, error) {
	if f.teamsCacheHit {
		return f.teamsCache, f.teamsCacheErr
	}
	f.teamsCache, f.teamsCacheErr = f.ApiClient.ListTeams(ctx)
	f.teamsCacheHit = true
	return f.teamsCache, f.teamsCacheErr
}

// NewFactory returns a new cmd factory
func NewFactory() *Factory {
	return &Factory{}
}
