package codec

import (
	"fmt"

	"github.com/streamingfast/bstream"
	"google.golang.org/protobuf/proto"

	pbcodec "github.com/streamingfast/dummy-blockchain/proto"
)

func blockDecoder(blk *bstream.Block) (interface{}, error) {
	if blk.Kind() != ProtocolNum {
		return nil, fmt.Errorf("expected kind %s, got %s", ProtocolNum, blk.Kind())
	}

	if blk.Version() != 1 {
		return nil, fmt.Errorf("this decoder only knows about version 1, got %d", blk.Version())
	}

	payload, err := blk.Payload.Get()
	if err != nil {
		return nil, fmt.Errorf("unable to get payload from block stream data: %v", err)
	}

	block := &pbcodec.Block{}
	if err := proto.Unmarshal(payload, block); err != nil {
		return nil, fmt.Errorf("unable to decode block stream data: %v", err)
	}

	return block, nil
}
