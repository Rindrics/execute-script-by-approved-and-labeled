package main

import "github.com/Rindrics/execute-script-with-merge/infrastructure"

func main() {
	logger := infrastructure.NewLogger("debug")

	logger.Info("starting application")
}
