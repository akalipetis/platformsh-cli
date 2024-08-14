package main

import (
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/viper"

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

	viper.SetEnvPrefix(strings.TrimSuffix(cnf.Application.EnvPrefix, "_"))
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	if err := commands.Execute(cnf); err != nil {
		os.Exit(1)
	}
}
