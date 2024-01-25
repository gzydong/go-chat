package encoding

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"sync"
)

var bufferPool = sync.Pool{
	New: func() any {
		return &bytes.Buffer{}
	},
}

// NewEncode 将消息编码
//
//	[x][x][x][x][x][x][x][x]...
//	|  (int32) || (binary)
//	|  4-byte  || N-byte
//	------------------------...
//	    size       data
func NewEncode(data []byte) ([]byte, error) {

	buf := bufferPool.Get().(*bytes.Buffer)

	// 写入消息头
	// 读取消息的长度，转换成int32类型（占4个字节）
	var length = int32(len(string(data)))
	if err := binary.Write(buf, binary.LittleEndian, length); err != nil {
		return nil, err
	}

	// 写入消息实体
	if err := binary.Write(buf, binary.LittleEndian, data); err != nil {
		return nil, err
	}

	buffer := buf.Bytes()
	buf.Reset()
	bufferPool.Put(buf)

	return buffer, nil
}

// NewDecode 从缓冲区里读取数据
func NewDecode(r io.Reader) ([]byte, error) {
	var length int32

	// message size
	err := binary.Read(r, binary.LittleEndian, &length)
	if err != nil {
		return nil, err
	}

	if length < 0 {
		return nil, fmt.Errorf("response msg size is negative: %v", length)
	}

	// message binary data
	buf := make([]byte, length)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, err
	}

	return buf, nil
}
