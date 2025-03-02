package sdk

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"nbim/pkg/protocol/pb"
	"net"
	"sync"
	"time"

	"github.com/goccy/go-json"
	"google.golang.org/protobuf/proto"
)

const (
	MsgTypeText      = "text"
	MsgTypeAck       = "ack"
	MsgTypeReConn    = "reConn"
	MsgTypeHeartbeat = "heartbeat"
	MsgLogin         = "loginMsg"
)

type Chat struct {
	Nick             string
	UserID           string
	SessionID        int64
	conn             *connect
	closeChan        chan struct{}
	MsgClientIDTable map[int64]int64
	sync.RWMutex
}

type Message struct {
	Type       string
	Name       string
	FormUserID string
	ToUserID   string
	Content    string
	Session    int64
}

func NewChat(ip net.IP, port int, nick, userID string, sessionID int64) *Chat {
	chat := &Chat{
		Nick:             nick,
		UserID:           userID,
		SessionID:        sessionID,
		conn:             newConnect(ip, port),
		closeChan:        make(chan struct{}),
		MsgClientIDTable: make(map[int64]int64),
	}
	go chat.loop()      //处理可读事件(接收消息)
	chat.login()        //登陆
	go chat.heartbeat() //周期性发送心跳
	return chat
}
func (chat *Chat) Send(msg *Message) {
	data, _ := json.Marshal(msg)
	//key := fmt.Sprintf("%d", chat.conn.connID)
	upMsg := &pb.UplinkMsg{
		ClientId: chat.getClientID(chat.SessionID),
		// ConnId:    chat.conn.connID,
		SessionId: chat.SessionID,

		UplinkBody: data,
	}
	palyload, _ := proto.Marshal(upMsg)
	chat.conn.send(pb.CMD_Uplink, palyload)
}

func (chat *Chat) GetCurClientID() int64 {
	if id, ok := chat.MsgClientIDTable[chat.SessionID]; ok {
		return id
	}
	return 0
}

// Close close chat
func (chat *Chat) Close() {
	chat.conn.close()
	close(chat.closeChan)
	close(chat.conn.recvChan)
	close(chat.conn.sendChan)
}

func (chat *Chat) ReConn() {
	chat.Lock()
	defer chat.Unlock()
	// 需要重置clientId
	chat.MsgClientIDTable = make(map[int64]int64)
	chat.conn.reConn() //重新启动一个tcp连接
	time.Sleep(5 * time.Second)
	chat.reConn()
}

// Recv receive message
func (chat *Chat) Recv() <-chan *Message {
	return chat.conn.recv()
}

func (chat *Chat) loop() {
Loop:
	for {
		select {
		case <-chat.closeChan:
			return
		default:
			mc := &pb.Data{}
			data, err := ReadData(chat.conn.conn)
			if err != nil {
				goto Loop
			}
			err = proto.Unmarshal(data, mc)
			if err != nil {
				panic(err)
			}
			var msg *Message
			switch mc.Cmd {
			case pb.CMD_Ack:
				msg = handAckMsg(chat.conn, mc.Payload)
			case pb.CMD_Downlink:
				msg = handPushMsg(chat.conn, mc.Payload)

			}
			chat.conn.recvChan <- msg
		}
	}
}

func (chat *Chat) getClientID(sessionID int64) int64 {
	chat.Lock()
	defer chat.Unlock()
	var res int64
	if id, ok := chat.MsgClientIDTable[sessionID]; ok {
		res = id
	}
	chat.MsgClientIDTable[sessionID] = res + 1
	return res
}

func (chat *Chat) login() {
	loginMsg := pb.LoginMsg{
		DeviceId: 123,
	}
	palyload, err := proto.Marshal(&loginMsg)
	if err != nil {
		panic(err)
	}
	chat.conn.send(pb.CMD_Login, palyload)
}

func (chat *Chat) reConn() {
	reConn := pb.ReconnMsg{
		ConnId: chat.conn.connID,
	}
	palyload, err := proto.Marshal(&reConn)
	if err != nil {
		panic(err)
	}
	chat.conn.send(pb.CMD_Reconn, palyload)
}

func (chat *Chat) heartbeat() {
	tc := time.NewTicker(1 * time.Second)
	defer func() {
		chat.heartbeat()
	}()
loop:
	for {
		select {
		case <-chat.closeChan:
			return
		case <-tc.C:
			hearbeat := pb.HeartbeatMsg{}
			palyload, err := proto.Marshal(&hearbeat)
			if err != nil {
				panic(err)
			}
			err = chat.conn.send(pb.CMD_Heartbeat, palyload)
			if err != nil {
				goto loop
			}
		}
	}
}

func ReadData(conn *net.TCPConn) ([]byte, error) {
	var dataLen uint32
	dataLenBuf := make([]byte, 4)
	if err := readFixedData(conn, dataLenBuf); err != nil {
		return nil, err
	}
	// fmt.Printf("readFixedData:%+v\n", dataLenBuf)
	buffer := bytes.NewBuffer(dataLenBuf)
	if err := binary.Read(buffer, binary.BigEndian, &dataLen); err != nil {
		return nil, fmt.Errorf("read headlen error:%s", err.Error())
	}
	if dataLen <= 0 {
		return nil, fmt.Errorf("wrong headlen :%d", dataLen)
	}
	dataBuf := make([]byte, dataLen)
	// fmt.Printf("readFixedData.dataLen:%+v\n", dataLen)
	if err := readFixedData(conn, dataBuf); err != nil {
		return nil, fmt.Errorf("read headlen error:%s", err.Error())
	}
	return dataBuf, nil
}

// 读取固定buf长度的数据
func readFixedData(conn *net.TCPConn, buf []byte) error {
	_ = (*conn).SetReadDeadline(time.Now().Add(time.Duration(120) * time.Second))
	var pos int = 0
	var totalSize int = len(buf)
	for {
		c, err := (*conn).Read(buf[pos:])
		if err != nil {
			return err
		}
		pos = pos + c
		if pos == totalSize {
			break
		}
	}
	return nil
}
