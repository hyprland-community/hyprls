package main

import (
	"fmt"
	"os"

	hyprls "github.com/ewen-lbh/hyprlang-lsp"
	"go.uber.org/zap"
)

func main() {
	logconf := zap.NewDevelopmentConfig()
	// logconf.OutputPaths = []string{"./logs/server.log", "stderr"}
	logger, err := logconf.Build()
	if err != nil {
		fmt.Printf("while building logger: %s", err)
		os.Exit(1)
	}

	logger.Debug("going to start server")
	hyprls.StartServer(logger, "")
}
