package system

import (
	"fmt"
	"os"

	"github.com/juliengk/go-utils/readinput"
	"github.com/kassisol/tsa/cli/session"
	"github.com/kassisol/tsa/client"
	"github.com/spf13/cobra"
	"log/slog"
)

func NewPasswdCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "passwd [server name]",
		Short: "Change admin password",
		Long:  passwdDescription,
		Run:   runPasswd,
	}

	return cmd
}

func runPasswd(cmd *cobra.Command, args []string) {
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

	srv, err := sess.GetServer(args[0])
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	oldPassword := readinput.ReadPassword("Old Password")
	newPassword := readinput.ReadPassword("New Password")
	confirmPassword := readinput.ReadPassword("Confirm Password")

	clt, err := client.New(srv.TSAURL)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	if err := clt.AdminChangePassword(oldPassword, newPassword, confirmPassword); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	fmt.Println("Password changed successfully")
}

var passwdDescription = `
Change user password

`
