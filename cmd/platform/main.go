package main

import (
	"log"
	"os"

	"github.com/fatih/color"

	"github.com/platformsh/cli/commands"
	"github.com/platformsh/cli/internal/config"
)

func main() {
	log.SetOutput(color.Error)

	// Load configuration.
	cnfYAML, err := config.LoadYAML()
	if err != nil {
		log.Fatal(err)
	}
	cnf, err := config.FromYAML(cnfYAML)
	if err != nil {
		log.Fatal(err)
	}

	if err := commands.Execute(cnf); err != nil {
		os.Exit(1)
	}
}
