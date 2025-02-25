package gateway

import (
	"sync"

	"github.com/gorilla/websocket"
	netreactors "github.com/lance-e/net-reactors"
)

const (
	ConnTypeTCP int8 = 1
	ConnTypeWS  int8 = 2
)

var (
	TcpConnToConnID = sync.Map{}
	IDtoConnInfo    = sync.Map{}
)

type ConnInfo struct {
	ConnID   int64
	ConnType int8
	TCP      *netreactors.TcpConnection
	WS       *websocket.Conn
	WSMutex  sync.Mutex
}
