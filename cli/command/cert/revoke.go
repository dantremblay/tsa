package cert

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/kassisol/tsa/cli/session"
	"github.com/kassisol/tsa/client"
	"github.com/spf13/cobra"
)

var revokeCN string

func newRevokeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "revoke [serial number]",
		Short: "Revoke certificate",
		Long:  revokeDescription,
		Run:   runRevoke,
	}

	flags := cmd.Flags()
	flags.StringVar(&revokeCN, "cn", "", "Revoke all valid certificates matching this CN (server name)")

	return cmd
}

func runRevoke(cmd *cobra.Command, args []string) {
	if len(revokeCN) == 0 && len(args) != 1 {
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

	if len(revokeCN) > 0 {
		if err := clt.CertRevokeByCN(srv.Token, revokeCN); err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
	} else {
		serialNumber, err := strconv.Atoi(args[0])
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}

		if err := clt.CertRevoke(srv.Token, serialNumber); err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
	}
}

var revokeDescription = `
Revoke certificate

`
