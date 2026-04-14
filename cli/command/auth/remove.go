package auth

import (
	"os"

	"github.com/kassisol/tsa/cli/session"
	"github.com/kassisol/tsa/client"
	"github.com/spf13/cobra"
	"log/slog"
)

func newRemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rm [key] [value]",
		Aliases: []string{"remove"},
		Short:   "Remove authentication configuration",
		Long:    removeDescription,
		Run:     runRemove,
	}

	return cmd
}

func runRemove(cmd *cobra.Command, args []string) {
	if len(args) < 1 || len(args) > 2 {
		cmd.Usage()
		os.Exit(-1)
	}

	sess, err := session.New()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	defer sess.End()

	srv, err := sess.Get()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	clt, err := client.New(srv.Server.TSAURL)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	if len(args) == 1 {
		args = append(args, "")
	}

	if err := clt.AuthDelete(srv.Token, args[0], args[1]); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

var removeDescription = `
Remove authentication configuration

`
