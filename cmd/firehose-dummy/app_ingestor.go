package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/streamingfast/bstream"
	"github.com/streamingfast/bstream/blockstream"
	"github.com/streamingfast/dgrpc"
	"github.com/streamingfast/dlauncher/launcher"
	"github.com/streamingfast/logging"
	"github.com/streamingfast/node-manager/mindreader"

	pbbstream "github.com/streamingfast/pbgo/sf/bstream/v1"
	pbheadinfo "github.com/streamingfast/pbgo/sf/headinfo/v1"
	"github.com/streamingfast/shutter"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	pbcodec "github.com/streamingfast/dummy-blockchain/proto"
	"github.com/streamingfast/firehose-dummy/codec"
)

const (
	defaultLineBufferSize = 10 * 1024 * 1024

	modeLogs  = "logs"  // Consume events from the log file(s)
	modeStdin = "stdin" // Consume events from the STDOUT of another process
	modeNode  = "node"  // Consume events from the spawned node process
)

func init() {
	appLogger := zap.NewNop()
	logging.Register("ingestor", &appLogger)

	registerFlags := func(cmd *cobra.Command) error {
		flags := cmd.Flags()

		flags.String("ingestor-mode", modeStdin, "Mode of operation, one of (stdin, logs, node)")
		flags.String("ingestor-logs-dir", "", "Event logs source directory")
		flags.Int("ingestor-line-buffer-size", defaultLineBufferSize, "Buffer size in bytes for the line reader")
		flags.String("ingestor-working-dir", "{sf-data-dir}/workdir", "Path where mindreader will stores its files")
		flags.String("ingestor-grpc-listen-addr", BlockStreamServingAddr, "GRPC server listen address")
		flags.Duration("ingestor-merge-threshold-block-age", time.Duration(math.MaxInt64), "When processing blocks with a blocktime older than this threshold, they will be automatically merged")
		flags.String("ingestor-node-path", "", "Path to node binary")
		flags.String("ingestor-node-dir", "", "Node working directory")
		flags.String("ingestor-node-args", "", "Node process arguments")

		return nil
	}

	initFunc := func(runtime *launcher.Runtime) (err error) {
		mode := viper.GetString("ingestor-mode")

		switch mode {
		case modeStdin:
			return nil
		case modeNode:
			return checkNodeBinPath(viper.GetString("ingestor-node-path"))
		case modeLogs:
			return checkLogsSource(viper.GetString("ingestor-logs-dir"))
		default:
			return fmt.Errorf("invalid mode: %v", mode)
		}
	}

	factoryFunc := func(runtime *launcher.Runtime) (launcher.App, error) {
		sfDataDir := runtime.AbsDataDir

		oneBlockStoreURL := mustReplaceDataDir(sfDataDir, viper.GetString("common-oneblock-store-url"))
		mergedBlockStoreURL := mustReplaceDataDir(sfDataDir, viper.GetString("common-blocks-store-url"))
		workingDir := mustReplaceDataDir(sfDataDir, viper.GetString("ingestor-working-dir"))
		gprcListenAdrr := viper.GetString("ingestor-grpc-listen-addr")
		mergeAndStoreDirectly := viper.GetBool("ingestor-merge-and-store-directly")
		mergeThresholdBlockAge := viper.GetDuration("ingestor-merge-threshold-block-age")
		batchStartBlockNum := viper.GetUint64("ingestor-start-block-num")
		batchStopBlockNum := viper.GetUint64("ingestor-stop-block-num")
		waitTimeForUploadOnShutdown := viper.GetDuration("ingestor-wait-upload-complete-on-shutdown")
		oneBlockFileSuffix := viper.GetString("ingestor-oneblock-suffix")
		blocksChanCapacity := viper.GetInt("ingestor-blocks-chan-capacity")

		tracker := bstream.NewTracker(50) // TODO: make a flag
		tracker.AddResolver(bstream.OffsetStartBlockResolver(1))

		consoleReaderFactory := func(lines chan string) (mindreader.ConsolerReader, error) {
			return codec.NewLogReader(lines, "")
		}

		consoleReaderTransformer := func(obj interface{}) (*bstream.Block, error) {
			return codec.BlockFromProto(obj.(*pbcodec.Block))
		}

		blockStreamServer := blockstream.NewUnmanagedServer(blockstream.ServerOptionWithLogger(appLogger))
		healthCheck := func(ctx context.Context) (isReady bool, out interface{}, err error) {
			return blockStreamServer.Ready(), nil, nil
		}

		server := dgrpc.NewServer2(
			dgrpc.WithLogger(appLogger),
			dgrpc.WithHealthCheck(dgrpc.HealthCheckOverGRPC|dgrpc.HealthCheckOverHTTP, healthCheck),
		)
		server.RegisterService(func(gs *grpc.Server) {
			pbheadinfo.RegisterHeadInfoServer(gs, blockStreamServer)
			pbbstream.RegisterBlockStreamServer(gs, blockStreamServer)
		})

		mrp, err := mindreader.NewMindReaderPlugin(
			oneBlockStoreURL,
			mergedBlockStoreURL,
			mergeAndStoreDirectly,
			mergeThresholdBlockAge,
			workingDir,
			consoleReaderFactory,
			consoleReaderTransformer,
			tracker,
			batchStartBlockNum,
			batchStopBlockNum,
			blocksChanCapacity,
			headBlockUpdater,
			func(error) {},
			true,
			waitTimeForUploadOnShutdown,
			oneBlockFileSuffix,
			blockStreamServer,
			appLogger,
		)
		if err != nil {
			log.Fatal("error initialising mind reader", zap.Error(err))
			return nil, nil
		}

		return &IngestorApp{
			Shutter:          shutter.New(),
			mrp:              mrp,
			mode:             viper.GetString("ingestor-mode"),
			lineBufferSize:   viper.GetInt("ingestor-line-buffer-size"),
			nodeBinPath:      viper.GetString("ingestor-node-path"),
			nodeDir:          viper.GetString("ingestor-node-dir"),
			nodeArgs:         viper.GetString("ingestor-node-args"),
			server:           server,
			serverListenAddr: gprcListenAdrr,
		}, nil
	}

	launcher.RegisterApp(&launcher.AppDef{
		ID:            "ingestor",
		Title:         "Ingestor",
		Description:   "Reads the log files produces by the instrumented node",
		MetricsID:     "ingestor",
		Logger:        launcher.NewLoggingDef("ingestor.*", nil),
		RegisterFlags: registerFlags,
		InitFunc:      initFunc,
		FactoryFunc:   factoryFunc,
	})
}

func headBlockUpdater(uint64, string, time.Time) {
	// TODO: will need to be implemented somewhere
}

func checkLogsSource(dir string) error {
	if dir == "" {
		return errors.New("ingestor logs dir must be set")
	}

	dir, err := expandDir(dir)
	if err != nil {
		return err
	}

	if !dirExists(dir) {
		return errors.New("ingestor logs dir must exist")
	}

	return nil
}

func checkNodeBinPath(binPath string) error {
	if binPath == "" {
		return errors.New("node path must be set")
	}

	stat, err := os.Stat(binPath)
	if err != nil {
		return fmt.Errorf("cant inspect node path: %w", err)
	}

	if stat.IsDir() {
		return fmt.Errorf("path %v is a directory", binPath)
	}

	return nil
}
