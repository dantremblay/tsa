package server

import (
	"os"

	"github.com/juliengk/go-utils/validation"
	"github.com/juliengk/stack/client"
	"github.com/kassisol/tsa/cli/storage"
	"github.com/kassisol/tsa/pkg/adf"
	"github.com/spf13/cobra"
	"log/slog"
)

func newAddCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [name] [tsa url]",
		Short: "Add TSA server",
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

	cfg := adf.NewServer()
	if err := cfg.Init(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	s, err := storage.NewDriver("sqlite", cfg.AppDir)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	defer s.End()

	// Input Validations
	// IV - Server name
	if err = validation.IsValidName(args[0]); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	// IV - TSA URL
	if _, err := client.ParseUrl(args[1]); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	s.AddServer(args[0], args[1], "")
}

var addDescription = `
Add TSA server

`
