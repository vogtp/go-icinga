package director

import (
	"fmt"
	"slices"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	WriteConfigFlagName  = "icinga.director.write"  // GenerateFlagName is the name of the flag used to trigger gerneration
	ImportConfigFlagName = "icinga.director.import" // ImportFlagName run icinga cli director basket restore
	CommandDIrFlagName   = "icinga.command.dir"     // Directory where the check command are stored
)

var (
	ignoredFlags []string = []string{WriteConfigFlagName, ImportConfigFlagName, "help"}
)

// IgnoreFlag adds a flag to the ignored flags
func IgnoreFlag(f string) {
	ignoredFlags = append(ignoredFlags, f)
}

// Flags adds a pflag for gerneration to the FlagSet
func Flags(s *pflag.FlagSet) {
	s.Bool(WriteConfigFlagName, false, "Generate a icinga director config and write it out")
	s.Bool(ImportConfigFlagName, false, "Generate a icinga director config and run icinga cli director basket restore")
	s.String(CommandDIrFlagName, "/usr/lib/nagios/plugins/", "Directory where the check command are stored")
}

// ShouldGenerate checks the flag if generation should be triggered
func ShouldGenerate() bool {
	return viper.GetBool(WriteConfigFlagName) || viper.GetBool(ImportConfigFlagName)
}

func (g *Generator) parsePFlags() {
	args := make(map[string]argument)
	datafields := make(map[int]datafield)
	fieldID := 1
	cmdFields := make([]cmdField, 0)
	g.CobraCmd.Flags().VisitAll(func(f *pflag.Flag) {
		if f.Hidden || slices.Contains(ignoredFlags, f.Name) {
			return
		}
		
		fName := idPrintf("%s_%s", g.id, f.Name)

		datatype := "Icinga\\Module\\Director\\DataType\\DataTypeString"
		if f.Value.Type() == "bool" {
			datatype = "Icinga\\Module\\Director\\DataType\\DataTypeBoolean"
			args[fmt.Sprintf("--%s", f.Name)] = argument{
				SetIf: fmt.Sprintf("$%s$", fName),
			}
		} else {
			args[fmt.Sprintf("--%s", f.Name)] = argument{
				Value: fmt.Sprintf("$%s$", fName),
			}
		}
		datafields[fieldID] = datafield{
			Varname:     fName,
			Caption:     fmt.Sprintf("%s: %s", g.name, strings.ReplaceAll(f.Name, "-", " ")),
			Description: f.Usage,
			Datatype:    datatype,
			UUID:        GenerateUUID(),
		}
		defVal := f.Value.String()
		if len(defVal) < 1{
			defVal = viper.GetString(f.Name)
		}
		if strings.HasSuffix(f.Value.Type(), "Slice") {
			defVal = strings.ReplaceAll(defVal, "[", "")
			defVal = strings.ReplaceAll(defVal, "]", "")
			defVal = strings.ReplaceAll(defVal, ", ", ",")
		}
		g.srvDef.Vars[fName] = defVal
		cmdFields = append(cmdFields, cmdField{
			DatafieldID: fieldID,
			IsRequired:  "n",
		})
		fieldID++
	})
	order := 1
	for i, cp := range g.cobraParams {
		fName := idPrintf("%s_cmd_%v", g.id, i)
		args[cp] = argument{
			//	SetIf:    fmt.Sprintf("$%s$", fName),
			Value:    fmt.Sprintf("$%s$", fName),
			Required: true,
			SkipKey:  true,
			Order:    order,
		}
		order++

		datafields[fieldID] = datafield{
			Varname:  fName,
			Caption:  fmt.Sprintf("%s: Command%v", g.name, i),
			Datatype: "Icinga\\Module\\Director\\DataType\\DataTypeString",
			UUID:     GenerateUUID(),
		}
		g.srvDef.Vars[fName] = cp
		cmdFields = append(cmdFields, cmdField{
			DatafieldID: fieldID,

			IsRequired: "y",
		})
		fieldID++
	}
	g.c.Datafield = datafields

	g.cmdDef.Arguments = args
	g.cmdDef.Fields = cmdFields
}
