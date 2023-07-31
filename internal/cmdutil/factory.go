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
	}
	// PersistentFlags are flags that are common to all commands
	PersistentFlags struct {
		Debug            bool // debug mode, default false
		Interactive      bool // interactive mode, default true
		AutoRefreshToken bool // auto refresh token, default true, only when token is from browser(OAuth2)
	}
)

// NewFactory returns a new cmd factory
func NewFactory() *Factory {
	return &Factory{}
}
