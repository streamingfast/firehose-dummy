package codec

import (
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/streamingfast/bstream"

	pbcodec "github.com/streamingfast/dummy-blockchain/proto"
)

func BlockFromProto(b *pbcodec.Block) (*bstream.Block, error) {
	content, err := proto.Marshal(b)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal to binary form: %s", err)
	}

	block := &bstream.Block{
		Id:             b.Hash,
		Number:         b.Height,
		PreviousId:     b.PrevHash,
		Timestamp:      time.Unix(0, int64(b.Timestamp)).UTC(),
		LibNum:         b.Height - 1,
		PayloadKind:    ProtocolNum,
		PayloadVersion: 1,
	}

	if block.Number == bstream.GetProtocolFirstStreamableBlock {
		block.LibNum = bstream.GetProtocolFirstStreamableBlock
		block.PreviousId = ""
	}

	return bstream.GetBlockPayloadSetter(block, content)
}
