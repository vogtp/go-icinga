package director_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/vogtp/go-icinga/pkg/director"
	"github.com/vogtp/go-icinga/pkg/icinga"
)

func cmpGenerate(t *testing.T, testName string, cmd *cobra.Command, crit icinga.Criticality) {
	viper.Set(director.WriteConfigFlagName, true)
	reset := initUUIDGenerator()
	defer reset()
	g := director.Generator{
		CobraCmd:       cmd,
		Description:    "Test Icinga Directory Bucket",
		DescriptionURL: "https://github.com/vogtp/go-icinga/",
		Criticality:    crit,
	}
	var out bytes.Buffer
	g.Output = &out
	if err := g.Generate(); err != nil {
		t.Errorf("Error generating config: %v", err)
	}
	should, err := os.ReadFile(testFileName(testName))
	if err != nil {
		t.Errorf("Cannot read test output: %v", err)
	}
	outStr := out.String()
	if string(should) != outStr {
		fmt.Println(out.String())
		f, err := os.OpenFile(shouldFileName(testName), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			t.Errorf("Cannot open output file: %v", err)
		}
		defer f.Close()
		if _, err := f.WriteString(outStr); err != nil {
			t.Errorf("Cannot write output file: %v", err)
		}
		t.Errorf("Output not identital.\nComape %s and %s", shouldFileName(testName), testFileName(testName))
	}
}

func getTestCmd() *cobra.Command {
	testCmd := &cobra.Command{
		Use: "testCmd",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if director.ShouldGenerate() {
				d := director.Generator{
					CobraCmd: cmd,
					Output:   os.Stdout,
				}
				d.Generate()
				os.Exit(0)
			}
		},
	}
	testCmd.Flags().Bool("testFlagBool", false, "A boolean test flag")
	testCmd.Flags().String("testFlagString", "", "A string test flag")
	director.GenerateDirectorConfigPFlag(testCmd.Flags())
	testCmd.Flags().VisitAll(func(f *pflag.Flag) {
		if err := viper.BindPFlag(f.Name, f); err != nil {
			panic(err)
		}
	})
	return testCmd
}

var (
	testUUIDs = []string{
		"4c9c22a6-80a5-4cfb-a857-896d84f24e04",
		"3300dc08-2db2-414a-badf-8f1f5bb266bb",
		"3e4a7d42-7e6c-4e07-8c8d-746bef32adf3",
		"f4880941-391f-4c98-9361-a70256fa6d5d",
	}
	testUUIDsIdx = -1
)

func initUUIDGenerator() func() {
	orig := director.GenerateUUID
	director.GenerateUUID = func() string {
		testUUIDsIdx++
		if testUUIDsIdx >= len(testUUIDs) {
			testUUIDsIdx = 0
		}
		return testUUIDs[testUUIDsIdx]
	}
	return func() {
		testUUIDsIdx = -1
		director.GenerateUUID = orig
	}
}

func testFileName(testName string) string {
	return fmt.Sprintf("testfiles/%s.json", testName)
}

func shouldFileName(testName string) string {
	return fmt.Sprintf("testfiles/ignore_test_output_%s.json", testName)
}

func Test_GenerateCommand(t *testing.T) {
	testName := "testCmdSimple"
	cmd := getTestCmd()
	var crit icinga.Criticality
	cmpGenerate(t, testName, cmd, crit)
}

func Test_GenerateSubCommand(t *testing.T) {
	testName := "testCmdSub"
	cmd := getTestCmd()
	testCmd := &cobra.Command{
		Use: "testSubCmd",
	}
	testCmd.Flags().Bool("testFlagBoolSubCmd", false, "A boolean test sub command flag")
	testCmd.Flags().String("testFlagStringSubCmd", "", "A string test sub command flag")
	testCmd.Flags().VisitAll(func(f *pflag.Flag) {
		if err := viper.BindPFlag(f.Name, f); err != nil {
			panic(err)
		}
	})
	cmd.AddCommand(testCmd)
	cmpGenerate(t, testName, testCmd, icinga.Criticality7x24)
}
