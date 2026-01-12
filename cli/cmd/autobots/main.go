package main

import (
	"os"
	
	"github.com/SheaHawkins/AutoBots/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
