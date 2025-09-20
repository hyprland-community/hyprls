package main

import (
	"fmt"
	"os"
	"path/filepath"

	hyprls "github.com/hyprland-community/hyprls"
	"go.uber.org/zap"
)

var OutputServerLogs string

func main() {
	var logconf zap.Config

	if os.Getenv("HYPRLS_DEBUG") != "" {
		logconf = zap.NewDevelopmentConfig()
	} else {
		logconf = zap.NewProductionConfig()
	}

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
