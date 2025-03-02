package sdk

import (
	"encoding/json"
	"fmt"
	"nbim/pkg/protocol/pb"
	"nbim/pkg/tcp"
	"net"
	"sync/atomic"
	"time"

	"google.golang.org/protobuf/proto"
)

type connect struct {
	sendChan, recvChan chan *Message
	conn               *net.TCPConn
	connID             int64
	ip                 net.IP
	port               int
}

func newConnect(ip net.IP, port int) *connect {
	clientConn := &connect{
		sendChan: make(chan *Message),
		recvChan: make(chan *Message),
		ip:       ip,
		port:     port,
	}
	addr := &net.TCPAddr{IP: ip, Port: port}
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		fmt.Printf("DialTCP.err=%+v", err)
		return nil
	}
	clientConn.conn = conn
	return clientConn
}

func handAckMsg(c *connect, data []byte) *Message {
	ackMsg := &pb.AckMsg{}
	proto.Unmarshal(data, ackMsg)
	switch ackMsg.ToType {
	case pb.CMD_Login, pb.CMD_Reconn:
		atomic.StoreInt64(&c.connID, ackMsg.ConnId)
	}
	return &Message{
		Type:       MsgTypeAck,
		Name:       "nbim",
		FormUserID: "1212121",
		ToUserID:   "222212122",
		Content:    ackMsg.GetMessage(),
	}
}
func handPushMsg(c *connect, data []byte) *Message {

	pushMsg := &pb.DownlinkMsg{}
	proto.Unmarshal(data, pushMsg)
	// if pushMsg.MsgID == c.maxMsgID+1 {
	//      c.maxMsgID++
	msg := &Message{}
	json.Unmarshal(pushMsg.DownlinkBody, msg)
	ackMsg := &pb.AckMsg{
		ToType: pb.CMD_Uplink,
		ConnId: c.connID,
	}
	ackData, _ := proto.Marshal(ackMsg)
	c.send(pb.CMD_Ack, ackData)
	return msg
	// }
}

func (c *connect) reConn() {
	c.conn.Close()
	time.Sleep(1 * time.Second)
	addr := &net.TCPAddr{IP: c.ip, Port: c.port}
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		fmt.Printf("DialTCP.err=%+v", err)
	}
	c.conn = conn
}
func (c *connect) send(ty pb.CMD, palyload []byte) error {
	// 直接发送给接收方
	msgCmd := pb.Data{
		Cmd:     ty,
		Payload: palyload,
	}
	msg, err := proto.Marshal(&msgCmd)
	if err != nil {
		panic(err)
	}
	dataPgk := tcp.Packing(msg)
	_, err = c.conn.Write(dataPgk.RetrieveAllString())
	return err
}

func (c *connect) recv() <-chan *Message {
	return c.recvChan
}

func (c *connect) close() {
	// 目前没啥值得回收的
	c.conn.Close()
}
