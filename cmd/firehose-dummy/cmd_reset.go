package main

import (
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	resetCommand = &cobra.Command{
		Use:   "reset",
		Short: "Reset local state",
		RunE:  runResetCommand,
	}
)

func runResetCommand(cmd *cobra.Command, args []string) error {
	// TODO: should be configurable?
	zlog.Info("removing data directory", zap.String("dir", "sf-data"))
	return os.RemoveAll("./sf-data")
}
