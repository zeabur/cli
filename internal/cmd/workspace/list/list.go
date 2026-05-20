// Package list implements `zeabur workspace list`.
package list

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
)

// NewCmdList builds `zeabur workspace list`.
func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "List the personal workspace and every team the caller belongs to",
		Aliases: []string{"ls"},
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(f)
		},
	}
}

func run(f *cmdutil.Factory) error {
	teams, err := f.ListTeams(context.Background())
	if err != nil {
		return fmt.Errorf("list teams: %w", err)
	}

	// Use the effective workspace so the `*` marker tracks a `--workspace`
	// flag override, not just the persisted state. Otherwise
	// `--workspace foo workspace list` would print `*` on the persisted
	// team rather than `foo`.
	currentID := f.CurrentOwnerID()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	// Personal always renders first. The 24-space placeholder keeps the
	// columns aligned with the team rows that follow.
	personalMarker := " "
	if currentID == "" {
		personalMarker = "*"
	}
	personalLabel := f.Config.GetUser()
	if personalLabel == "" {
		personalLabel = f.Config.GetUsername()
	}
	if personalLabel == "" {
		personalLabel = "(you)"
	}
	fmt.Fprintf(w, "%s\tpersonal\t\t\t(%s)\n", personalMarker, personalLabel)

	for _, t := range teams {
		marker := " "
		if t.ID == currentID {
			marker = "*"
		}
		role := ""
		if t.MyRole != nil {
			role = t.MyRole.Display()
		}
		fmt.Fprintf(w, "%s\t%s\t%s\tteam\t%s\n", marker, t.ID, t.Name, role)
	}

	if err := w.Flush(); err != nil {
		return fmt.Errorf("render table: %w", err)
	}

	if currentID != "" {
		// If the persisted workspace is no longer in the membership list
		// (e.g. the team was deleted / the caller was removed), surface it
		// so the user knows the next command may behave unexpectedly. We
		// don't auto-clear here — that's the lazy-verify path in root.
		seen := false
		for _, t := range teams {
			if t.ID == currentID {
				seen = true
				break
			}
		}
		if !seen {
			fmt.Fprintf(os.Stderr,
				"\nwarning: the persisted workspace %s is not in your memberships any more — run `zeabur workspace clear` or switch to another workspace.\n",
				currentID,
			)
		}
	}

	return nil
}
