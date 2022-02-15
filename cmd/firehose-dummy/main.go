package main

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/streamingfast/derr"
	"github.com/streamingfast/dlauncher/flags"
	"github.com/streamingfast/dlauncher/launcher"

	"github.com/streamingfast/firehose-dummy/codec"
)

var (
	rootCmd = &cobra.Command{
		Use:   "firehose-dummy",
		Short: "Dummy Chain Firehose",
	}
)

func main() {
	cobra.OnInitialize(func() {
		//
		// This will allow setting flags via environment variables
		// Example: SF_GLOBAL_VERBOSE=3 will set `global-verbose` flag
		//
		flags.AutoBind(rootCmd, "SF")
	})

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		launcher.SetupLogger(&launcher.LoggingOptions{
			WorkingDir:    viper.GetString("data-dir"),
			Verbosity:     viper.GetInt("verbose"),
			LogFormat:     viper.GetString("log-format"),
			LogToFile:     viper.GetBool("log-to-file"),
			LogListenAddr: viper.GetString("log-level-switcher-listen-addr"),
		})
		return nil
	}

	rootCmd.AddCommand(
		initCommand,
		startCommand,
		setupCommand,
	)

	codec.SetProtocolFirstStreamableBlock(FirstStreamableBlock)

	derr.Check("validating codec settings", codec.Validate())
	derr.Check("registering application flags", launcher.RegisterFlags(startCommand))
	derr.Check("executing root command", rootCmd.Execute())
}
