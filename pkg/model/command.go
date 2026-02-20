package model

// CommandResult represents the result of a command execution.
type CommandResult struct {
	ExitCode int    `graphql:"exitCode"`
	Output   string `graphql:"output"`
}
