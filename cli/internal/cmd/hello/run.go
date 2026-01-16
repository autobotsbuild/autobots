package hello

import (
	"context"
	"fmt"

	"github.com/SheaHawkins/AutoBots/internal/cmd/shared"
)

func Run(ctx context.Context, deps shared.Dependencies, flags *Flags, args []string) error {
	name := "World"
	if len(args) > 0 {
		name = args[0]
	}

	if flags.Verbose {
		fmt.Println("Running in verbose mode...")
	}

	fmt.Printf("Hello, %s! Welcome to AutoBots.\n", name)
	return nil
}
