package filter

import (
	"google.golang.org/protobuf/proto"
	connect "iot-x/conn"
	"iot-x/protobuf"
)

type AuthFilter struct {
}

func (auth *AuthFilter) PreFilter(conn *connect.Conn, input *protobuf.MessageInput) bool {
	if input.Type == protobuf.MessageInput_LOGIN {
		return true
	}

	if input.Token != "aaa.bbb.ccc" {
		msg := &protobuf.MessageOutput{
			RequestId: connect.IotS.UUID.GenerateID(),
			Type:      protobuf.MessageOutput_AUTH_FAILED,
			Msg:       "client auth failed.",
		}
		bytes, _ := proto.Marshal(msg)
		conn.Write(bytes)
		return false
	}

	return true
}

func (auth *AuthFilter) Order() uint8 {
	return 0
}
