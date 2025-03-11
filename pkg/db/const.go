package db

import "time"

const (
	//最新的客户端消息序列号:1->多,一个用户有多个session
	MaxClientIDKey = "max_client_id_{%d}_%d_%d" //max_client_id_{'slot'}_'connID'_'sessionID'
	//最后发送的消息:多->1,多个session都会通过这一个接口发送消息给同一个用户
	LastMessageKey = "last_message_{%d}_%d" //last_message_{'slot'}_'connID'
	//设备id与连接id的映射
	DeviceIdToConnId = "device_conn_%d" //device_conn_'deviceID'

	//群组的下一条消息序列号
	GroupSeqIDKey = "group_seqid_%d" //group_seqid_'group_id'

	TTL7Day = 7 * 24 * time.Hour //一周过期时间
)
