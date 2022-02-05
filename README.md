# sf-chain

StreamingFast CLI for the Dummy Chain

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
Usage:
  firehose-dummy [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  init        Initialize local configuration
  setup       Configures and initializes the project files
  start       Starts all services at once

Flags:
  -h, --help   help for firehose-dummy

Use "firehose-dummy [command] --help" for more information about a command.
```

## Running

TBD

## Contributors

- [Figment](https://github.com/figment-networks): Initial Implementation

## License

TBD
