package message

import (
	"context"
	"nbim/internal/logic/proxy"
	"nbim/pkg/id"
	"nbim/pkg/protocol/pb"
	"nbim/pkg/rpc"
	"time"
)

type service struct{}

var Service = new(service)
var SnowflakeId *id.Snowflake //用于tcp连接的id生成

func init() {
	var err error
	SnowflakeId, err = id.NewSnowflake(0) //TODO:use nodeId
	if err != nil {
		//FIXME
		panic(err)
	}
}

// 接收到上行消息，先生成唯一消息序列号ID并落库,再异步发送下行消息
func (*service) SendToGroup(senderId int64, deviceId int64, sessionId uint64, content []byte, sendTime int64) error {
	groupId := int64(sessionId & 0x7fffffffffffffff)

	//step1:获取seqID
	id, err := Seq.GetGroupMsgSeq(groupId)
	if err != nil {
		return err
	}

	//step2:落库(TODO:异步落库)
	msgs := []*Messages{
		{
			Seq:         id,
			SenderId:    senderId,
			SessionId:   sessionId,
			Content:     content,
			SendTime:    sendTime,
			MessageType: "text", //TODO:目前仅支持文本
			IsDeleted:   0,
			CreateTime:  time.Now(),
			UpdateTime:  time.Now(),
		},
	}
	err = Dao.Save(msgs[0])
	if err != nil {
		return err
	}

	//step3:分发(TODO:异步分发)
	ctx := context.TODO()
	//--先获取全部群成员
	members, err := proxy.GroupProxy.GetGroupAllMember(ctx, groupId)
	if err != nil {
		return err
	}
	//--再获取每个群成员登陆过的设备
	for _, member := range members {
		devices, err := proxy.DevcieProxy.ListAllOnlineDeviceByUserId(ctx, member.UserId)
		if err != nil {
			return err
		}
		for _, device := range devices {
			if device.DeviceId == deviceId {
				continue
			}
			Service.SendToDevice(device, msgs)
		}
	}
	return nil
}

func (*service) SendToUser(senderId int64, deviceId int64, userId int64, content []byte, sendTime int64) error {
	//step1:获取seqID
	id := SnowflakeId.Generate()
	//step2:落库(TODO:异步落库)
	msgs := []*Messages{
		{
			Seq:         id,
			SenderId:    senderId,
			SessionId:   uint64(userId),
			Content:     content,
			SendTime:    sendTime,
			MessageType: "text", //TODO:目前仅支持文本
			IsDeleted:   0,
			CreateTime:  time.Now(),
			UpdateTime:  time.Now(),
		},
	}

	//--先写入总表
	err := Dao.Save(msgs[0])
	if err != nil {
		return err
	}
	//--再写入用户写信箱(写扩散)
	err = Dao.SaveUserMsg(msgs)
	if err != nil {
		return err
	}
	//step3:分发(TODO:异步分发)

	devices, err := proxy.DevcieProxy.ListAllOnlineDeviceByUserId(context.TODO(), userId)
	if err != nil {
		return err
	}
	for _, device := range devices {
		if device.DeviceId == deviceId {
			continue
		}
		Service.SendToDevice(device, msgs)
	}
	return nil
}

func (*service) SendToDevice(device *pb.Device, msgs []*Messages) error {
	downs := make([]*pb.DownlinkMsg, len(msgs))
	for _, msg := range msgs {
		downs = append(downs, &pb.DownlinkMsg{
			Seq:          msg.Seq,
			SenderId:     msg.SenderId,
			SessionId:    msg.SessionId,
			DownlinkBody: msg.Content,
		})
	}
	_, err := rpc.GetStateClient().DelieverDownlinkMessage(context.TODO(), &pb.DelieverDownlinkMessageReq{
		DeviceId: device.DeviceId,
		Downs:    downs,
	})
	if err != nil {
		return err
	}
	return nil
}
