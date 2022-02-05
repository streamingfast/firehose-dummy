.PHONY: build
build:
	go build -o ./build/firehose-dummy ./cmd/firehose-dummy

.PHONY: setup
setup:
	go mod download

.PHONY: clean
clean:
	@rm -rf ./tmp/
	@rm -rf ./data/
	@rm -rf ./sf-data/
	@rm -rf ./build

.PHONY: update-stack-deps
update-stack-deps:
	go get github.com/streamingfast/bstream@develop
	go get github.com/streamingfast/dauth@develop
	go get github.com/streamingfast/dbin@develop
	go get github.com/streamingfast/derr@develop
	go get github.com/streamingfast/dlauncher@develop
	go get github.com/streamingfast/dmetering@develop
	go get github.com/streamingfast/dstore@develop
	go get github.com/streamingfast/firehose@develop
	go get github.com/streamingfast/merger@develop
	go get github.com/streamingfast/pbgo@develop
	go get github.com/streamingfast/relayer@develop
	go get github.com/streamingfast/node-manager@develop
	go get github.com/streamingfast/shutter@develop
	go mod tidy