package commands

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/viper"

	"github.com/platformsh/cli/internal/config"
	"github.com/platformsh/cli/internal/legacy"
)

func listLegacyCommands(ctx context.Context, cnf *config.Config, category string, all bool) (*List, error) {
	arguments := []string{"list", "--format=json"}
	if all {
		arguments = append(arguments, "--all")
	}
	if category != "" {
		arguments = append(arguments, category)
	}

	var b bytes.Buffer
	c := &legacy.CLIWrapper{
		Config:         cnf,
		Version:        version,
		CustomPharPath: viper.GetString("phar-path"),
		Debug:          viper.GetBool("debug"),
		Stdout:         &b,
		Stderr:         nil,
		Stdin:          nil,
	}

	if err := c.Init(); err != nil {
		return nil, fmt.Errorf("could not initialize legacy CLI: %w", err)
	}

	if err := c.Exec(ctx, arguments...); err != nil {
		return nil, fmt.Errorf("could not list legacy CLI commands: %w", err)
	}

	list := &List{}
	if err := json.Unmarshal(b.Bytes(), list); err != nil {
		return nil, fmt.Errorf("could not parse legacy CLI command list: %w", err)
	}

	list.Application.Name = cnf.Application.Name
	list.Application.Executable = cnf.Application.Executable

	return list, nil
}
