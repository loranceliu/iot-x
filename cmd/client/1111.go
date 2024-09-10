package main

import (
	"google.golang.org/protobuf/proto"
	"iot-x/core/codec"
	"iot-x/protobuf"
	"log"
	"net"
	"time"
)

var (
// decoder = codec.NewUvarintDecoder()
// encoder = codec.NewUvarintEncoder(1024)
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Println("error dialing", err.Error())
		return
	}

	msg := &protobuf.MessageInput{
		MagicNum: 0x11,
		Version:  protobuf.MessageInput_VERSION_1,
		Type:     protobuf.MessageInput_UNKNOWN_TYPE,
		Token:    "aaaa",
	}
	bytes, _ := proto.Marshal(msg)

	go func() {
		i := 0
		for {
			err := encoder.EncodeToWriter(conn, bytes)
			time.Sleep(time.Second * 1)
			if err != nil {
				log.Println("err:", err.Error())
				return
			}
			i++
		}
	}()

	go handleConn(conn)

	select {}
}

func whandleConn(conn net.Conn) {
	buffer := codec.NewBuffer(make([]byte, 1024))
	var handler = func(bytes []byte) {
		msg := &protobuf.MessageOutput{}
		proto.Unmarshal(bytes, msg)
		log.Println(msg)
	}

	for {
		_, err := buffer.ReadFromReader(conn)
		if err != nil {
			log.Println("error", err)
			return
		}

		err = decoder.Decode(buffer, handler)
		if err != nil {
			log.Println(err)
			return
		}
	}
}
