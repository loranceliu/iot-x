package connect

import (
	"iot-x/protobuf"
)

type Filter interface {
	PreFilter(conn *Conn, input *protobuf.MessageInput) bool

	Order() uint8
}
