package mq

import (
	"cloud_storage/config"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"log"
)

var conn *amqp.Connection
var channel *amqp.Channel

// 如果异常关闭，会接收通知
var notifyClose chan *amqp.Error

// UpdateRabbitHost : 更新mq host
func UpdateRabbitHost(host string) {
	config.RabbitURL = host
}

// Init 初始化
func Init() {
	// 是否开启异步转移功能，开启时才初始化rabbitMQ连接
	if !config.AsyncTransferEnable {
		return
	}
	if initChannel(config.RabbitURL) {
		channel.NotifyClose(notifyClose)
	}
	// 断线自动重连
	go func() {
		for {
			select {
			case msg := <-notifyClose:
				conn = nil
				channel = nil
				log.Printf("onNotifyChannelClosed: %+v\n", msg)
				initChannel(config.RabbitURL)
			}
		}
	}()
}

func initChannel(rabbitHost string) bool {
	if channel != nil {
		return true
	}

	conn, err := amqp.Dial(rabbitHost)
	if err != nil {
		zap.S().Error(err)
		return false
	}

	channel, err = conn.Channel()
	if err != nil {
		zap.S().Error(err)
		return false
	}
	return true
}
