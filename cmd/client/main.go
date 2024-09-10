package main

import (
	"bufio"
	"fmt"
	"github.com/siddontang/go-log/log"
	"google.golang.org/protobuf/proto"
	"iot-x/core/codec"
	"iot-x/protobuf"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	decoder          = codec.NewUvarintDecoder()
	encoder          = codec.NewUvarintEncoder(1024)
	token            = ""
	deviceId   int64 = 10002
	toDeviceId int64 = 10003
)

func main() {
	// 服务器地址和端口
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

inputLoop:
	for {
		// 从命令行读取用户输入
		time.Sleep(time.Millisecond * 200)
		fmt.Print("Enter request type (e.g., 1 for login, 2 for report, 3 for SendTo): ")
		reader := bufio.NewReader(os.Stdin)
		requestType, _ := reader.ReadString('\n')
		requestType = strings.TrimSpace(requestType)

		// 根据用户输入发送不同类型的请求
		switch requestType {
		case "1":
			go signIn(conn)
			goto inputLoop
		case "2":
			go reportInfo(conn)
			goto inputLoop
		case "3":
			go sendTo(conn)
			goto inputLoop
		default:
			fmt.Println("Unsupported request type.")
		}
	}

}

func signIn(conn net.Conn) {
	device, _ := proto.Marshal(&protobuf.DeviceInfo{
		Sn:       "BX10000010",
		Model:    "X30",
		DeviceId: deviceId,
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

func sendTo(conn net.Conn) {
	i := 0
	for {
		device := &protobuf.SendDevice{ToDeviceId: toDeviceId, FromDeviceId: deviceId, Type: 1, Data: []byte(fmt.Sprintf("%d", i) + "来自设备:" + strconv.FormatInt(deviceId, 10))}

		bytes, _ := proto.Marshal(device)

		msg := &protobuf.MessageInput{
			RequestId: "fdskfkdsfksk333333333",
			MagicNum:  0x11,
			Version:   protobuf.MessageInput_VERSION_1,
			Type:      protobuf.MessageInput_SEND_TO_DEVICE,
			Token:     token,
			Data:      bytes,
		}
		sendRequest(conn, msg)
		time.Sleep(time.Millisecond * 10)
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
