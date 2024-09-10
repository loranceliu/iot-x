package connect

import (
	"iot-x/types"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

const (
	DeviceKey    = "device:"
	DeviceExpire = 2 * time.Hour
)

type DeviceContext struct{}

// Get 获取指定用户的所有在线设备
func (c *DeviceContext) GetDevice(deviceId int64) (*types.Device, error) {
	device := new(types.Device)
	err := IotS.Redis.GetObj(DeviceKey+strconv.FormatInt(deviceId, 10), device)
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if err == redis.Nil {
		return nil, nil
	}
	return device, nil
}

// Set 将指定用户的所有在线设备存入缓存
func (c *DeviceContext) SetDevice(deviceId int64, device types.Device) error {
	err := IotS.Redis.SetObj(DeviceKey+strconv.FormatInt(deviceId, 10), device, DeviceExpire)
	return err
}

// Del 删除用户的在线设备列表
func (c *DeviceContext) DelDevice(deviceId int64) error {
	key := DeviceKey + strconv.FormatInt(deviceId, 10)
	err := IotS.Redis.Del(key)
	return err
}
