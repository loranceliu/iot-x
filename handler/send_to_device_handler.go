package handler

import (
	"fmt"
	"github.com/siddontang/go-log/log"
	"google.golang.org/protobuf/proto"
	connect "iot-x/conn"
	"iot-x/protobuf"
)

type SendToDeviceHandler struct {
}

func (h *SendToDeviceHandler) Handle(conn *connect.Conn, message *protobuf.MessageInput) {
	sendFromDevice := &protobuf.SendDevice{}
	proto.Unmarshal(message.Data, sendFromDevice)

	out := &protobuf.MessageOutput{
		RequestId: connect.IotS.UUID.GenerateID(),
		Type:      protobuf.MessageOutput_SEND_FROM_DEVICE,
		Data:      message.Data,
	}

	o, _ := proto.Marshal(out)

	device, _ := connect.IotS.DeviceContext.GetDevice(sendFromDevice.GetToDeviceId())

	if device == nil {
		log.Info("device is not online")
		return
	}

	if device.Instance.Node != connect.IotS.Instance.Node {
		t, _ := proto.Marshal(message)
		conn.SendMsg(device.Instance.Node, t)
	} else {
		if conn := connect.GetConn(sendFromDevice.ToDeviceId); conn != nil {
			err := conn.Write(o)
			if err != nil {
				return
			}
		}
	}

	fmt.Printf("Handler: %v\n", message)
}

func (h *SendToDeviceHandler) Type() protobuf.MessageInput_Type {
	return protobuf.MessageInput_SEND_TO_DEVICE
}
