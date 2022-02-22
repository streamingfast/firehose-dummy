package main

import (
	"github.com/spf13/cobra"
	"github.com/streamingfast/dlauncher/launcher"
)

var (
	// GRPC Service Addresses
	BlockStreamServingAddr  = ":9000"
	RelayerServingAddr      = ":9010"
	MergerServingAddr       = ":9020"
	FirehoseGRPCServingAddr = ":9030"

	// Blocks store
	MergedBlocksStoreURL string = "file://{sf-data-dir}/storage/merged-blocks"
	OneBlockStoreURL     string = "file://{sf-data-dir}/storage/one-blocks"

	// Protocol defaults
	FirstStreamableBlock uint64 = 1
)

func init() {
	launcher.RegisterCommonFlags = func(cmd *cobra.Command) error {
		flags := cmd.Flags()

		// Logging
		flags.Int("verbose", 1, "Logging verbosity")
		flags.String("log-format", "text", "Logging format")
		flags.String("config", "firehose.yml", "Path to configuration file")

		// Common stores configuration flags
		flags.String("common-blocks-store-url", MergedBlocksStoreURL, "Store URL (with prefix) where to read/write")
		flags.String("common-oneblock-store-url", OneBlockStoreURL, "Store URL (with prefix) to read/write one-block files")
		flags.String("common-blockstream-addr", RelayerServingAddr, "GRPC endpoint to get real-time blocks")
		flags.Uint64("common-first-streamable-block", FirstStreamableBlock, "First streamable block number")

		// Authentication, metering and rate limiter plugins
		flags.String("common-auth-plugin", "null://", "Auth plugin URI, see streamingfast/dauth repository")
		flags.String("common-metering-plugin", "null://", "Metering plugin URI, see streamingfast/dmetering repository")

		// System Behavior
		flags.Duration("common-shutdown-delay", 5, "Add a delay between receiving SIGTERM signal and shutting down apps. Apps will respond negatively to /healthz during this period")

		return nil
	}
}
