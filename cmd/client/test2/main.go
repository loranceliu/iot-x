package main

import (
	"fmt"
	"github.com/siddontang/go-log/log"
	"google.golang.org/protobuf/proto"
	"iot-x/core/codec"
	"iot-x/protobuf"
	"net"
	"strconv"
	"time"
)

var (
	decoder          = codec.NewUvarintDecoder()
	encoder          = codec.NewUvarintEncoder(1024)
	token            = ""
	deviceId   int64 = 10005
	toDeviceId int64 = 10003
)

func main() {
	// 服务器地址和端口

	for i := 10005; i <= 20000; i++ {
		serverAddr := "192.168.2.201:8081"

		// 连接到服务器
		conn, err := net.Dial("tcp", serverAddr)
		if err != nil {
			fmt.Println("Error connecting to the server:", err)
			return
		}
		defer conn.Close()

		fmt.Println("Connected to the server.")

		go handleConn(conn)

		go signIn(int64(i), conn)

		go sendTo(int64(i), conn)
	}

	select {}

}

func signIn(a int64, conn net.Conn) {
	device, _ := proto.Marshal(&protobuf.DeviceInfo{
		Sn:       "BX10000010",
		Model:    "X30",
		DeviceId: a,
		Token:    "aaa.bbb.ccc",
	})

	msg := &protobuf.MessageInput{
		RequestId: "fdskfkdsfksk11111111",
		MagicNum:  0x11,
		Version:   protobuf.MessageInput_VERSION_1,
		Type:      protobuf.MessageInput_LOGIN,
		Data:      device,
	}
	sendRequest(conn, msg)
}

func reportInfo(conn net.Conn) {
	msg := &protobuf.MessageInput{
		RequestId: "fdskfkdsfksk222222222",
		MagicNum:  0x11,
		Version:   protobuf.MessageInput_VERSION_1,
		Type:      protobuf.MessageInput_REPORT,
		Token:     token,
		Data:      []byte("上报信息"),
	}
	sendRequest(conn, msg)
}

func sendTo(a int64, conn net.Conn) {
	i := 0
	for {
		device := &protobuf.SendDevice{ToDeviceId: toDeviceId, FromDeviceId: a, Type: 1, Data: []byte(fmt.Sprintf("%d", i) + "来自设备:" + strconv.FormatInt(a, 10))}

		bytes, _ := proto.Marshal(device)

		msg := &protobuf.MessageInput{
			RequestId: "fdskfkdsfksk333333333",
			MagicNum:  0x11,
			Version:   protobuf.MessageInput_VERSION_1,
			Type:      protobuf.MessageInput_SEND_TO_DEVICE,
			Token:     "aaa.bbb.ccc",
			Data:      bytes,
		}
		sendRequest(conn, msg)
		time.Sleep(time.Millisecond * 1000)
		i++
	}
}

func sendRequest(conn net.Conn, msg *protobuf.MessageInput) {
	bytes, _ := proto.Marshal(msg)

	err := encoder.EncodeToWriter(conn, bytes)
	if err != nil {
		log.Errorf("err:", err.Error())
		return
	}
}

func handleConn(conn net.Conn) {
	buffer := codec.NewBuffer(make([]byte, 1024))
	var handler = func(bytes []byte) {
		msg := &protobuf.MessageOutput{}
		proto.Unmarshal(bytes, msg)

		if msg.Type == protobuf.MessageOutput_AUTH_FAILED {
			log.Error(msg.Msg)
		}

		if msg.Type == protobuf.MessageOutput_AUTH_SUCCESS {
			deviceInfo := &protobuf.DeviceInfo{}
			proto.Unmarshal(msg.Data, deviceInfo)
			token = deviceInfo.Token
			log.Info(deviceInfo)
		}

		if msg.Type == protobuf.MessageOutput_ACK {
			log.Info("message ack.")
		}

		if msg.Type == protobuf.MessageOutput_SEND_FROM_DEVICE {
			sendDevice := &protobuf.SendDevice{}
			proto.Unmarshal(msg.Data, sendDevice)
			log.Info(string(sendDevice.Data))
		}

	}

	for {
		_, err := buffer.ReadFromReader(conn)
		if err != nil {
			log.Errorf("error", err)
			return
		}

		err = decoder.Decode(buffer, handler)
		if err != nil {
			log.Error(err)
			return
		}
	}
}
