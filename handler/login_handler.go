package handler

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	connect "iot-x/conn"
	"iot-x/protobuf"
	"iot-x/types"
)

type LoginHandler struct {
}

func (h *LoginHandler) Handle(conn *connect.Conn, message *protobuf.MessageInput) {

	deviceInfo := &protobuf.DeviceInfo{}

	proto.Unmarshal(message.Data, deviceInfo)

	conn.DeviceId = deviceInfo.DeviceId

	connect.SetConn(conn.DeviceId, conn)

	storeDeviceCache(deviceInfo)

	out := &protobuf.MessageOutput{
		RequestId: connect.IotS.UUID.GenerateID(),
		Type:      protobuf.MessageOutput_AUTH_SUCCESS,
		Data:      message.Data,
	}

	aa, _ := proto.Marshal(out)

	err := conn.Write(aa)
	if err != nil {
		return
	}
	fmt.Printf("Handler: %v\n", message)
}

func (h *LoginHandler) Type() protobuf.MessageInput_Type {
	return protobuf.MessageInput_LOGIN
}

func storeDeviceCache(deviceInfo *protobuf.DeviceInfo) {
	device := types.Device{
		DeviceId: deviceInfo.DeviceId,
		Sn:       deviceInfo.Sn,
		Model:    deviceInfo.Model,
		Instance: *connect.IotS.Instance,
	}
	connect.IotS.DeviceContext.SetDevice(deviceInfo.DeviceId, device)
}
