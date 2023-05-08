package commands

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"

	orderedmap "github.com/wk8/go-ordered-map/v2"
)

var (
	ProjectInitCommand = Command{
		Name: CommandName{
			Namespace: "project",
			Command:   "init",
		},
		Usage: []string{
			"platform project:init",
			"ify",
		},
		Description: "Initialize the project",
		Help:        "",
		Definition: Definition{
			Arguments: nil,
			Options: orderedmap.New[string, Option](orderedmap.WithInitialData[string, Option](
				orderedmap.Pair[string, Option]{
					Key:   HelpOption.GetName(),
					Value: HelpOption,
				},
				orderedmap.Pair[string, Option]{
					Key:   VerboseOption.GetName(),
					Value: VerboseOption,
				},
				orderedmap.Pair[string, Option]{
					Key:   VersionOption.GetName(),
					Value: VersionOption,
				},
				orderedmap.Pair[string, Option]{
					Key:   YesOption.GetName(),
					Value: YesOption,
				},
				orderedmap.Pair[string, Option]{
					Key:   NoInteractionOption.GetName(),
					Value: NoInteractionOption,
				},
				orderedmap.Pair[string, Option]{
					Key:   AnsiOption.GetName(),
					Value: AnsiOption,
				},
				orderedmap.Pair[string, Option]{
					Key:   NoAnsiOption.GetName(),
					Value: NoAnsiOption,
				},
				orderedmap.Pair[string, Option]{
					Key:   NoOption.GetName(),
					Value: NoOption,
				},
				orderedmap.Pair[string, Option]{
					Key:   QuietOption.GetName(),
					Value: QuietOption,
				},
			)),
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
	Name        CommandName `json:"name"`
	Usage       []string    `json:"usage"` // + aliases
	Description string      `json:"description"`
	Help        string      `json:"help"`
	Definition  Definition  `json:"definition"`
	Hidden      bool        `json:"hidden"`
}

type CommandName struct {
	Namespace string
	Command   string
}

func (n *CommandName) String() string {
	if n.Namespace == "" {
		return n.Command
	}
	return n.Namespace + ":" + n.Command
}

func (n *CommandName) ContainsNamespace() bool {
	return n.Namespace != ""
}

func (n *CommandName) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.String())
}

func (n *CommandName) UnmarshalJSON(text []byte) error {
	var command string
	err := json.Unmarshal(text, &command)
	if err != nil {
		return err
	}
	names := strings.SplitN(command, ":", 2)
	switch {
	case len(names) == 1:
		n.Command = names[0]
	case len(names) > 1:
		n.Namespace = names[0]
		n.Command = names[1]
	}
	return nil
}

type Definition struct {
	Arguments *Arguments                             `json:"arguments"`
	Options   *orderedmap.OrderedMap[string, Option] `json:"options"`
}

// Arguments is a custom type of "orderedmap.OrderedMap[string, Argument]" that we need
// for the one reason only - to support the Legacy-CLI behavior. The Legacy-CLI can return
// both an empty array and an object/map in the "arguments" field. This type can be
// deleted and replaced with map[string]Argument in Definition structure when
// the json-schema will be fixed in the Legacy-CLI.
type Arguments struct {
	orderedmap.OrderedMap[string, Argument]
}

func (a *Arguments) MarshalJSON() ([]byte, error) {
	if a.Len() == 0 {
		return []byte("[]"), nil
	}
	return json.Marshal(&a.OrderedMap)
}

func (a *Arguments) UnmarshalJSON(text []byte) error {
	if string(text) == "[]" {
		return nil
	}

	var arguments orderedmap.OrderedMap[string, Argument]
	err := json.Unmarshal(text, &arguments)
	if err != nil {
		return err
	}
	*a = Arguments{OrderedMap: arguments}

	return nil
}

type Argument struct {
	Name        string `json:"name"`
	IsRequired  YesNo  `json:"is_required"`
	IsArray     YesNo  `json:"is_array"`
	Description string `json:"description"`
	Default     Any    `json:"default"`
}

type Option struct {
	Name            string `json:"name"`
	Shortcut        string `json:"shortcut"`
	AcceptValue     YesNo  `json:"accept_value"`
	IsValueRequired YesNo  `json:"is_value_required"`
	IsMultiple      YesNo  `json:"is_multiple"`
	Description     string `json:"description"`
	Default         Any    `json:"default"`
}

func (o *Option) GetName() string {
	return strings.TrimPrefix(o.Name, "--")
}

type YesNo bool

func (y YesNo) String() string {
	if y {
		return "yes"
	}
	return "no"
}

type Any struct {
	any
}

func (a *Any) String() string {
	if a.any == nil {
		return "NULL"
	}
	switch t := a.any.(type) {
	case bool:
		return fmt.Sprintf("%t", a.any)
	case float32, float64:
		s := a.any.(float64) //nolint:errcheck
		if s == math.Trunc(s) {
			return fmt.Sprintf("%d", int64(s))
		}
		return fmt.Sprintf("%f", s)
	case string:
		s := a.any.(string) //nolint:errcheck
		return fmt.Sprintf("'%s'", s)
	case []any, []string, []int, []float64:
		return "array ()"
	default:
		panic(fmt.Sprintf("options: unsupported type: %T", t))
	}
}

func (a Any) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.any)
}

func (a *Any) UnmarshalJSON(text []byte) error {
	return json.Unmarshal(text, &a.any)
}

type Namespace struct {
	ID       string   `json:"id"`
	Commands []string `json:"commands"` // the same as Command.Name
}

func (l *List) DescribesNamespace() bool {
	return l.Namespace != ""
}

func (l *List) AddCommand(cmd *Command) {
	for i := range l.Namespaces {
		name := &l.Namespaces[i]
		if name.ID == cmd.Name.Namespace {
			name.Commands = append(name.Commands, cmd.Name.String())
			sort.Strings(name.Commands)
		}
	}

	l.Commands = append(l.Commands, *cmd)
	sort.Slice(l.Commands, func(i, j int) bool {
		switch {
		case !l.Commands[i].Name.ContainsNamespace() && l.Commands[j].Name.ContainsNamespace():
			return true
		case l.Commands[i].Name.ContainsNamespace() && !l.Commands[j].Name.ContainsNamespace():
			return false
		default:
			return l.Commands[i].Name.String() < l.Commands[j].Name.String()
		}
	})
}
