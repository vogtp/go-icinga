package director

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vogtp/go-icinga/pkg/icinga"
	"github.com/vogtp/go-icinga/pkg/icingacli"
)

// Generator for icinga directory basket config files
// CobraCmd is mandatory and all the config will be genertated from it
type Generator struct {
	NamePrefix     string         // NamePrefix is a optional prefix for the icinga names
	CobraCmd       *cobra.Command // CobraCmd is mandatory and all the config will be genertated from it
	Description    string         // Description (optional) is the icinga Notes field
	DescriptionURL string         // DescriptionURL (optional) is the icinga NotesURL field
	Output         io.Writer      // Output is a io.Writer where the string output is written to
	Criticality    icinga.Criticality

	name        string
	id          string
	cobraParams []string
	srvDef      service
	cmdDef      CommandDefinition
	c           *Config
}

// Generate writes the icinga basket config to the passed writer
func (g *Generator) Generate() error {
	g.cobraParams = getCobraParams(g.CobraCmd)
	g.name = strings.Join(g.cobraParams, " ")
	g.id = strings.Join(g.cobraParams, "-")
	if len(g.NamePrefix) > 0 && !strings.HasSuffix(g.NamePrefix, "-") {
		g.NamePrefix = fmt.Sprintf("%v-", g.NamePrefix)
	}
	c := g.generate()

	var out bytes.Buffer
	var w io.Writer = &out
	if g.Output != nil && viper.GetBool(WriteConfigFlagName) {
		w = io.MultiWriter(w, g.Output)
	}
	if err := prettyPrint(c, w); err != nil {
		return fmt.Errorf("cannot write config: %w", err)
	}
	if viper.GetBool(ImportConfigFlagName) {
		slog.Info("Importing icinga director config", "plugin", g.name)
		if err := icingacli.ImportDirectorBasket(&out); err != nil {
			return fmt.Errorf("cannot import basket into director: %w", err)
		}
	}
	return nil
}

func prettyPrint(c *Config, w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(&c); err != nil {
		return fmt.Errorf("cannot encode config: %w", err)
	}
	return nil
}

func (g *Generator) generate() *Config {
	g.c = &Config{
		Command:         make(command),
		ServiceTemplate: make(serviceTemplate),
	}
	cmdID := fmt.Sprintf("%vcmd-check-%s", g.NamePrefix, g.id)
	g.cmdDef = CommandDefinition{
		Command:        fmt.Sprintf("%s%s", viper.GetString(CommandDIrFlagName), g.cobraParams[0]),
		Imports:        make([]interface{}, 0),
		MethodsExecute: "PluginCheck",
		ObjectName:     cmdID,
		ObjectType:     "object",
		Timeout:        30,
		UUID:           GenerateUUID(),
	}
	g.cobraParams = g.cobraParams[1:]
	srvID := fmt.Sprintf("%vtpl-service-%s", g.NamePrefix, g.id)
	g.srvDef = service{
		ObjectName:    srvID,
		CheckCommand:  cmdID,
		CheckInterval: 300,
		RetryInterval: 60,
		Notes:         g.Description,
		NotesURL:      g.DescriptionURL,
		//CommandEndpoint:     CommandEndpoint, //FIXME make generic
		UseAgent:            false,
		EnablePerfdata:      true,
		EnableNotifications: true,
		MaxCheckAttempts:    3,
		//IconImage: ,
		//Imports:    []string{"tpl-service-generic"},
		ObjectType: "template",
		Fields:     make([]interface{}, 0),
		Vars:       make(map[string]any),
		UUID:       GenerateUUID(),
	}
	g.srvDef.Vars["criticality"] = g.Criticality.Get()

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
