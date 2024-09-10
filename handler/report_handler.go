package handler

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	connect "iot-x/conn"
	"iot-x/protobuf"
)

type ReportHandler struct {
}

func (h *ReportHandler) Handle(conn *connect.Conn, message *protobuf.MessageInput) {

	out := &protobuf.MessageOutput{
		RequestId: connect.IotS.UUID.GenerateID(),
		Type:      protobuf.MessageOutput_ACK,
	}

	aa, _ := proto.Marshal(out)

	err := conn.Write(aa)
	if err != nil {
		return
	}
	fmt.Printf("Handler: %v\n", message)
}

func (h *ReportHandler) Type() protobuf.MessageInput_Type {
	return protobuf.MessageInput_REPORT
}
