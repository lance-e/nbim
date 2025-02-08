package device

import (
	"nbim/pkg/protocol/pb"
	"time"
)

const (
	DeviceOnline  = 1 //设备在线
	DeviceOffline = 0 //设备离线
)

// 设备信息
type Device struct {
	Id            int64     //设备id
	UserId        int64     //用户id
	Type          int64     //设备类型,1.Android;2.IOS;3.Windows;4.Macos;5:Web
	Brand         string    //手机厂商
	Model         string    //机型
	SystemVersion string    //系统版本
	SDKVersion    string    //SDK版本
	Status        int32     //在线状态,0:不在线;1:在线
	ServerAddr    string    //服务器地址
	ClientAddr    string    //客户端地址
	CreateTime    time.Time //创建时间
	UpdateTime    time.Time //更新时间
}

func (d *Device) ToProtoDevice() *pb.Device {
	return &pb.Device{
		DeviceId:      d.Id,
		UserId:        d.UserId,
		Type:          d.Type,
		Brand:         d.Brand,
		Model:         d.Model,
		SystemVersion: d.SystemVersion,
		SdkVersion:    d.SDKVersion,
		Status:        d.Status,
		ServerAddr:    d.ServerAddr,
		ClientAddr:    d.ClientAddr,
		CreateTime:    d.CreateTime.UnixMilli(),
		UpdateTime:    d.UpdateTime.UnixMilli(),
	}
}

// 判断设备信息是否合法
func (d *Device) IsLegal() bool {
	if d.Type == 0 || d.Brand == "" || d.Model == "" || d.SystemVersion == "" || d.SDKVersion == "" {
		return false
	}
	return true
}

// 设备上线
func (d *Device) Online(userId int64, serverAddr string, clientAdd string) {
	d.UserId = userId
	d.ServerAddr = serverAddr
	d.ClientAddr = clientAdd
	d.Status = DeviceOnline
}

// 设备离线
func (d *Device) Offline(userId int64, serverAddr string, clientAdd string) {
	d.UserId = userId
	d.ServerAddr = serverAddr
	d.ClientAddr = clientAdd
	d.Status = DeviceOffline
}
