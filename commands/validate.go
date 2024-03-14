package commands

import (
	"fmt"
	"os"

	"github.com/platformsh/platformify/validator"
	"github.com/platformsh/platformify/vendorization"
	"github.com/symfony-cli/console"
)

func validateCommand(assets *vendorization.VendorAssets) *console.Command {
	return &console.Command{
		Name:        "config-validate",
		Aliases:     []*console.Alias{{Name: "validate"}},
		Description: fmt.Sprintf("%s app:config-validate [options]", assets.Binary),
		Usage:       "Validate the YAML files of your project",
		Category:    "app",
		Action: func(ctx *console.Context) error {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("Cannot get current working directory: %w", err)
			}

			if err := validator.ValidateConfig(cwd, assets.ConfigFlavor); err != nil {
				return fmt.Errorf("Invalid %s config: %w", assets.ServiceName, err)
			}

			return nil
		},
	}
}
