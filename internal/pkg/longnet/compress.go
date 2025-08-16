package longnet

import (
	"github.com/golang/snappy"
)

var _ ICompress = (*SnappyCompress)(nil)

type SnappyCompress struct {
}

func NewSnappyCompress() ICompress {
	return &SnappyCompress{}
}

func (s *SnappyCompress) Compress(data []byte) ([]byte, error) {
	return snappy.Encode(nil, data), nil
}

func (s *SnappyCompress) Decompress(data []byte) ([]byte, error) {
	return snappy.Decode(nil, data)
}
