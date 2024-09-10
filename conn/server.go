package connect

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/siddontang/go-log/log"
	"github.com/streadway/amqp"
	"google.golang.org/protobuf/proto"
	gn "iot-x/core"
	"iot-x/core/codec"
	"iot-x/protobuf"
	"iot-x/types"
	"iot-x/utils"
	"sort"
	"time"
)

var IotS *IotServer

type Handler struct{}

func (*Handler) OnConnect(c *gn.Conn) {
	con := &Conn{
		Conn: c,
		ch:   IotS.ch,
	}
	c.SetData(con)
}
func (*Handler) OnMessage(c *gn.Conn, bytes []byte) {
	conn := c.GetData().(*Conn)
	conn.HandleMessage(bytes)
}
func (*Handler) OnClose(c *gn.Conn, err error) {
	conn := c.GetData().(*Conn)
	DeleteConn(conn.DeviceId)
	IotS.DeviceContext.DelDevice(conn.DeviceId)
	log.Infof("fd:%v close: %v deviceId: %v", c.GetFd(), err, conn.DeviceId)
}

type IotServer struct {
	Instance      *types.Instance
	handlers      []EventHandler
	filters       []Filter
	UUID          *utils.Snowflake
	Redis         *utils.RedisUtil
	DeviceContext *DeviceContext
	ch            *amqp.Channel
}

var exchange string

var queue string

var routing string

var consumer string

func NewServer(instance *types.Instance, eventHandler ...EventHandler) (*IotServer, error) {
	if len(eventHandler) == 0 {
		return nil, errors.New("handler can't be null")
	}
	IotS = &IotServer{Instance: instance, handlers: eventHandler}
	return IotS, nil
}

func (is *IotServer) StartServer() {
	server, err := gn.NewServer(fmt.Sprintf(":%v", is.Instance.Port), &Handler{},
		gn.WithDecoder(codec.NewUvarintDecoder()),
		gn.WithEncoder(codec.NewUvarintEncoder(1024)),
		gn.WithTimeout(10*time.Minute),
		gn.WithReadBufferLen(256),
		gn.WithAcceptGNum(10),
		gn.WithIOGNum(10))
	if err != nil {
		log.Info("err")
		return
	}
	go is.StartConsume()
	server.Run()
}

func (is *IotServer) StartConsume() {
	msgs, err := is.ch.Consume(
		queue,    // 队列名称
		consumer, // 消费者
		true,     // 自动应答
		false,    // 独占
		false,    // 不等待
		false,    // 参数
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	for d := range msgs {
		msg := &protobuf.MessageInput{}
		proto.Unmarshal(d.Body, msg)
		for _, h := range IotS.getHandler() {
			if h.Type() == msg.Type {
				h.Handle(nil, msg)
			}
		}
		fmt.Printf("Received a message: %s\n", d.Body)
	}

}

func (is *IotServer) StartProducer() {
	_, err := is.ch.QueueDeclare(
		queue, // 队列名称
		false, // 持久性
		false, // 自动删除
		false, // 独占
		false, // 不等待
		nil,   // 参数
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}
}

func (is *IotServer) UseMQ() {
	conn, err := amqp.Dial("amqp://test:test@192.168.2.204:5672/vhost")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	//defer conn.Close()

	exchange = fmt.Sprintf("%s_Exchange", is.Instance.Node)
	queue = fmt.Sprintf("%s_Queue", is.Instance.Node)
	routing = fmt.Sprintf("%s_Routing", is.Instance.Node)
	consumer = fmt.Sprintf("%s_Consumer", is.Instance.Node)

	// 创建一个通道
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	//defer ch.Close()

	// 声明一个Exchange
	err = ch.ExchangeDeclare(
		exchange, // 交换器名称
		"direct", // 交换器类型
		true,     // 持久性
		false,    // 自动删除
		false,    // 不等待
		false,    // 参数
		nil,
	)

	if err != nil {
		log.Fatalf("Failed to declare an exchange: %v", err)
	}

	// 声明一个Queue
	q, err := ch.QueueDeclare(
		queue, // 队列名称（留空表示由RabbitMQ自动生成）
		false, // 持久性
		false, // 自动删除
		true,  // 独占
		false, // 不等待
		nil,   // 参数
	)

	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	// 将Queue绑定到Exchange
	err = ch.QueueBind(
		q.Name,   // 队列名称
		routing,  // Routing key
		exchange, // 交换器名称
		false,    // 不等待
		nil,      // 参数
	)
	if err != nil {
		log.Fatalf("Failed to bind a queue: %v", err)
	}
	is.ch = ch
}

func (is *IotServer) UseFilter(filter ...Filter) {
	sort.Slice(filter, func(i, j int) bool {
		return filter[i].Order() < filter[j].Order()
	})

	is.filters = filter
}

func (is *IotServer) UseUUID(snowflake *utils.Snowflake) {
	is.UUID = snowflake
}

func (is *IotServer) UseRedis(client *redis.Client) {
	is.Redis = utils.NewRedisUtil(client)
}

func (is *IotServer) UseDeviceContext(ctx *DeviceContext) {
	is.DeviceContext = ctx
}

func (is *IotServer) getHandler() []EventHandler {
	return is.handlers
}

func (is *IotServer) getFilter() []Filter {
	return is.filters
}

func (is *IotServer) GetInstanceInfo() *types.Instance {
	return is.Instance
}

func (is *IotServer) GetMqCh() *amqp.Channel {
	return is.ch
}
