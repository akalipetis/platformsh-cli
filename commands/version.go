package commands

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/symfony-cli/console"
	"github.com/symfony-cli/terminal"

	"github.com/platformsh/cli/internal/config"
	"github.com/platformsh/cli/internal/legacy"
)

var (
	version = "0.0.0"
	channel = "dev"
	date    = ""
	commit  = "local"
	builtBy = "local"
)

func newVersionCommand(cnf *config.Config) *console.Command {
	return &console.Command{
		Name:  "version",
		Usage: "Print the version number of the " + cnf.Application.Name,
		Action: func(_ *console.Context) error {
			fmt.Fprintf(color.Output, "%s %s\n", cnf.Application.Name, color.CyanString(version))
			if terminal.GetLogLevel() > 1 {
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
					color.CyanString(config.Commit),
					color.CyanString(config.Date),
					color.CyanString(config.BuiltBy),
				)
			}
			return nil
		},
	}
}
