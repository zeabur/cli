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

	// Two different IDs to keep two different concerns honest:
	//   * effectiveID drives the `*` marker — it must reflect a
	//     `--workspace` flag override so the user sees which workspace
	//     this invocation is acting under.
	//   * persistedID drives the stale-workspace warning — that warning
	//     is about "your saved workspace is no longer valid", which is
	//     a property of the persisted state, not of the override. Using
	//     effectiveID for both would suppress the warning whenever a
	//     valid override happens to be active, hiding a stale persisted
	//     state from the user.
	effectiveID := f.CurrentOwnerID()
	persistedID := ""
	if f.Config != nil {
		persistedID = f.Config.GetContext().GetWorkspace().ID
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	// Personal always renders first. The 24-space placeholder keeps the
	// columns aligned with the team rows that follow.
	personalMarker := " "
	if effectiveID == "" {
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
		if t.ID == effectiveID {
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

	if persistedID != "" {
		// If the persisted workspace is no longer in the membership list
		// (e.g. the team was deleted / the caller was removed), surface it
		// so the user knows the next command without `--workspace` may
		// behave unexpectedly. We deliberately check persistedID, not
		// effectiveID, so a transient `--workspace` override doesn't
		// suppress the warning. We don't auto-clear here either — that's
		// the lazy-verify path in root.
		seen := false
		for _, t := range teams {
			if t.ID == persistedID {
				seen = true
				break
			}
		}
		if !seen {
			fmt.Fprintf(os.Stderr,
				"\nwarning: the persisted workspace %s is not in your memberships any more — run `zeabur workspace clear` or switch to another workspace.\n",
				persistedID,
			)
		}
	}

	return nil
}
