package help

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func NewCmdHelp(rootCmd *cobra.Command) *cobra.Command {
	var all bool

	cmd := &cobra.Command{
		Use:   "help [command]",
		Short: "Help about any command",
		RunE: func(cmd *cobra.Command, args []string) error {
			if all {
				printAllCommands(rootCmd, "")
				return nil
			}

			// default: find the target command and show its help
			target, _, err := rootCmd.Find(args)
			if err != nil {
				return err
			}
			return target.Help()
		},
	}

	cmd.Flags().BoolVar(&all, "all", false, "Show all commands with their flags")

	return cmd
}

func printAllCommands(cmd *cobra.Command, prefix string) {
	fullName := prefix + cmd.Name()

	if cmd.Runnable() || len(cmd.Commands()) == 0 {
		fmt.Printf("%s - %s\n", fullName, cmd.Short)
		printFlags(cmd, fullName)
	}

	for _, child := range cmd.Commands() {
		if child.Hidden || child.Name() == "help" {
			continue
		}
		printAllCommands(child, fullName+" ")
	}
}

func printFlags(cmd *cobra.Command, fullName string) {
	var flags []string

	cmd.LocalFlags().VisitAll(func(f *pflag.Flag) {
		if f.Hidden {
			return
		}
		entry := "  --" + f.Name
		if f.Shorthand != "" {
			entry = "  -" + f.Shorthand + ", --" + f.Name
		}
		if f.DefValue != "" && f.DefValue != "false" {
			entry += fmt.Sprintf(" (default: %s)", f.DefValue)
		}
		entry += "  " + f.Usage
		flags = append(flags, entry)
	})

	if len(flags) > 0 {
		fmt.Println(strings.Join(flags, "\n"))
		fmt.Println()
	}
}
