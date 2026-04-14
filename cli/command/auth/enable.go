package auth

import (
	"os"

	"github.com/kassisol/tsa/cli/session"
	"github.com/kassisol/tsa/client"
	"github.com/spf13/cobra"
	"log/slog"
)

func newEnableCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enable [type]",
		Short: "Enable authentication",
		Long:  enableDescription,
		Run:   runEnable,
	}

	return cmd
}

func runEnable(cmd *cobra.Command, args []string) {
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

	if err := clt.AuthEnable(srv.Token, args[0]); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

var enableDescription = `
Enable authentication

`
