package main

import (
	"github.com/zeabur/cli/internal/cmd/root"
	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/pkg/auth"
	"github.com/zeabur/cli/pkg/config"
	"github.com/zeabur/cli/pkg/prompt"
)

func main() {
	factory := initFactory()

	// todo: get version from build
	version := "0.0.1"
	rootCmd, err := root.NewCmdRoot(factory, version)
	if err != nil {
		panic(err)
	}

	// cobra will handle the error
	_ = rootCmd.Execute()
}

// init factory, including config, auth, etc.
func initFactory() *cmdutil.Factory {
	factory := cmdutil.NewFactory()

	configPath, err := config.DefaultConfigFilePath()
	if err != nil {
		panic(err)
	}
	factory.Config = config.New(configPath)

	factory.AuthClient = auth.NewZeaburWebAppOAuthClient()

	factory.Prompter = prompt.New()

	return factory
}
