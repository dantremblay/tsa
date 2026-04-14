package access

import (
	"os"
	"strconv"

	"github.com/kassisol/tsa/cli/session"
	"github.com/spf13/cobra"
	"log/slog"
)

func newUnuseCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unuse [session_id]",
		Short: "Unuse session",
		Long:  unuseDescription,
		Run:   runUnuse,
	}

	return cmd
}

func runUnuse(cmd *cobra.Command, args []string) {
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

	if err := sess.Unuse(uint(id)); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

var unuseDescription = `
Unuse session

`
