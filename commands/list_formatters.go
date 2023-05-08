package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/fatih/color"

	"github.com/platformsh/cli/internal/md"
)

type Formatter[T any] interface {
	Format(T) ([]byte, error)
}

type JSONListFormatter struct{}

func (f *JSONListFormatter) Format(list *List) ([]byte, error) {
	return json.Marshal(list)
}

type TXTListFormatter struct{}

func (f *TXTListFormatter) Format(list *List) ([]byte, error) {
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

	cmds := make(map[string][]Command)
	for _, cmd := range list.Commands {
		cmds[cmd.Name.Namespace] = append(cmds[cmd.Name.Namespace], cmd)
	}

	namespaces := make([]string, 0, len(cmds))
	for namespace := range cmds {
		namespaces = append(namespaces, namespace)
	}
	sort.Strings(namespaces)

	for _, namespace := range namespaces {
		if namespace != "" {
			fmt.Fprintln(writer, color.YellowString("%s\t", namespace))
		}
		for _, cmd := range cmds[namespace] {
			name := color.GreenString(cmd.Name.String())
			if len(cmd.Usage) > 1 {
				name = name + " (" + strings.Join(cmd.Usage[1:], ", ") + ")"
			}
			fmt.Fprintf(writer, "  %s\t%s\n", name, cmd.Description)
		}
	}
	writer.Flush()

	return b.Bytes(), nil
}

type RawListFormatter struct{}

func (f *RawListFormatter) Format(list *List) ([]byte, error) {
	var b bytes.Buffer
	writer := tabwriter.NewWriter(&b, 0, 8, 16, ' ', 0)
	for _, cmd := range list.Commands {
		fmt.Fprintf(writer, "%s\t%s\n", cmd.Name.String(), cmd.Description)
	}
	writer.Flush()

	return b.Bytes(), nil
}

type MDListFormatter struct{}

func (f *MDListFormatter) Format(list *List) ([]byte, error) {
	b := md.NewBuilder()
	b.H1(list.Application.Name + " " + list.Application.Version)

	cmds := make(map[string][]Command)
	for _, cmd := range list.Commands {
		cmds[cmd.Name.Namespace] = append(cmds[cmd.Name.Namespace], cmd)
	}

	namespaces := make([]string, 0, len(cmds))
	for namespace := range cmds {
		namespaces = append(namespaces, namespace)
	}
	sort.Strings(namespaces)

	for _, namespace := range namespaces {
		if namespace != "" {
			b.Paragraph(md.Bold(namespace)).Ln()
		}
		for _, cmd := range cmds[namespace] {
			b.ListItem(md.Link(md.Code(cmd.Name.String()), md.Anchor(cmd.Name.String())))
		}
		b.Ln()
	}

	for _, cmd := range list.Commands {
		b.H2(md.Code(cmd.Name.String()))
		b.Paragraph(cmd.Description).Ln()

		b.H3("Usage")
		for _, u := range cmd.Usage {
			b.ListItem(md.Code(u))
		}
		b.Ln()
		if cmd.Help != "" {
			b.Paragraph(cmd.Help).Ln()
		}

		if cmd.Definition.Arguments != nil && cmd.Definition.Arguments.Len() > 0 {
			b.H3("Arguments")
			for pair := cmd.Definition.Arguments.Oldest(); pair != nil; pair = pair.Next() {
				arg := pair.Value
				b.H4(md.Code(arg.Name))
				b.Paragraph(arg.Description).Ln()
				b.ListItem(fmt.Sprintf("Is required: %s", arg.IsRequired))
				b.ListItem(fmt.Sprintf("Is array: %s", arg.IsArray))
				b.ListItem(fmt.Sprintf("Default: %s", md.Code(arg.Default.String())))
				b.Ln()
			}
		}

		b.H3("Options")
		for pair := cmd.Definition.Options.Oldest(); pair != nil; pair = pair.Next() {
			opt := pair.Value
			name := opt.Name
			if opt.Shortcut != "" {
				name += "|" + opt.Shortcut
			}
			b.H4(md.Code(name))
			b.Paragraph(opt.Description).Ln()
			b.ListItem(fmt.Sprintf("Accept value: %s", opt.AcceptValue))
			b.ListItem(fmt.Sprintf("Is value required: %s", opt.IsValueRequired))
			b.ListItem(fmt.Sprintf("Is multiple: %s", opt.IsMultiple))
			b.ListItem(fmt.Sprintf("Default: %s", md.Code(opt.Default.String())))
			b.Ln()
		}
	}

	return []byte(b.String()), nil
}
