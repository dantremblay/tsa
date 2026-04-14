package access

import (
	"fmt"
	"os"

	"github.com/kassisol/tsa/cli/session"
	"github.com/spf13/cobra"
	"log/slog"
)

func newStatusCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Session status",
		Long:  statusDescription,
		Run:   runStatus,
	}

	return cmd
}

func runStatus(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		cmd.Usage()
		os.Exit(-1)
	}

	sess, err := session.New()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	defer sess.End()

	fmt.Println(sess.Status())
}

var statusDescription = `
Session status

`
