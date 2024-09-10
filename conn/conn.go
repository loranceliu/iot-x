package connect

import (
	"fmt"
	log "github.com/siddontang/go-log/log"
	"github.com/streadway/amqp"
	"google.golang.org/protobuf/proto"
	gn "iot-x/core"
	"iot-x/protobuf"
)

type Conn struct {
	CoonType int8     // 连接类型
	Conn     *gn.Conn // tcp连接
	UserId   int64    // 用户ID
	DeviceId int64    // 设备ID
	ch       *amqp.Channel
}

func (c *Conn) Write(bytes []byte) error {
	return c.Conn.WriteWithEncoder(bytes)
}

func (c *Conn) SendMsg(node string, bytes []byte) {
	// 发送消息到队列
	err := c.ch.Publish(
		fmt.Sprintf("%s_Exchange", node), // 交换器
		fmt.Sprintf("%s_Routing", node),  // 队列名称
		false,                            // 强制
		false,                            // 立即
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        bytes,
		},
	)
	if err != nil {
		log.Errorf("send message error: %v", err)
	}
}

func (c *Conn) Close() error {

	c.Conn.Close()

	return nil
}

func (c *Conn) GetAddr() string {
	return c.Conn.GetAddr()
}

// HandleMessage 消息处理
func (c *Conn) HandleMessage(bytes []byte) {
	msg := &protobuf.MessageInput{}
	err := proto.Unmarshal(bytes, msg)

	if err != nil {
		log.Errorf("unmarshal error: %v", err)
		return
	}

	for _, f := range IotS.getFilter() {
		if !f.PreFilter(c, msg) {
			return
		}
	}

	for _, h := range IotS.getHandler() {
		if h.Type() == msg.Type {
			h.Handle(c, msg)
		}
	}
}
