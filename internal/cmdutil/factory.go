package cmdutil

import (
	"github.com/zeabur/cli/pkg/fill"
	"github.com/zeabur/cli/pkg/printer"
	"github.com/zeabur/cli/pkg/selector"
	"go.uber.org/zap"

	"github.com/zeabur/cli/pkg/api"
	"github.com/zeabur/cli/pkg/auth"
	"github.com/zeabur/cli/pkg/config"
	"github.com/zeabur/cli/pkg/prompt"
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

		// workspaceOverrideID is the team ObjectID resolved from the
		// --workspace flag during PersistentPreRunE. Empty when the flag is
		// not set; CurrentOwnerID then falls back to the persisted workspace.
		workspaceOverrideID string
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
//  1. --workspace flag (resolved to ID during PersistentPreRunE)
//  2. Persisted workspace in the config file
//
// Returning an empty string is the canonical "personal" signal that the
// project API uses to fall back to the un-owner-scoped GraphQL query.
func (f *Factory) CurrentOwnerID() string {
	if f.workspaceOverrideID != "" {
		return f.workspaceOverrideID
	}
	if f.Config == nil {
		return ""
	}
	return f.Config.GetContext().GetWorkspace().ID
}

// SetWorkspaceOverride records the resolved ID for a --workspace flag value.
// Called from PersistentPreRunE after the flag string has been disambiguated
// against the list of teams.
func (f *Factory) SetWorkspaceOverride(id string) {
	f.workspaceOverrideID = id
}

// NewFactory returns a new cmd factory
func NewFactory() *Factory {
	return &Factory{}
}
