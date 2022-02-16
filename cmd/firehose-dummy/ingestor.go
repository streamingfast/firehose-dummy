package main

import (
	"context"
	"os"
	"strings"

	"github.com/streamingfast/dgrpc"
	"github.com/streamingfast/firehose-dummy/noderunner"
	"github.com/streamingfast/node-manager/mindreader"
	"github.com/streamingfast/shutter"
	"go.uber.org/zap"
)

type IngestorApp struct {
	*shutter.Shutter

	mode             string
	lineBufferSize   int
	logsDir          string
	serverListenAddr string
	nodeBinPath      string
	nodeDir          string
	nodeArgs         string

	mrp    *mindreader.MindReaderPlugin
	server *dgrpc.Server
}

func (app *IngestorApp) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Configure shutdown flow
	app.OnTerminating(app.mrp.Shutdown)
	app.mrp.OnTerminated(app.Shutdown)

	zlog.Info("starting ingestor", zap.String("mode", app.mode))
	defer zlog.Info("ingestor stopped")

	zlog.Info("starting ingestor blockstream server")
	go app.server.Launch(app.serverListenAddr)

	zlog.Info("starting ingestor reader plugin")
	go app.mrp.Launch()

	go func() {
		var err error

		switch app.mode {
		case modeStdin:
			err = app.startFromStdin(ctx)
		case modeNode:
			err = app.startFromNode(ctx)
		case modeLogs:
			err = app.startFromLogs(ctx)
		}

		zlog.Info("reader finished", zap.Error(err))
		app.mrp.Shutdown(err)
	}()

	<-app.Terminated()
	return app.Err()
}

func (app *IngestorApp) startFromStdin(ctx context.Context) error {
	return noderunner.StartLineReader(os.Stdin, app.mrp.LogLine, zlog)
}

func (app *IngestorApp) startFromNode(ctx context.Context) error {
	args := strings.Split(app.nodeArgs, " ")

	runner := noderunner.New(app.nodeBinPath, args, true)
	runner.SetLogger(zlog)
	runner.SetLineReader(app.mrp.LogLine)
	runner.SetDir(app.nodeDir)

	return runner.Start(ctx)
}

func (app *IngestorApp) startFromLogs(ctx context.Context) error {
	return nil
}
