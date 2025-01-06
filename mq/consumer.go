package mq

import "go.uber.org/zap"

var done chan bool

func StartConsume(qName, cName string, callback func(msg []byte) bool) {
	msgs, err := channel.Consume(
		qName,
		cName,
		true,  //自动应答
		false, // 非唯一的消费者
		false, // rabbitMQ只能设置为false
		false, // noWait, false表示会阻塞直到有消息过来
		nil)
	if err != nil {
		zap.S().Fatalf(err.Error())
	}
	done = make(chan bool)

	go func() {
		for d := range msgs {
			processErr := callback(d.Body)
			if processErr != nil {
				// TODO:将错误写入消息队列，待后续处理
			}

		}
	}()

	<-done
	channel.Close()

}

// StopConsume : 停止监听队列
func StopConsume() {
	done <- true
}
