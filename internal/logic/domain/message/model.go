package message

import "time"

// 消息总表
type Messages struct {
	Seq         int64
	SenderId    int64
	SessionId   uint64
	Content     []byte
	SendTime    int64 //时间戳
	MessageType string
	IsDeleted   int
	CreateTime  time.Time
	UpdateTime  time.Time
}

// 用户消息表(索引扩散而非数据扩散):用户收件箱,所有用户共享一个表,查完seq后去Messages总表中查
type UserMessages struct {
	UserId      int64     //消息接收者id
	Seq         int64     //消息序列号
	ReceiveTime int64     //接收时间
	Status      int       //状态, 0:未读,1:已读
	IsDeleted   int       //是否删除,0:未删除,1:已删除
	CreateTime  time.Time //创建时间
	UpdateTime  time.Time //更新时间
}

// 所有的用户群组消息状态
type UserGroupMessageStatus struct {
	UserId            int64
	GroupId           int64
	LastReadMessageId int64
	LastReadTime      time.Time
	CreateTime        time.Time
	UpdateTime        time.Time
}
