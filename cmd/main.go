package main

import (
	"fmt"
	"os"
	"path/filepath"

	hyprls "github.com/ewen-lbh/hyprlang-lsp"
	"go.uber.org/zap"
)

var OutputServerLogs string

func main() {
	logconf := zap.NewDevelopmentConfig()
	if OutputServerLogs != "" {
		logconf.OutputPaths = []string{OutputServerLogs, "stderr"}
	}
	logger, err := logconf.Build()
	if err != nil {
		fmt.Printf("while building logger: %s", err)
		os.Exit(1)
	}

	logger.Debug("going to start server")
	if OutputServerLogs != "" {
		hyprls.StartServer(logger, filepath.Dir(OutputServerLogs))
	} else {
		hyprls.StartServer(logger, "")
	}
}
