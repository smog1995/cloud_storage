package mq

import (
	"cloud_storage/config"
	"go.uber.org/zap"


	"github.com/streadway/amqp"
)

ipmort (
	)

func Publish(exchange, routingKey string, msg []byte) bool {
	if !initChannel(config.RabbitURL) {
		zap.S().Error("初始化信道 失败")
	}

	//发布消息没有返回错误
	if nil == channel.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body: msg}) {
		return true
	}
	return false

}