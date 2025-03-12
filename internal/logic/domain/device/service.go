package device

import (
	"context"
	"fmt"
	"nbim/pkg/gerror"
	"nbim/pkg/logger"
	"nbim/pkg/protocol/pb"
	"nbim/pkg/rpc"
	"time"

	"go.uber.org/zap"
)

type service struct{}

var Service = new(service)

func (*service) Register(ctx context.Context, req *pb.RegisterDeviceReq) (int64, error) {
	device := Device{
		Type:          req.Type,
		Brand:         req.Brand,
		Model:         req.Model,
		SystemVersion: req.SystemVersion,
		SDKVersion:    req.SdkVersion,
	}
	//判断设备信息是否合法
	if !device.IsLegal() {
		return 0, gerror.ErrBadRequest
	}

	//mysql保存设备信息
	err := Dao.Save(&device)
	if err != nil {
		return 0, err
	}
	//redis清除缓存
	if device.UserId != 0 {
		err = Cathe.Del(device.UserId)
		if err != nil {
			return 0, err
		}
	}

	return device.Id, nil
}

func (*service) ConnSignIn(ctx context.Context, req *pb.ConnSignInReq) error {
	_, err := rpc.GetLogicIntClient().Auth(ctx, &pb.AuthReq{
		UserId:   req.UserId,
		DeviceId: req.DeviceId,
		Token:    req.Token,
	})
	if err != nil {
		return err
	}

	//标记用户在设备上登陆(该设备在线)
	device, err := Dao.Get(req.DeviceId)
	if err != nil {
		return err
	}
	if device == nil {
		return nil
	}
	device.Online(req.UserId, req.ServerAddr, req.ClientAddr)
	err = Dao.Save(device)
	if err != nil {
		return err
	}
	return nil
}

func (*service) Offline(ctx context.Context, req *pb.OfflineReq) error {
	device, err := Dao.Get(req.DeviceId)
	if err != nil {
		return err
	}
	if device == nil {
		return nil
	}
	if device.ClientAddr != req.ClientAddr {
		return nil
	}

	//下线
	device.Status = DeviceOffline

	err = Dao.Save(device)
	if err != nil {
		return err
	}
	return nil
}

func (*service) GetDevice(ctx context.Context, req *pb.GetDeviceReq) (*pb.Device, error) {
	device, err := Dao.Get(req.DeviceId)
	if err != nil {
		return nil, err
	}
	if device == nil {
		return nil, gerror.ErrDeviceNotExist
	}
	return device.ToProtoDevice(), nil
}

// ServerStop:长连接层服务停止，需要将连接在这台长连接层服务器上的所有设备下线
func (*service) ServerStop(ctx context.Context, req *pb.ServerStopReq) error {
	devices, err := Dao.ListAllOnlineDeviceByServerAddr(req.ServerAddr)
	if err != nil {
		return err
	}
	for i := range devices {
		//异步修改设备状态
		if n, err := Dao.UpdateStatus(devices[i].Id, req.ServerAddr, DeviceOffline); err == nil && n == 1 && devices[i].UserId != 0 {
			err = Cathe.Del(devices[i].UserId)
		}
		if err != nil {
			logger.Logger.Error("device.server.serverstop error", zap.Any("device", devices[i]), zap.Error(err))
		}
		time.Sleep(2 * time.Millisecond)
	}
	return nil
}

func (*service) ListAllOnlineDeviceByUserId(ctx context.Context, userid int64) ([]*pb.Device, error) {
	//先从缓存里面找
	devices, err := Cathe.Get(userid)
	if err != nil {
		return nil, err
	}
	if devices != nil {
		fmt.Printf("%d devices\n", len(devices))
		pbDevices := make([]*pb.Device, len(devices))
		for i := range devices {
			pbDevices[i] = devices[i].ToProtoDevice()
		}
		return pbDevices, nil
	}
	//没找到去数据库里面找
	devices, err = Dao.ListAllOnlineDeviceByUserId(userid)
	if err != nil {
		return nil, err
	}

	//更新缓存
	err = Cathe.Set(userid, devices)
	if err != nil {
		return nil, err
	}

	pbDevices := make([]*pb.Device, len(devices))
	for i := range devices {
		pbDevices[i] = devices[i].ToProtoDevice()
	}
	return pbDevices, nil
}
