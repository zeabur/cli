// Package main provides the entry point of the cli
package main

import (
	"github.com/zeabur/cli/internal/cmd/root"
	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/pkg/auth"
	"github.com/zeabur/cli/pkg/config"
	"github.com/zeabur/cli/pkg/log"
	"github.com/zeabur/cli/pkg/printer"
	"github.com/zeabur/cli/pkg/prompt"
	"time"
)

var (
	version = "dev"
	commit  = "none"
	date    = time.Now().Format(time.RFC3339)
)

func main() {
	factory := initFactory()

	rootCmd, err := root.NewCmdRoot(factory, version, commit, date)
	if err != nil {
		panic(err)
	}

	// log errors
	if err := rootCmd.Execute(); err != nil {
		// when some errors occur(such as args dis-match), the log may not be initialized
		if factory.Log == nil {
			factory.Log = log.NewInfoLevel()
		}
		factory.Log.Error(err)
	}
}

// init factory, including config, auth, etc.
func initFactory() *cmdutil.Factory {
	factory := cmdutil.NewFactory()

	configPath, err := config.DefaultConfigFilePath()
	if err != nil {
		panic(err)
	}
	factory.Config = config.New(configPath)

	factory.Printer = printer.New()

	factory.AuthClient = auth.NewZeaburWebAppOAuthClient()

	factory.Prompter = prompt.New()

	return factory
}
