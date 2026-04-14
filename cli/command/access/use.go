package access

import (
	"os"
	"strconv"

	"github.com/kassisol/tsa/cli/session"
	"github.com/spf13/cobra"
	"log/slog"
)

func newUseCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "use [session_id]",
		Short: "Use session",
		Long:  useDescription,
		Run:   runUse,
	}

	return cmd
}

func runUse(cmd *cobra.Command, args []string) {
	if len(args) < 1 || len(args) > 1 {
		cmd.Usage()
		os.Exit(-1)
	}

	sess, err := session.New()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	defer sess.End()

	id, err := strconv.Atoi(args[0])
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	if err := sess.Use(uint(id)); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

var useDescription = `
Use session

`
