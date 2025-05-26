package director

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

type Generator struct {
	NamePrefix string
	name       string
	id         string
	//Flags          *pflag.FlagSet
	CobraCmd       *cobra.Command
	Description    string
	DescriptionURL string

	cobraParams []string
	srvDef      Service
	cmdDef      CommandDefinition
	c           *Config
}

func (g *Generator) Generate(w io.Writer) {
	g.cobraParams = getCobraParams(g.CobraCmd)
	g.name = strings.Join(g.cobraParams, " ")
	g.id = strings.Join(g.cobraParams, "-")
	if len(g.NamePrefix) > 0 && !strings.HasSuffix(g.NamePrefix, "-") {
		g.NamePrefix = fmt.Sprintf("%v-", g.NamePrefix)
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	c := g.generate()
	if err := enc.Encode(&c); err != nil {
		slog.Error("Cannot encode", "err", err)
	}
}

func (g *Generator) generate() *Config {
	g.c = &Config{
		Command:         make(Command),
		ServiceTemplate: make(ServiceTemplate),
	}
	cmdID := fmt.Sprintf("%vcmd-check-%s", g.NamePrefix, g.id)
	g.cmdDef = CommandDefinition{
		Command:        fmt.Sprintf("/usr/lib64/nagios/plugins/%s", g.cobraParams[0]),
		Imports:        make([]interface{}, 0),
		MethodsExecute: "PluginCheck",
		ObjectName:     cmdID,
		ObjectType:     "object",
		Timeout:        30,
		UUID:           uuid.NewString(),
	}
	g.cobraParams = g.cobraParams[1:]
	srvID := fmt.Sprintf("%vtpl-service-%s", g.NamePrefix, g.id)
	g.srvDef = Service{
		ObjectName:    srvID,
		CheckCommand:  cmdID,
		CheckInterval: 86400,
		RetryInterval: 3600,
		Notes:         g.Description,
		NotesURL:      g.DescriptionURL,
		//IconImage: ,
		Imports:    []string{"tpl-service-generic"},
		ObjectType: "template",
		Fields:     make([]interface{}, 0),
		Vars:       make(map[string]any),
		UUID:       uuid.NewString(),
	}
	g.srvDef.Vars["criticality"] = "C"

	g.parsePFlags()

	g.c.Command[cmdID] = g.cmdDef
	g.c.ServiceTemplate[srvID] = g.srvDef
	return g.c
}

func idPrintf(format string, a ...any) string {
	s := fmt.Sprintf(format, a...)
	s = strings.ReplaceAll(s, "-", "_")
	s = strings.ReplaceAll(s, ".", "_")
	return s
}

func getCobraParams(cmd *cobra.Command) []string {
	n := cmd.Name()
	p := cmd.Parent()
	if p == nil {
		return []string{n}
	}
	cmds := getCobraParams(cmd.Parent())
	return append(cmds, n)
}
