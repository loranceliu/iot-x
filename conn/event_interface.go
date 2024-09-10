package connect

import (
	"iot-x/protobuf"
)

type EventHandler interface {
	Handle(conn *Conn, message *protobuf.MessageInput)

	Type() protobuf.MessageInput_Type
}
