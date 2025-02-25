package tcp_test

import (
	"nbim/pkg/tcp"
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
	buf := tcp.Packing(data)
	msg, err := tcp.Unpacking(buf)
	assert.NoError(t, err)
	assert.Equal(t, []byte("hello"), msg)

	// 消息头长度不足
	data = []byte{0x00, 0x00, 0x00}
	buf = netreactors.NewBuffer()
	buf.Append(data)
	msg, err = tcp.Unpacking(buf)
	assert.Error(t, err, "readable bytes not enough to unpack")
	assert.Nil(t, msg)

	// 消息体长度不足
	data = []byte{0x00, 0x00, 0x00, 0x05, 'h', 'e', 'l', 'l'}
	buf = netreactors.NewBuffer()
	buf.Append(data)
	msg, err = tcp.Unpacking(buf)
	assert.Error(t, err, "message not complete")
	assert.Nil(t, msg)

	// 消息长度异常
	data = []byte{0xff, 0xff, 0xff, 0xff} // 长度过大
	buf = netreactors.NewBuffer()
	buf.Append(data)
	msg, err = tcp.Unpacking(buf)
	assert.Error(t, err, "wrong header")
	assert.Nil(t, msg)

	data = []byte{0x80, 0x00, 0x00, 0x00} // 最高位为 1 ,最大
	buf = netreactors.NewBuffer()
	buf.Append(data)
	msg, err = tcp.Unpacking(buf)
	assert.Error(t, err, "wrong header")
	assert.Nil(t, msg)

	// 空消息
	data = []byte{0x00, 0x00, 0x00, 0x00}
	buf = netreactors.NewBuffer()
	buf.Append(data)
	msg, err = tcp.Unpacking(buf)
	assert.NoError(t, err)
	assert.Equal(t, []byte(""), msg)

	// 多个消息包
	data = []byte{0x00, 0x00, 0x00, 0x05, 'h', 'e', 'l', 'l', 'o', 0x00, 0x00, 0x00, 0x03, 'w', 'o'}
	buf = netreactors.NewBuffer()
	buf.Append(data)

	msg, err = tcp.Unpacking(buf)
	assert.NoError(t, err)
	assert.Equal(t, []byte("hello"), msg)

	msg, err = tcp.Unpacking(buf)
	assert.Error(t, err, "message not complete")
	assert.Nil(t, msg)

}

func TestPacking(t *testing.T) {
	data := []byte("hello")
	buf := tcp.Packing(data)
	expected := []byte{0x00, 0x00, 0x00, 0x05, 'h', 'e', 'l', 'l', 'o'}
	assert.Equal(t, expected, buf.RetrieveAllString())
}
