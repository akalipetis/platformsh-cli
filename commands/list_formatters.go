package commands

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/fatih/color"
)

type Formatter interface {
	Format(list *List) ([]byte, error)
}

type JSONFormatter struct{}

func (f *JSONFormatter) Format(list *List) ([]byte, error) {
	return json.Marshal(list)
}

type XMLFormatter struct{}

func (f *XMLFormatter) Format(list *List) ([]byte, error) {
	return xml.Marshal(list)
}

type TXTFormatter struct{}

func (f *TXTFormatter) Format(list *List) ([]byte, error) {
	var b bytes.Buffer
	writer := tabwriter.NewWriter(&b, 0, 8, 1, ' ', 0)
	fmt.Fprintf(writer, "%s %s\n", list.Application.Name, color.GreenString(list.Application.Version))
	fmt.Fprintln(writer)
	fmt.Fprintln(writer, color.YellowString("Global options:"))
	for _, opt := range GlobalOptions {
		shortcut := opt.Shortcut
		if shortcut == "" {
			shortcut = "  "
		}
		fmt.Fprintf(writer, "  %s\t%s %s\n", color.GreenString(opt.Name), color.GreenString(shortcut), opt.Description)
	}
	fmt.Fprintln(writer)

	writer.Init(&b, 0, 8, 4, ' ', 0)
	if list.DescribesNamespace() {
		fmt.Fprintln(writer, color.YellowString("Available commands for the \"%s\" namespace:", list.Namespace))
	} else {
		fmt.Fprintln(writer, color.YellowString("Available commands:"))
	}

	var namespace string
	for _, cmd := range list.Commands {
		if !list.DescribesNamespace() {
			names := strings.SplitN(cmd.Name, ":", 2)
			if len(names) == 0 {
				continue
			}

			if len(names) > 1 && names[0] != namespace {
				fmt.Fprintln(writer, color.YellowString("%s\t", names[0]))
				namespace = names[0]
			}
		}

		name := color.GreenString(cmd.Name)
		if len(cmd.Usage) > 1 {
			name = name + " (" + strings.Join(cmd.Usage[1:], ", ") + ")"
		}
		fmt.Fprintf(writer, "  %s\t%s\n", name, cmd.Description)
	}
	writer.Flush()

	return b.Bytes(), nil
}

type RawFormatter struct{}

func (f *RawFormatter) Format(list *List) ([]byte, error) {
	var b bytes.Buffer
	writer := tabwriter.NewWriter(&b, 0, 8, 16, ' ', 0)
	for _, cmd := range list.Commands {
		fmt.Fprintf(writer, "%s\t%s\n", cmd.Name, cmd.Description)
	}
	writer.Flush()

	return b.Bytes(), nil
}

type MDFormatter struct{}

func (f *MDFormatter) Format(_ *List) ([]byte, error) {
	// TODO: implement the method
	return []byte("The format has not been implemented yet."), nil
}
