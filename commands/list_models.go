package commands

import (
	"sort"
	"strings"
)

var (
	ProjectInitCommand = Command{
		Name: "project:init",
		Usage: []string{
			"platform project:init",
			"ify",
		},
		Description: "Platformify the project", // TODO: replace with something meaningful
		Definition: Definition{
			Options: map[string]Option{},
		},
		Hidden: false,
	}
)

type List struct {
	Application Application `json:"application"`
	Commands    []Command   `json:"commands"`
	Namespace   string      `json:"namespace,omitempty"`
	Namespaces  []Namespace `json:"namespaces,omitempty"`
}

type Application struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Command struct {
	Name        string     `json:"name"`
	Usage       []string   `json:"usage"` // aliases
	Description string     `json:"description"`
	Definition  Definition `json:"definition"`
	Hidden      bool       `json:"hidden"`
}

type Definition struct {
	Options map[string]Option `json:"options" xml:"-"`
}

type Option struct {
	Name        string `json:"name"`
	Shortcut    string `json:"shortcut"`
	Description string `json:"description"`
}

type Namespace struct {
	ID       string   `json:"id"`
	Commands []string `json:"commands"` // the same as Command.Name
}

func (l *List) DescribesNamespace() bool {
	return l.Namespace != ""
}

func (l *List) AddCommand(namespace string, cmd Command) {
	for i := range l.Namespaces {
		name := &l.Namespaces[i]
		if name.ID == namespace {
			name.Commands = append(name.Commands, cmd.Name)
			sort.Strings(name.Commands)
		}
	}

	l.Commands = append(l.Commands, cmd)
	sort.Slice(l.Commands, func(i, j int) bool {
		const column = ":"
		switch {
		case !strings.Contains(l.Commands[i].Name, column) &&
			strings.Contains(l.Commands[j].Name, column):
			return true
		case strings.Contains(l.Commands[i].Name, column) &&
			!strings.Contains(l.Commands[j].Name, column):
			return false
		default:
			return l.Commands[i].Name < l.Commands[j].Name
		}
	})
}

func (l *List) RemoveHiddenCommands() {
	i := 0
	for _, cmd := range l.Commands {
		if !cmd.Hidden {
			l.Commands[i] = cmd
			i++
		}
	}
	l.Commands = l.Commands[:i]
}
