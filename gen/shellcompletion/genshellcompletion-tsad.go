package main

import (
	"fmt"
	"os"

	"github.com/juliengk/go-utils/filedir"
	"github.com/kassisol/tsa/daemon"
)

func main() {
	scPath := os.Getenv("TSAD_SC_PATH")
	if scPath == "" {
		scPath = "dist/tsad/shellcompletion"
	}
	bashTarget := fmt.Sprintf("%s/bash", scPath)

	if err := filedir.CreateDirIfNotExist(scPath, true, 0755); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cmd := daemon.NewCommand()
	cmd.DisableAutoGenTag = true

	if err := cmd.GenBashCompletionFile(bashTarget); err != nil {
		fmt.Println(err)
	}
}
