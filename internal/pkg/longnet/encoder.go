package longnet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

type FlagPos uint

const (
	FlagCompressed FlagPos = 0 // 压缩标志
	FlagEncrypted  FlagPos = 1 // 加密标志
	FlagNeedAck    FlagPos = 2 // 是否需要 ACK 回执
	FlagReserve    FlagPos = 3
)

// 封包后的结构
// +-------------------+-------------------+-------------------+-----------------------------------+
// |      Command      |       Flags       |       MsgID       |             Payload              |
// |     (4 bytes)     |     (4 bytes)     |     (8 bytes)     |        (Length bytes)            |
// +-------------------+-------------------+-------------------+-----------------------------------+

var _ IEncoder = (*Encoder)(nil)

type EncoderOptions struct {
	MaxPacketSize         int // 最大数据包大小
	MinCompressPacketSize int // 压缩阈值
}

func (e EncoderOptions) init() EncoderOptions {
	if e.MaxPacketSize <= 0 {
		e.MaxPacketSize = PacketMaxSize
	}

	if e.MinCompressPacketSize <= 0 {
		e.MinCompressPacketSize = 1024 * 500 // 500KB
	}

	return e
}

type Encoder struct {
	options  EncoderOptions
	compress ICompress
	encrypt  IEncrypter
}

func NewEncoder(options EncoderOptions, compress ICompress, encrypt IEncrypter) *Encoder {
	return &Encoder{
		options:  options.init(),
		compress: compress,
		encrypt:  encrypt,
	}
}

func (e *Encoder) Pack(m *Packet) ([]byte, error) {
	var buf bytes.Buffer
	var err error
	var body = m.Payload

	if len(body) > e.options.MaxPacketSize {
		return nil, fmt.Errorf("message too large: %d", len(body))
	}

	if err := binary.Write(&buf, binary.BigEndian, m.Cmd); err != nil {
		return nil, err
	}

	var flags Flag
	if e.compress != nil && len(body) >= e.options.MinCompressPacketSize {
		body, err = e.compress.Compress(body)
		if err != nil {
			return nil, err
		}

		flags.SetBit(FlagCompressed)
	}

	// 是否加密
	if e.encrypt != nil {
		body, err = e.encrypt.Encrypt(body)
		if err != nil {
			return nil, err
		}

		flags.SetBit(FlagEncrypted)
	}

	if err := binary.Write(&buf, binary.BigEndian, flags); err != nil {
		return nil, err
	}

	if err := binary.Write(&buf, binary.BigEndian, m.Msgid); err != nil {
		return nil, err
	}

	if _, err := buf.Write(body); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (e *Encoder) UnPack(data []byte) (*Packet, error) {
	buf := bytes.NewReader(data)
	packet := &Packet{}

	if err := binary.Read(buf, binary.BigEndian, &packet.Cmd); err != nil {
		return nil, err
	}

	var flags Flag
	if err := binary.Read(buf, binary.BigEndian, &flags); err != nil {
		return nil, err
	}

	var msgID uint64
	if err := binary.Read(buf, binary.BigEndian, &msgID); err != nil {
		return nil, err
	}

	packet.Msgid = int64(msgID)

	packet.Payload = make([]byte, len(data)-16)
	if _, err := buf.Read(packet.Payload); err != nil {
		return nil, err
	}

	// 注意顺序：先解密，后解压（与 Encrypt 顺序相反）

	// 1. 解密
	if flags.HasBit(FlagEncrypted) {
		if e.encrypt == nil {
			return nil, errors.New("encrypt is nil")
		}

		body, err := e.encrypt.Decrypt(packet.Payload)
		if err != nil {
			return nil, err
		}
		packet.Payload = body
	}

	// 2. 解压
	if flags.HasBit(FlagCompressed) {
		if e.compress == nil {
			return nil, errors.New("compress is nil")
		}

		var err error
		packet.Payload, err = e.compress.Decompress(packet.Payload)
		if err != nil {
			return nil, err
		}
	}

	return packet, nil
}

type Flag uint32

// SetBit 设置指定位置的 bit
func (f *Flag) SetBit(pos FlagPos) {
	*f |= 1 << pos
}

// HasBit 判断指定位置的 bit 是否为 1
func (f *Flag) HasBit(pos FlagPos) bool {
	return (*f)&(1<<pos) != 0
}
