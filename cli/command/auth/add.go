package auth

import (
	"os"

	"github.com/kassisol/tsa/cli/session"
	"github.com/kassisol/tsa/client"
	"github.com/spf13/cobra"
	"log/slog"
)

func newAddCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [key] [value]",
		Short: "Add auth configuration",
		Long:  addDescription,
		Run:   runAdd,
	}

	return cmd
}

func runAdd(cmd *cobra.Command, args []string) {
	if len(args) < 2 || len(args) > 2 {
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

	if _, err := clt.AuthCreate(srv.Token, args[0], args[1]); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

var addDescription = `
Add auth configuration

`
