package main

import (
	"fmt"
	"os"

	"github.com/juliengk/go-utils/filedir"
	"github.com/kassisol/tsa/cli/command/commands"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func main() {
	manPath := os.Getenv("TSA_MAN_PATH")
	if manPath == "" {
		manPath = "dist/tsa/man"
	}
	man8 := fmt.Sprintf("%s/man8", manPath)

	if err := filedir.CreateDirIfNotExist(man8, true, 0755); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	header := &doc.GenManHeader{
		Title:   "TSA",
		Section: "8",
		Source:  "Harbormaster",
	}
	opts := doc.GenManTreeOptions{
		Header:           header,
		Path:             man8,
		CommandSeparator: "-",
	}

	cmd := &cobra.Command{Use: "tsa"}
	commands.AddCommands(cmd)
	cmd.DisableAutoGenTag = true

	if err := doc.GenManTreeFromOpts(cmd, opts); err != nil {
		fmt.Println(err)
	}
}
