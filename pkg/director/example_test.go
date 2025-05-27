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
)

func getTestCmd() *cobra.Command {
	testCmd := &cobra.Command{
		Use: "testCmd",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if director.ShouldGenerate() {
				d := director.Generator{
					CobraCmd: cmd,
				}
				d.Generate(os.Stdout)
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
		return testUUIDs[testUUIDsIdx]
	}
	return func() {
		testUUIDsIdx = -1
		director.GenerateUUID = orig
	}
}

func Test_Generate(t *testing.T) {
	reset := initUUIDGenerator()
	defer reset()
	g := director.Generator{
		CobraCmd: getTestCmd(),
	}
	var out bytes.Buffer
	g.Generate(&out)
	should, err := os.ReadFile("testfiles/testCmdSimple.json")
	if err != nil {
		t.Fatalf("Cannot read test output: %v", err)
	}
	outStr := out.String()
	if string(should) != outStr {
		fmt.Println(out.String())
		of := "testfiles/ignore_test_output_testCmdSimple.json"
		f, err := os.OpenFile(of, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			t.Errorf("Cannot open output file: %v", err)
		}
		defer f.Close()
		if _, err := f.WriteString(outStr); err != nil {
			t.Errorf("Cannot write output file: %v", err)
		}
		t.Error("Output not identital")
	}

}
