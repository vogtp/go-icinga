# go Icinga

## director

Generate director basket configs from a cobra.Command:

    func init() {
        testCmd.Flags().Bool("testFlagBool", false, "A boolean test flag")
        testCmd.Flags().String("testFlagString", "", "A string test flag")
        director.GenerateDirectorConfigPFlag(testCmd.Flags())
        testCmd.Flags().VisitAll(func(f *pflag.Flag) {
            if err := viper.BindPFlag(f.Name, f); err != nil {
                panic(err)
            }
        })
    }

    var (
        testCmd = &cobra.Command{
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
    )

This can be used as follows:

    testCmd  --icinga.director.import
