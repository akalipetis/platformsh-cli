package commands

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strings"
	"text/tabwriter"
	"text/template"

	"github.com/platformsh/platformify/vendorization"
	"github.com/spf13/viper"
	"github.com/symfony-cli/console"

	"github.com/platformsh/cli/internal/config"
	"github.com/platformsh/cli/internal/legacy"
)

var overrideCommands = []string{"list", "help"}

// Execute is the main entrypoint to run the CLI.
func Execute(cnf *config.Config) error {
	assets := &vendorization.VendorAssets{
		Use:          "project:init",
		Binary:       cnf.Application.Executable,
		ConfigFlavor: cnf.Service.ProjectConfigFlavor,
		EnvPrefix:    strings.TrimSuffix(cnf.Service.EnvPrefix, "_"),
		ServiceName:  cnf.Service.Name,
		DocsBaseURL:  cnf.Service.DocsURL,
	}
	console.CommandHelpTemplate = `<comment>Command:</> {{.FullName}}
{{if .Aliases}}<comment>Aliases:</> {{joinAliases .Aliases ", "}}
{{end}}{{if .Usage}}<comment>Description:</> {{.Usage}}

{{end}}<comment>Usage:</>
  {{.HelpName}}{{if .VisibleFlags}} [options]{{end}}{{.Arguments.Usage}}{{if .Arguments}}

<comment>Arguments:</>
  {{range .Arguments}}{{.}}
  {{end}}{{end}}{{if .VisibleFlags}}

<comment>Options:</>
  {{range .VisibleFlags}}{{.}}
  {{end}}{{end}}{{if .Description}}

<comment>Help:</>
  {{.Description}}
{{end}}
`
	console.HelpPrinter = func(out io.Writer, templ string, data interface{}) {
		funcMap := template.FuncMap{
			"join": strings.Join,
			"joinAliases": func(aliases []*console.Alias, sep string) string {
				aliasesStr := make([]string, 0, len(aliases))
				for _, a := range aliases {
					aliasesStr = append(aliasesStr, a.String())
				}

				return strings.Join(aliasesStr, sep)
			},
		}

		w := tabwriter.NewWriter(out, 1, 8, 2, ' ', 0)
		t := template.Must(template.New("help").Funcs(funcMap).Parse(templ))

		err := t.Execute(w, data)
		if err != nil {
			panic(fmt.Errorf("CLI TEMPLATE ERROR: %#v", err.Error()))
		}
		w.Flush()
	}

	// ctx := vendorization.WithVendorAssets(config.ToContext(context.Background(), cnf), assets)
	app, err := newApp(cnf, assets)
	if err != nil {
		return err
	}
	return app.Run(os.Args)
}

