package commands

import (
	"github.com/fatih/color"
)

type InputOption struct {
	Name        string
	Shortcut    string
	Description string
}

var (
	HelpOption = InputOption{
		Name:        "--help",
		Shortcut:    "-h",
		Description: "Display this help message",
	}
	VerboseOption = InputOption{
		Name:        "--verbose",
		Shortcut:    "-v|vv|vvv",
		Description: "Increase the verbosity of messages",
	}
	VersionOption = InputOption{
		Name:        "--version",
		Shortcut:    "-V",
		Description: "Display this application version",
	}
	YesOption = InputOption{
		Name:        "--yes",
		Shortcut:    "-y",
		Description: "Answer \"yes\" to confirmation questions; accept the default value for other questions; disable interaction",
	}
	NoInteractionOption = InputOption{
		Name:        "--no-interaction",
		Description: "Do not ask any interactive questions; accept default values. Equivalent to using the environment variable: " + color.YellowString("PLATFORMSH_CLI_NO_INTERACTION=1"),
	}
	AnsiOption = InputOption{
		Name:        "--ansi",
		Description: "Force ANSI output",
	}
	NoAnsiOption = InputOption{
		Name:        "--no-ansi",
		Description: "Disable ANSI output",
	}
	NoOption = InputOption{
		Name:        "--no",
		Shortcut:    "-n",
		Description: "Answer \"no\" to confirmation questions; accept the default value for other questions; disable interaction",
	}
	QuietOption = InputOption{
		Name:        "--quiet",
		Shortcut:    "-q",
		Description: "Do not output any message",
	}
)

var (
	GlobalOptions = []InputOption{
		HelpOption,
		VerboseOption,
		VersionOption,
		YesOption,
		NoInteractionOption,
		AnsiOption,
		NoAnsiOption,
		NoOption,
		QuietOption,
	}
)
