package access

import (
	"os"
	"strconv"

	"github.com/kassisol/tsa/cli/session"
	"github.com/spf13/cobra"
	"log/slog"
)

var sessionRemoveAll bool

func newRemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rm [session_id]",
		Aliases: []string{"remove"},
		Short:   "Remove session",
		Long:    removeDescription,
		Run:     runRemove,
	}

	flags := cmd.Flags()
	flags.BoolVarP(&sessionRemoveAll, "all", "a", false, "Remove all sessions")

	return cmd
}

func runRemove(cmd *cobra.Command, args []string) {
	if !sessionRemoveAll && (len(args) < 1 || len(args) > 1) {
		cmd.Usage()
		os.Exit(-1)
	}

	sess, err := session.New()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	defer sess.End()

	if sessionRemoveAll {
		if err := sess.Clear(); err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
	} else {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}

		if err := sess.Remove(uint(id)); err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
	}
}

var removeDescription = `
Remove session

`
