package main

import (
	"fmt"
	"os"

	"github.com/juliengk/go-utils/filedir"
	"github.com/kassisol/tsa/cli/command/commands"
	"github.com/spf13/cobra"
)

func main() {
	scPath := os.Getenv("TSA_SC_PATH")
	if scPath == "" {
		scPath = "dist/tsa/shellcompletion"
	}
	bashTarget := fmt.Sprintf("%s/bash", scPath)

	if err := filedir.CreateDirIfNotExist(scPath, true, 0755); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cmd := &cobra.Command{Use: "tsa"}
	commands.AddCommands(cmd)
	cmd.DisableAutoGenTag = true

	if err := cmd.GenBashCompletionFile(bashTarget); err != nil {
		fmt.Println(err)
	}
}
