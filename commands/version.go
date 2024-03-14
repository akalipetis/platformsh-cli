package commands

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/viper"
	"github.com/symfony-cli/console"

	"github.com/platformsh/cli/internal/config"
	"github.com/platformsh/cli/internal/legacy"
)

var (
	version = "0.0.0"
	commit  = "local"
	date    = ""
	builtBy = "local"
)

func newVersionCommand(cnf *config.Config) *console.Command {
	return &console.Command{
		Name:  "version",
		Usage: "Print the version number of the " + cnf.Application.Name,
		Action: func(ctx *console.Context) error {
			fmt.Fprintf(color.Output, "%s %s\n", cnf.Application.Name, color.CyanString(version))

			if viper.GetBool("verbose") {
				fmt.Fprintf(
					color.Output,
					"Embedded PHP version %s\n",
					color.CyanString(legacy.PHPVersion),
				)
				fmt.Fprintf(
					color.Output,
					"Embedded Legacy CLI version %s\n",
					color.CyanString(legacy.LegacyCLIVersion),
				)
				fmt.Fprintf(
					color.Output,
					"Commit %s (built %s by %s)\n",
					color.CyanString(commit),
					color.CyanString(date),
					color.CyanString(builtBy),
				)
			}
			return nil
		},
	}
}
