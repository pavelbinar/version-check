package cmd

import (
	"os"

	internal "github.com/pavelbinar/version-check/internal"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use:   "version-checker",
		Short: "A tool to check versions of installed software",
		Long: `Version Checker is a CLI tool that checks the versions of installed software
against expected versions specified in a configuration file.`,
		Run: func(cmd *cobra.Command, args []string) {
			internal.RunVersionCheck(cmd, args, cfgFile)
		},
	}
)

func initConfig() {
	if cfgFile == "" {
		cfgFile = "config.yaml"
	}
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is ./config.yaml)")
	rootCmd.Version = "0.1.0"

}
