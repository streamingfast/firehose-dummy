package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/streamingfast/bstream"
	dauthAuthenticator "github.com/streamingfast/dauth/authenticator"
	_ "github.com/streamingfast/dauth/authenticator/null"
	"github.com/streamingfast/dlauncher/launcher"
	"github.com/streamingfast/dmetering"
	"github.com/streamingfast/dmetrics"
	firehoseApp "github.com/streamingfast/firehose/app/firehose"
	"github.com/streamingfast/logging"

	"github.com/streamingfast/firehose-dummy/codec"
)

var (
	metricset           = dmetrics.NewSet()
	headBlockNumMetric  = metricset.NewHeadBlockNumber("firehose")
	headTimeDriftmetric = metricset.NewHeadTimeDrift("firehose")
)

func init() {
	appLogger := zap.NewNop()
	logging.Register("firehose", &appLogger)

	registerFlags := func(cmd *cobra.Command) error {
		cmd.Flags().String("firehose-grpc-listen-addr", FirehoseGRPCServingAddr, "Address on which the firehose will listen")
		cmd.Flags().StringSlice("firehose-blocks-store-urls", nil, "If non-empty, overrides common-blocks-store-url with a list of blocks stores")
		cmd.Flags().Duration("firehose-real-time-tolerance", 1*time.Minute, "firehose will became alive if now - block time is smaller then tolerance")
		return nil
	}

	factoryFunc := func(runtime *launcher.Runtime) (launcher.App, error) {
		sfDataDir := runtime.AbsDataDir
		tracker := runtime.Tracker.Clone()

		// Configure block stream connection (to node or relayer)
		blockstreamAddr := viper.GetString("common-blockstream-addr")
		if blockstreamAddr != "" {
			tracker.AddGetter(bstream.BlockStreamLIBTarget, bstream.StreamLIBBlockRefGetter(blockstreamAddr))
		}

		// Configure authentication (default is no auth)
		authenticator, err := dauthAuthenticator.New(viper.GetString("common-auth-plugin"))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize dauth: %w", err)
		}

		// Metering? (TODO: what is this for?)
		metering, err := dmetering.New(viper.GetString("common-metering-plugin"))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize dmetering: %w", err)
		}
		dmetering.SetDefaultMeter(metering)

		// Configure firehose data sources
		firehoseBlocksStoreURLs := viper.GetStringSlice("firehose-blocks-store-urls")
		if len(firehoseBlocksStoreURLs) == 0 {
			firehoseBlocksStoreURLs = []string{viper.GetString("common-blocks-store-url")}
		} else if len(firehoseBlocksStoreURLs) == 1 && strings.Contains(firehoseBlocksStoreURLs[0], ",") {
			firehoseBlocksStoreURLs = strings.Split(firehoseBlocksStoreURLs[0], ",")
		}
		for i, url := range firehoseBlocksStoreURLs {
			firehoseBlocksStoreURLs[i] = mustReplaceDataDir(sfDataDir, url)
		}

		// Configure graceful shutdown
		shutdownSignalDelay := viper.GetDuration("common-system-shutdown-signal-delay")
		grcpShutdownGracePeriod := time.Duration(0)
		if shutdownSignalDelay.Seconds() > 5 {
			grcpShutdownGracePeriod = shutdownSignalDelay - (5 * time.Second)
		}

		return firehoseApp.New(appLogger,
			&firehoseApp.Config{
				BlockStoreURLs:          firehoseBlocksStoreURLs,
				BlockStreamAddr:         blockstreamAddr,
				GRPCListenAddr:          viper.GetString("firehose-grpc-listen-addr"),
				GRPCShutdownGracePeriod: grcpShutdownGracePeriod,
				RealtimeTolerance:       viper.GetDuration("firehose-real-time-tolerance"),
			},
			&firehoseApp.Modules{
				Authenticator:             authenticator,
				FilterPreprocessorFactory: codec.FilterPreprocessorFactory(),
				HeadTimeDriftMetric:       headTimeDriftmetric,
				HeadBlockNumberMetric:     headBlockNumMetric,
				Tracker:                   tracker,
			}), nil
	}

	launcher.RegisterApp(&launcher.AppDef{
		ID:            "firehose",
		Title:         "Block Firehose",
		Description:   "Provides on-demand filtered blocks, depends on common-blocks-store-url and common-blockstream-addr",
		MetricsID:     "firehose",
		Logger:        launcher.NewLoggingDef("firehose.*", nil),
		RegisterFlags: registerFlags,
		FactoryFunc:   factoryFunc,
	})
}
