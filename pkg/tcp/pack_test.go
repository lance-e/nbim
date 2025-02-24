package tcp_test

import (
	"errors"
	"testing"

	netreactors "github.com/lance-e/net-reactors"
	"github.com/stretchr/testify/assert"
)

const (
	MessageHeader int = 4 //消息头为4字节
)

func TestUnpacking(t *testing.T) {
	// 正常情况
	data := []byte{'h', 'e', 'l', 'l', 'o'}
	buf := Packing(data)
	msg, err := Unpacking(buf)
	assert.NoError(t, err)
	assert.Equal(t, []byte("hello"), msg)

	// 消息头长度不足
	data = []byte{0x00, 0x00, 0x00}
	buf = netreactors.NewBuffer()
	buf.Append(data)
	msg, err = Unpacking(buf)
	assert.Error(t, err, "readable bytes not enough to unpack")
	assert.Nil(t, msg)

	// 消息体长度不足
	data = []byte{0x00, 0x00, 0x00, 0x05, 'h', 'e', 'l', 'l'}
	buf = netreactors.NewBuffer()
	buf.Append(data)
	msg, err = Unpacking(buf)
	assert.Error(t, err, "message not complete")
	assert.Nil(t, msg)

	// 消息长度异常
	data = []byte{0xff, 0xff, 0xff, 0xff} // 长度过大
	buf = netreactors.NewBuffer()
	buf.Append(data)
	msg, err = Unpacking(buf)
	assert.Error(t, err, "wrong header")
	assert.Nil(t, msg)

	data = []byte{0x80, 0x00, 0x00, 0x00} // 最高位为 1 ,最大
	buf = netreactors.NewBuffer()
	buf.Append(data)
	msg, err = Unpacking(buf)
	assert.Error(t, err, "wrong header")
	assert.Nil(t, msg)

	// 空消息
	data = []byte{0x00, 0x00, 0x00, 0x00}
	buf = netreactors.NewBuffer()
	buf.Append(data)
	msg, err = Unpacking(buf)
	assert.NoError(t, err)
	assert.Equal(t, []byte(""), msg)

	// 多个消息包
	data = []byte{0x00, 0x00, 0x00, 0x05, 'h', 'e', 'l', 'l', 'o', 0x00, 0x00, 0x00, 0x03, 'w', 'o'}
	buf = netreactors.NewBuffer()
	buf.Append(data)

	msg, err = Unpacking(buf)
	assert.NoError(t, err)
	assert.Equal(t, []byte("hello"), msg)

	msg, err = Unpacking(buf)
	assert.Error(t, err, "message not complete")
	assert.Nil(t, msg)

}

func TestPacking(t *testing.T) {
	data := []byte("hello")
	buf := Packing(data)
	expected := []byte{0x00, 0x00, 0x00, 0x05, 'h', 'e', 'l', 'l', 'o'}
	assert.Equal(t, expected, buf.RetrieveAllString())
}

// tcp
func Unpacking(buf *netreactors.Buffer) ([]byte, error) {
	if buf.ReadableBytes() < MessageHeader {
		return nil, errors.New("readable bytes not enough to unpack")
	}
	//读出消息头
	header := buf.RetrieveAsString(MessageHeader)
	length := uint32(header[0])<<24 | uint32(header[1])<<16 | uint32(header[2])<<8 | uint32(header[3])
	if length < 0 || length > 65536 {
		//异常消息头
		return nil, errors.New("wrong header")
	} else if buf.ReadableBytes() >= int(length) {
		//足够读取出一条消息
		return buf.RetrieveAsString(int(length)), nil
	} else {
		//剩余数据不够一条消息，说明tcp拆包了.把读出的消息头拼接回buf,等待下一次消息读取
		buf.Prepend(header)
		return nil, errors.New("message not complete")
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