func newApp(cnf *config.Config, assets *vendorization.VendorAssets) (*console.Application, error) {
	legacyAction := func(ctx *console.Context) error {
		var legacyErr error
		defer func() {
			if legacyErr != nil {
			}
		}()

		c := &legacy.CLIWrapper{
			Config:         cnf,
			Version:        version,
			CustomPharPath: viper.GetString("phar-path"),
			Debug:          viper.GetBool("debug"),
			Stdout:         os.Stdout,
			Stderr:         os.Stderr,
			Stdin:          os.Stdin,
		}
		if legacyErr = c.Init(); legacyErr != nil {
			debugLog("failed to initialize legacy CLI: %s", legacyErr)
			return nil
		}

		if legacyErr := c.Exec(context.TODO(), os.Args[1:]...); legacyErr != nil {
			debugLog("failed to run legacy CLI command: %s", legacyErr)
			return nil
		}

		return nil
	}

	list, err := listLegacyCommands(context.TODO(), cnf, "", true)
	if err != nil {
		return nil, err
	}

	console.VersionFlag = &console.BoolFlag{
		Name:    "version",
		Aliases: []string{"V"},
		Usage:   console.VersionFlag.Usage,
	}
	globalFlags := []console.Flag{
		console.VersionFlag,
		console.LogLevelFlag,
		console.QuietFlag,
		console.NoInteractionFlag,
		console.AnsiFlag,
		console.NoAnsiFlag,
		console.HelpFlag,
		console.VerbosityFlag("verbose", "v", "v"),
	}

	cmds := make([]*console.Command, 0, len(list.Commands))
	for _, legacyCmd := range list.Commands {
		if legacyCmd.Name.Namespace == "" && slices.Contains(overrideCommands, legacyCmd.Name.Command) {
			continue
		}

		cmd := &console.Command{
			Name:        legacyCmd.Name.Command,
			Usage:       legacyCmd.Description.String(),
			Description: legacyCmd.Help.String(),
			Category:    legacyCmd.Name.Namespace,
			Action:      legacyAction,
		}
		if legacyCmd.Hidden {
			cmd.Hidden = console.Hide
		}

		for _, alias := range legacyCmd.Aliases {
			cmd.Aliases = append(cmd.Aliases, &console.Alias{
				Name: alias,
			})
		}

		for pair := legacyCmd.Definition.Arguments.Oldest(); pair != nil; pair = pair.Next() {
			cmd.Args = append(cmd.Args, &console.Arg{
				Name:        pair.Value.Name,
				Description: pair.Value.Description.String(),
				Optional:    true,
				Slice:       bool(pair.Value.IsArray),
			})
		}

		for pair := legacyCmd.Definition.Options.Oldest(); pair != nil; pair = pair.Next() {
			isGlobal := false
			for _, fl := range globalFlags {
				if fl.Names()[0] == pair.Key {
					isGlobal = true
				}
			}
			if isGlobal {
				continue
			}

			if !pair.Value.AcceptValue {
				flag := &console.BoolFlag{
					Name:         strings.TrimPrefix(pair.Value.Name, "--"),
					Hidden:       pair.Value.Hidden,
					Usage:        pair.Value.Description.String(),
					Required:     false,
					DefaultValue: pair.Value.Default.BoolDefault(),
				}

				if pair.Value.Shortcut != "" {
					flag.Aliases = []string{strings.TrimPrefix(pair.Value.Shortcut, "-")}
				}
				cmd.Flags = append(cmd.Flags, flag)
				continue
			}

			if pair.Value.IsMultiple {
				flag := &console.StringSliceFlag{
					Name:     strings.TrimPrefix(pair.Value.Name, "--"),
					Hidden:   pair.Value.Hidden,
					Usage:    pair.Value.Description.String(),
					Required: false,
				}
				if pair.Value.Shortcut != "" {
					flag.Aliases = []string{strings.TrimPrefix(pair.Value.Shortcut, "-")}
				}
				cmd.Flags = append(cmd.Flags, flag)
				continue
			}

			flag := &console.StringFlag{
				Name:         strings.TrimPrefix(pair.Value.Name, "--"),
				Hidden:       pair.Value.Hidden,
				Usage:        pair.Value.Description.String(),
				Required:     false,
				DefaultValue: pair.Value.Default.StringDefault(),
			}
			if pair.Value.Shortcut != "" {
				flag.Aliases = []string{strings.TrimPrefix(pair.Value.Shortcut, "-")}
			}
			cmd.Flags = append(cmd.Flags, flag)
		}

		cmds = append(cmds, cmd)
	}

	return &console.Application{
		Name:        cnf.Application.Name,
		HelpName:    cnf.Application.Executable,
		Usage:       "",
		Version:     version,
		Channel:     "",
		Description: "",
		Commands: append(
			cmds,
			projectInitCommand(assets),
			validateCommand(assets),
			completionCommand(cnf),
		),
	}, nil
}

func debugLog(format string, v ...any) {
	if !viper.GetBool("debug") {
		return
	}

	log.Printf(format, v...)
}
