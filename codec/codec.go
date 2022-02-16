package codec

import (
	"errors"
	"io"
	"time"

	"github.com/streamingfast/bstream"
	pbbstream "github.com/streamingfast/pbgo/sf/bstream/v1"
)

const (
	ProtocolNum = pbbstream.Protocol_UNKNOWN
)

func init() {
	bstream.GetBlockReaderFactory = bstream.BlockReaderFactoryFunc(BlockReaderFactory)
	bstream.GetBlockDecoder = bstream.BlockDecoderFunc(blockDecoder)
	bstream.GetBlockWriterFactory = bstream.BlockWriterFactoryFunc(BlockWriterFactory)
	bstream.GetBlockWriterHeaderLen = 10
	bstream.GetBlockPayloadSetter = bstream.MemoryBlockPayloadSetter
	bstream.GetMemoizeMaxAge = 200 * 15 * time.Second
}

func Validate() error {
	if err := bstream.ValidateRegistry(); err != nil {
		return err
	}

	if bstream.GetProtocolFirstStreamableBlock == 0 {
		return errors.New("protocol first streamable block must be set")
	}

	return nil
}

func SetProtocolFirstStreamableBlock(height uint64) {
	bstream.GetProtocolFirstStreamableBlock = height
}

func BlockReaderFactory(reader io.Reader) (bstream.BlockReader, error) {
	return NewBlockReader(reader)
}

func BlockWriterFactory(writer io.Writer) (bstream.BlockWriter, error) {
	return NewBlockWriter(writer)
}
