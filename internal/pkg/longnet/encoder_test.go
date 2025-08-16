package longnet

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncoder_Encode(t *testing.T) {
	encoder := NewEncoder(EncoderOptions{}, &SnappyCompress{}, nil)
	message := NewCustomizePacket(1001, []byte("hello world"))
	body, err := encoder.Pack(message)
	fmt.Println(body)
	assert.NoError(t, err)
}

func TestEncoder_EncodeDecode(t *testing.T) {
	compressor := &SnappyCompress{}
	encoder := NewEncoder(EncoderOptions{}, compressor, nil)

	t.Run("no compression", func(t *testing.T) {
		msg := &Packet{
			Cmd:     1001,
			Payload: []byte("plain text"),
			Msgid:   1234567890,
		}

		encoded, err := encoder.Pack(msg)
		assert.NoError(t, err)

		decoded, err := encoder.UnPack(encoded)
		assert.NoError(t, err)

		assert.Equal(t, msg.Cmd, decoded.Cmd)
		assert.Equal(t, msg.Msgid, decoded.Msgid)
		assert.Equal(t, msg.Payload, decoded.Payload)
	})

	t.Run("with compression", func(t *testing.T) {
		msg := &Packet{
			Cmd:     1002,
			Payload: []byte("this is a test message for compression"),
			Msgid:   9876543210,
		}

		encoded, err := encoder.Pack(msg)
		assert.NoError(t, err)

		decoded, err := encoder.UnPack(encoded)
		assert.NoError(t, err)
		assert.Equal(t, msg.Cmd, decoded.Cmd)
		assert.Equal(t, msg.Msgid, decoded.Msgid)
		assert.True(t, bytes.Equal(msg.Payload, decoded.Payload))
	})
}
