package access

import (
	"os"

	"github.com/juliengk/go-utils"
	"github.com/juliengk/go-utils/readinput"
	"github.com/kassisol/tsa/cli/session"
	"github.com/spf13/cobra"
	"log/slog"
)

var (
	tsaTTL      int
	tsaUsername string
	tsaPassword string
)

func newLoginCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login [server]",
		Short: "Get TSA access token",
		Long:  loginDescription,
		Run:   runLogin,
	}

	flags := cmd.Flags()

	flags.IntVarP(&tsaTTL, "ttl", "t", 1440, "Token TTL")
	flags.StringVarP(&tsaUsername, "username", "u", "admin", "Username")
	flags.StringVarP(&tsaPassword, "password", "p", "", "Password")

	return cmd
}

func runLogin(cmd *cobra.Command, args []string) {
	var tsattl int
	var username string
	var password string

	go utils.RecoverFunc()

	if len(args) < 1 || len(args) > 1 {
		cmd.Usage()
		os.Exit(-1)
	}

	tsattl = tsaTTL

	if len(tsaUsername) <= 0 {
		username = readinput.ReadInput("Username")
	} else {
		username = tsaUsername
	}

	if len(tsaPassword) <= 0 {
		password = readinput.ReadPassword("Password")
	} else {
		password = tsaPassword
	}

	// Input validations
	// IV - Username
	if len(username) <= 0 {
		slog.Error("Empty username is not allowed")
		os.Exit(1)
	}

	// IV - Password
	if len(password) <= 0 {
		slog.Error("Empty password is not allowed")
		os.Exit(1)
	}

	sess, err := session.New()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	defer sess.End()

	if err := sess.Create(args[0], username, password, tsattl); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

var loginDescription = `
Getting TSA access token.

`
