package tcp

import (
	"errors"

	netreactors "github.com/lance-e/net-reactors"
)

const (
	MessageHeader int = 4 //消息头为4字节
)

const (
	TcpPackErrorWrongHeader     = "wrong header"
	TcpPackErrorNotComplete     = "message not complete"
	TcpPackErrorBytesLessHeader = "readable bytes not enough to unpack"
)

// tcp
func Unpacking(buf *netreactors.Buffer) ([]byte, error) {
	if buf.ReadableBytes() < MessageHeader {
		return nil, errors.New(TcpPackErrorBytesLessHeader)
	}
	//读出消息头
	header := buf.RetrieveAsString(MessageHeader)
	length := uint32(header[0])<<24 | uint32(header[1])<<16 | uint32(header[2])<<8 | uint32(header[3])
	if length < 0 || length > 65536 {
		//异常消息头
		return nil, errors.New(TcpPackErrorWrongHeader)
	} else if buf.ReadableBytes() >= int(length) {
		//足够读取出一条消息
		return buf.RetrieveAsString(int(length)), nil
	} else {
		//剩余数据不够一条消息，说明tcp拆包了.把读出的消息头拼接回buf,等待下一次消息读取
		buf.Prepend(header)
		return nil, errors.New(TcpPackErrorNotComplete)
	}
}

// tcp
func Packing(data []byte) *netreactors.Buffer {
	buf := netreactors.NewBuffer()
	buf.Append(data)

	//prepend header
	length := len(data)
	header := make([]byte, 4)
	header[0] = byte(length >> 24)
	header[1] = byte(length >> 16)
	header[2] = byte(length >> 8)
	header[3] = byte(length)
	buf.Prepend(header)

	return buf
}
