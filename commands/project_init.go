package commands

import (
	"context"
	"fmt"

	"github.com/platformsh/platformify/commands"
	"github.com/platformsh/platformify/vendorization"
	"github.com/symfony-cli/console"
)

func projectInitCommand(assets *vendorization.VendorAssets) *console.Command {
	return &console.Command{
		Name:        "init",
		Aliases:     []*console.Alias{{Name: "ify"}},
		Usage:       fmt.Sprintf("%s project:init [options]", assets.Binary),
		Description: "Create the starter YAML files for your project",
		Category:    "project",
		Action: func(ctx *console.Context) error {
			return commands.Platformify(context.Background(), assets)
		},
	}
}
