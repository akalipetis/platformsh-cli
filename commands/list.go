package commands

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/platformsh/cli/internal/config"
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
	c := makeLegacyCLIWrapper(cnf, &b, nil, nil)

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
