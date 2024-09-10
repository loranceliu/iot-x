package connect

import (
	"sync"
)

var ConnContext = sync.Map{}

// SetConn 存储
func SetConn(deviceId int64, conn *Conn) {
	ConnContext.Store(deviceId, conn)
}

// GetConn 获取
func GetConn(deviceId int64) *Conn {
	value, ok := ConnContext.Load(deviceId)
	if ok {
		return value.(*Conn)
	}
	return nil
}

// DeleteConn 删除
func DeleteConn(deviceId int64) {
	ConnContext.Delete(deviceId)
}
