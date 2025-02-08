// Package config provides the config of cli
package config

import (
	"github.com/spf13/viper"

	"github.com/zeabur/cli/pkg/zcontext"
)

// Keys about auth
const (
	KeyTokenString = "token"
	KeyUser        = "user"
	KeyUsername    = "username"
)

// Keys about CLI behavior
const (
	KeyInteractive     = "interactive"
	KeyAutoCheckUpdate = "auto_check_update"
)

type Config interface {
	GetTokenString() string // token string is the single token string, it may be set by user or our login function
	SetTokenString(token string)

	GetUser() string // nickname of user
	SetUser(user string)
	GetUsername() string // it is kind like id of user
	SetUsername(username string)

	GetContext() zcontext.Context

	Write() error
}

type config struct {
	ctx  zcontext.Context
	path string
}

func New(path string) Config {
	// create the config file and init viper
	initViper(path)

	return &config{
		path: path,
		ctx:  zcontext.NewViperContext(viper.GetViper()),
	}
}

func (c *config) GetTokenString() string {
	return viper.GetString(KeyTokenString)
}

func (c *config) SetTokenString(token string) {
	viper.Set(KeyTokenString, token)
}

func (c *config) GetUser() string {
	return viper.GetString(KeyUser)
}

func (c *config) SetUser(user string) {
	viper.Set(KeyUser, user)
}

func (c *config) GetUsername() string {
	return viper.GetString(KeyUsername)
}

func (c *config) SetUsername(username string) {
	viper.Set(KeyUsername, username)
}

func (c *config) GetContext() zcontext.Context {
	return c.ctx
}

func (c *config) Write() error {
	return viper.WriteConfig()
}

var _ Config = &config{}
