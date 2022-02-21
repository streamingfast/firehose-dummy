package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/streamingfast/derr"
	"github.com/streamingfast/dlauncher/flags"
	"github.com/streamingfast/dlauncher/launcher"
	"github.com/streamingfast/logging"
	"go.uber.org/zap"

	"github.com/streamingfast/firehose-dummy/codec"
)

const (
	version = "0.0.1"
)

var (
	userLog  = launcher.UserLog
	zlog     *zap.Logger
	allFlags = map[string]bool{}

	rootCmd = &cobra.Command{
		Use:     "firehose-dummy",
		Short:   "Dummy Chain Firehose",
		Version: version,
	}
)

func init() {
	logging.Register("main", &zlog)
	logging.Set(logging.MustCreateLogger())
}

func main() {
	cobra.OnInitialize(func() {
		// This will allow setting flags via environment variables
		// Example: FH_GLOBAL_VERBOSE=3 will set `global-verbose` flag
		allFlags = flags.AutoBind(rootCmd, "FH")
	})

	rootCmd.PersistentPreRunE = preRun

	rootCmd.AddCommand(
		initCommand,
		startCommand,
		setupCommand,
		resetCommand,
	)

	codec.SetProtocolFirstStreamableBlock(FirstStreamableBlock)

	derr.Check("registering application flags", launcher.RegisterFlags(startCommand))
	derr.Check("executing root command", rootCmd.Execute())
}

func preRun(cmd *cobra.Command, args []string) error {
	if configFile := viper.GetString("config"); configFile != "" {
		if fileExists(configFile) {
			if err := launcher.LoadConfigFile(configFile); err != nil {
				return err
			}
		}
	}

	cfg := launcher.DfuseConfig[cmd.Name()]
	if cfg != nil {
		for k, v := range cfg.Flags {
			validFlag := false
			if _, ok := allFlags[k]; ok {
				viper.SetDefault(k, v)
				validFlag = true
			}
			if !validFlag {
				return fmt.Errorf("invalid flag %v", k)
			}
		}
	}

	launcher.SetupLogger(&launcher.LoggingOptions{
		WorkingDir:    viper.GetString("data-dir"),
		Verbosity:     viper.GetInt("verbose"),
		LogFormat:     viper.GetString("log-format"),
		LogToFile:     viper.GetBool("log-to-file"),
		LogListenAddr: viper.GetString("log-level-switcher-listen-addr"),
	})

	codec.SetProtocolFirstStreamableBlock(viper.GetUint64("common-first-streamable-block"))

	return codec.Validate()
}
