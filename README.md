# firehose-dummy

Firehose CLI for the [Dummy Chain](https://github.com/streamingfast/dummy-blockchain)

## Building

Clone the repository:

```bash
git clone git@github.com:streamingfast/firehose-dummy.git
```

Install dependencies:

```bash
go mod download
```

Then build the binary:

```bash
make build
```

## Usage

To see usage example, run: `./build/firehose-dummy --help`:

```
Dummy Chain Firehose

Usage:
  firehose-dummy [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  init        Initialize local configuration
  reset       Reset local state
  setup       Configures and initializes the project files
  start       Starts all services at once

Flags:
  -h, --help   help for firehose-dummy

Use "firehose-dummy [command] --help" for more information about a command.
```

## Configuration

All available configuration options for the `start` command:

```
Starts all services at once

Usage:
  firehose-dummy start [flags]

Flags:
      --common-auth-plugin string                        Auth plugin URI, see streamingfast/dauth repository (default "null://")
      --common-blocks-store-url string                   Store URL (with prefix) where to read/write (default "file://{sf-data-dir}/storage/merged-blocks")
      --common-blockstream-addr string                   GRPC endpoint to get real-time blocks (default ":9010")
      --common-first-streamable-block uint               First streamable block number (default 1)
      --common-metering-plugin string                    Metering plugin URI, see streamingfast/dmetering repository (default "null://")
      --common-oneblock-store-url string                 Store URL (with prefix) to read/write one-block files (default "file://{sf-data-dir}/storage/one-blocks")
      --common-shutdown-delay duration                   Add a delay between receiving SIGTERM signal and shutting down apps. Apps will respond negatively to /healthz during this period (default 5ns)
      --firehose-blocks-store-urls strings               If non-empty, overrides common-blocks-store-url with a list of blocks stores
      --firehose-grpc-listen-addr string                 Address on which the firehose will listen (default ":9030")
      --firehose-real-time-tolerance duration            firehose will became alive if now - block time is smaller then tolerance (default 1m0s)
  -h, --help                                             help for start
      --ingestor-grpc-listen-addr string                 GRPC server listen address (default ":9000")
      --ingestor-line-buffer-size int                    Buffer size in bytes for the line reader (default 10485760)
      --ingestor-logs-dir string                         Event logs source directory
      --ingestor-merge-threshold-block-age duration      When processing blocks with a blocktime older than this threshold, they will be automatically merged (default 2562047h47m16.854775807s)
      --ingestor-mode string                             Mode of operation, one of (stdin, logs, node) (default "stdin")
      --ingestor-node-args string                        Node process arguments
      --ingestor-node-dir string                         Node working directory
      --ingestor-node-path string                        Path to node binary
      --ingestor-working-dir string                      Path where mindreader will stores its files (default "{sf-data-dir}/workdir")
      --log-format string                                Logging format (default "text")
      --merger-grpc-listen-addr string                   Address to listen for incoming gRPC requests (default ":9020")
      --merger-max-one-block-operations-batch-size int   max number of 'good' (mergeable) files to look up from storage in one polling operation (default 2000)
      --merger-next-exclusive-highest-block-limit int    for next bundle boundary
      --merger-one-block-deletion-threads int            number of parallel threads used to delete one-block-files (more means more stress on your storage backend) (default 10)
      --merger-state-file string                         Path to file containing last written block number, as well as a map of all 'seen blocks' in the 'max-fixable-fork' range (default "{sf-data-dir}/merger/merger.seen.gob")
      --merger-time-between-store-lookups duration       delay between source store polling (should be higher for remote storage) (default 5s)
      --merger-writers-leeway duration                   how long we wait after seeing the upper boundary, to ensure that we get as many blocks as possible in a bundle (default 10s)
      --relayer-buffer-size int                          Number of blocks that will be kept and sent immediately on connection (default 350)
      --relayer-grpc-listen-addr string                  Address to listen for incoming gRPC requests (default ":9010")
      --relayer-max-source-latency duration              Max latency tolerated to connect to a source (default 10m0s)
      --relayer-merger-addr string                       Address for grpc merger service (default ":9020")
      --relayer-min-start-offset uint                    Number of blocks before HEAD where we want to start for faster buffer filling (missing blocks come from files/merger) (default 120)
      --relayer-source strings                           List of Blockstream sources (mindreaders) to connect to for live block feeds (repeat flag as needed) (default [:9000])
      --relayer-source-request-burst int                 Block burst requested by relayer (useful when chaining relayers together, because normally a mindreader won't have a block buffer) (default 90)
      --verbose int                                      Logging verbosity (default 1)
```

## Running

Firehose CLI supports running Ingestor (event log consumer) in multiple modes.

### Node exec mode

In this mode, the Ingestor component will spawn the blockchain node subprocess and
consume all events from the STDOUT.

Example:

```bash
firehose-dummy start ingestor \
  --ingestor-mode=node \
  --ingestor-node-path=/path/to/blockchain/bin \
  --ingestor-node-dir=/path/to/blockchain/home \
  --ingestor-node-args="start" \
```

### Standard input mode

In this mode, the Ingestor component will consule blockchain event logs from the 
STDIN. Source could be another process, or a pipe.

Example:

```bash
dummy_blockchain start | firehose-dummy start ingestor --ingestor-mode=stdin
```

### Logs mode

TODO

## Contributors

- [Figment](https://github.com/figment-networks): Initial Implementation

## License

TBD
