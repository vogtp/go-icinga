package director

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	GenerateFlagName = "icinga.director"
)

func GenerateDirectorConfigPFlag(s *pflag.FlagSet) {
	s.Bool(GenerateFlagName, false, "Generate a icinga director config")
}

func ShouldGenerate() bool {
	b := viper.GetBool(GenerateFlagName)
	return b
}

func (g *Generator) parsePFlags() {
	args := make(map[string]argument)
	datafields := make(map[int]datafield)
	fieldID := 1
	cmdFields := make([]cmdField, 0)
	g.CobraCmd.Flags().VisitAll(func(f *pflag.Flag) {
		if f.Hidden || f.Name == GenerateFlagName || f.Name == "help" {
			return
		}
		fName := idPrintf("%s_%s", g.id, f.Name)
		args[fmt.Sprintf("--%s", f.Name)] = argument{
			SetIf: fmt.Sprintf("$%s$", fName),
		}
		datatype := "Icinga\\Module\\Director\\DataType\\DataTypeString"
		if f.Value.Type() == "bool" {
			datatype = "Icinga\\Module\\Director\\DataType\\DataTypeBoolean"
		}
		datafields[fieldID] = datafield{
			Varname:     fName,
			Caption:     fmt.Sprintf("%s: %s", g.name, strings.ReplaceAll(f.Name, "-", " ")),
			Description: f.Usage,
			Datatype:    datatype,
			UUID:        uuid.NewString(),
		}
		g.srvDef.Vars[fName] = f.DefValue
		cmdFields = append(cmdFields, cmdField{
			DatafieldID: fieldID,
			IsRequired:  "n",
		})
		fieldID++
	})
	for i, cp := range g.cobraParams {
		fName := idPrintf("%s_cmd_%v", g.id, i)
		args[fmt.Sprintf("%s", cp)] = argument{
			SetIf:    fmt.Sprintf("$%s$", fName),
			Required: true,
			SkeyKey:  true,
		}

		datafields[fieldID] = datafield{
			Varname:  fName,
			Caption:  fmt.Sprintf("%s: Command%v", g.name, i),
			Datatype: "Icinga\\Module\\Director\\DataType\\DataTypeString",
			UUID:     uuid.NewString(),
		}
		g.srvDef.Vars[fName] = cp
		cmdFields = append(cmdFields, cmdField{
			DatafieldID: fieldID,
			IsRequired:  "y",
		})
		fieldID++
	}
	g.c.Datafield = datafields

	g.cmdDef.Arguments = args
	g.cmdDef.Fields = cmdFields
}
