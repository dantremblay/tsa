package auth

import (
	"os"

	"github.com/kassisol/tsa/cli/session"
	"github.com/kassisol/tsa/client"
	"github.com/spf13/cobra"
	"log/slog"
)

func newDisableCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable",
		Short: "Disable authentication",
		Long:  disableDescription,
		Run:   runDisable,
	}

	return cmd
}

func runDisable(cmd *cobra.Command, args []string) {
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

	if err := clt.AuthDisable(srv.Token); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

var disableDescription = `
Disable authentication

`
