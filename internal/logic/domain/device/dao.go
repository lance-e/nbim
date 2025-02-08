package device

import (
	"errors"
	"nbim/pkg/db"
	"nbim/pkg/gerror"
	"time"

	"github.com/jinzhu/gorm"
)

type dao struct{}

var Dao = new(dao)

// Save:保存device所有字段(UPDATE),该id不存在就创建(INSERT)
func (*dao) Save(device *Device) error {
	device.CreateTime = time.Now()
	device.UpdateTime = time.Now()
	err := db.DB.Save(&device).Error
	if err != nil {
		return gerror.WrapError(err)
	}
	return nil
}

// Get:获取对于deviceId的设备信息
func (*dao) Get(id int64) (*Device, error) {
	var device = Device{Id: id}
	err := db.DB.First(&device).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, gerror.WrapError(err)
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &device, nil
}

// ListAllOnlineDeviceByUserId:查询用户所有的在线设备
func (*dao) ListAllOnlineDeviceByUserId(userid int64) ([]Device, error) {
	var devices []Device
	err := db.DB.Find(&devices, "user_id = ? and status = ?", userid, DeviceOnline).Error
	if err != nil {
		return nil, gerror.WrapError(err)
	}
	return devices, nil
}

// ListAllOnlineDeviceByServerAddr:查询用户所有的在线设备
func (*dao) ListAllOnlineDeviceByServerAddr(serverAddr string) ([]Device, error) {
	var devices []Device
	err := db.DB.Find(&devices, "server_addr = ? and status = ?", serverAddr, DeviceOnline).Error
	if err != nil {
		return nil, gerror.WrapError(err)
	}
	return devices, nil
}

// UpdateStatus:更新设备在线状态
func (*dao) UpdateStatus(deviceId int64, serverAddr string, status int) (int64, error) {
	db_ := db.DB.Model(&Device{}).Where("id = ? and server_addr = ?", deviceId, serverAddr).Update("status", status)
	if db_.Error != nil {
		return 0, gerror.WrapError(db_.Error)
	}
	return db_.RowsAffected, nil
}
