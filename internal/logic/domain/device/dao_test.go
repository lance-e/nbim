package device_test

import (
	"fmt"
	"nbim/internal/logic/domain/device"
	"testing"
)

func TestDevice_Add(t *testing.T) {
	d := device.Device{
		UserId:        1,
		Type:          1,
		Brand:         "iphone",
		Model:         "ipone13 pm",
		SystemVersion: "8.0.0",
		SDKVersion:    "1.0.0",
		Status:        1,
	}

	err := device.Dao.Save(&d)
	fmt.Println(err)
	fmt.Println(d)
}

func TestDevice_Get(t *testing.T) {
	d, err := device.Dao.Get(1)
	fmt.Printf("%+v\n %+v\n", d, err)
}

func TestDevice_GetOnline(t *testing.T) {
	ds, err := device.Dao.ListAllOnlineDeviceByUserId(1)
	fmt.Println(err)
	fmt.Printf("%+v \n", ds)
}
