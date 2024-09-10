package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/siddontang/go-log/log"
	connect "iot-x/conn"
	"iot-x/filter"
	"iot-x/handler"
	"iot-x/types"
	"iot-x/utils"
)

func main() {

	conf := utils.InitConf()

	iotServer, err := connect.NewServer(&types.Instance{
		Node: conf.Server.Node,
		IP:   conf.Server.IP,
		Port: conf.Server.Port,
	}, &handler.LoginHandler{}, &handler.ReportHandler{}, &handler.SendToDeviceHandler{})

	if err != nil {
		log.Fatalf("iotServer running failed: %v", err)
	}

	iotServer.UseFilter(&filter.AuthFilter{})

	iotServer.UseUUID(utils.NewSnowflake("ZPI7DQROM123"))

	iotServer.UseRedis(redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", conf.Redis.Host, conf.Redis.Port),
		Password: "",
		DB:       conf.Redis.DB,
	}))

	iotServer.UseDeviceContext(new(connect.DeviceContext))

	iotServer.UseMQ()

	iotServer.StartServer()
}
