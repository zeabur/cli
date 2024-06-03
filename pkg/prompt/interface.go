// Package prompt provides a prompter interface for UI interactions
package prompt

// Prompter is the interface for prompting the user for input
type Prompter interface {
	Select(message string, defaultValue string, options []string) (int, error)
	MultiSelect(message string, defaultValues, options []string) ([]int, error)
	Input(prompt, defaultValue string) (string, error)
	InputWithHelp(prompt, help, defaultValue string) (string, error)
	Confirm(prompt string, defaultValue bool) (bool, error)
	ConfirmDeletion(requiredValue string) error
}
