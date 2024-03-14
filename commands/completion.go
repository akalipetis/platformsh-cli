package commands

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/viper"
	"github.com/symfony-cli/console"

	"github.com/platformsh/cli/internal/config"
	"github.com/platformsh/cli/internal/legacy"
)

func completionCommand(cnf *config.Config) *console.Command {
	return &console.Command{
		Name:  "completion",
		Usage: "Print the completion script for your shell",
		Args: console.ArgDefinition{
			&console.Arg{Name: "shell_type", Optional: true},
		},
		Action: func(ctx *console.Context) error {
			completionArgs := []string{"_completion", "-g", "--program", cnf.Application.Executable}
			if shellType := ctx.Args().Get("shell_type"); shellType != "" {
				completionArgs = append(completionArgs, "--shell-type", shellType)
			}
			var b bytes.Buffer
			c := &legacy.CLIWrapper{
				Config:         cnf,
				Version:        version,
				CustomPharPath: viper.GetString("phar-path"),
				Debug:          viper.GetBool("debug"),
				Stdout:         &b,
				Stderr:         os.Stderr,
				Stdin:          os.Stdin,
			}

			if err := c.Init(); err != nil {
				debugLog("%s\n", color.RedString(err.Error()))
				return fmt.Errorf("Cannot initialize Legacy CLI: %w", err)
			}

			if err := c.Exec(context.Background(), completionArgs...); err != nil {
				debugLog("%s\n", color.RedString(err.Error()))
				return fmt.Errorf("Failed to run completion command: %w", err)
			}

			completions := strings.ReplaceAll(
				strings.ReplaceAll(
					b.String(),
					c.PharPath(),
					cnf.Application.Executable,
				),
				path.Base(c.PharPath()),
				cnf.Application.Executable,
			)
			fmt.Fprintln(os.Stdout, "#compdef "+cnf.Application.Executable)
			fmt.Fprintln(os.Stdout, completions)
			return nil
		},
	}
}
