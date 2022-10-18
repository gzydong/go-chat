package sdk

import (
	"bufio"
	"bytes"
	"encoding/binary"
)

// Encode 将消息编码
func Encode(message string) ([]byte, error) {

	var pkg = new(bytes.Buffer)

	// 写入消息头
	// 读取消息的长度，转换成int32类型（占4个字节）
	var length = int32(len(message))
	err := binary.Write(pkg, binary.LittleEndian, length)
	if err != nil {
		return nil, err
	}

	// 写入消息实体
	err = binary.Write(pkg, binary.LittleEndian, []byte(message))
	if err != nil {
		return nil, err
	}

	return pkg.Bytes(), nil
}

// Decode 解码消息
func Decode(reader *bufio.Reader) (string, error) {
	// 读取消息的长度
	lengthByte, _ := reader.Peek(4) // 读取前4个字节的数据
	lengthBuff := bytes.NewBuffer(lengthByte)

	var length int32
	err := binary.Read(lengthBuff, binary.LittleEndian, &length)
	if err != nil {
		return "", err
	}

	// Buffered返回缓冲中现有的可读取的字节数。
	if int32(reader.Buffered()) < length+4 {
		return "", err
	}

	// 读取真正的消息数据
	pack := make([]byte, int(4+length))
	if _, err = reader.Read(pack); err != nil {
		return "", err
	}

	return string(pack[4:]), nil
}
