package auth

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/kassisol/tsa/cli/session"
	"github.com/kassisol/tsa/client"
	"github.com/spf13/cobra"
	"log/slog"
)

func newListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls",
		Aliases: []string{"list"},
		Short:   "List authentication configurations",
		Long:    listDescription,
		Run:     runList,
	}

	return cmd
}

func runList(cmd *cobra.Command, args []string) {
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

	configs, err := clt.AuthList(srv.Token)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	if len(configs) > 0 {
		w := tabwriter.NewWriter(os.Stdout, 20, 1, 2, ' ', 0)
		fmt.Fprintln(w, "KEY\tVALUE")

		for _, config := range configs {
			fmt.Fprintf(w, "%s\t%s\n", config.Key, config.Value)
		}

		w.Flush()
	}
}

var listDescription = `
List authentication configurations

`
