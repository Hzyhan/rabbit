package rabbit

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

// PublishPub 订阅模式生产
func (r *MQRabbit) PublishPub(message string) {
	//1.尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.ExchangeName, // 交换机名称，若无则创建
		"fanout",       // 交换机类型
		true,           // 持久化
		false,          // 自动删除
		false,          //true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,          // 无需确认
		nil,
	)

	failOnError(err, "Failed to declare an exchange")

	//2.发送消息
	err = r.channel.Publish(
		r.ExchangeName,
		r.Key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

// ReceiveSub 订阅模式消费端代码
func (r *MQRabbit) ReceiveSub() {
	//1.试探性创建交换机
	err := r.channel.ExchangeDeclare(
		r.ExchangeName, // 交换机名称，若无则创建
		"fanout",       // 交换机类型
		true,           // 持久化
		false,          // 自动删除
		false,          //true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,          // 无需确认
		nil,
	)
	failOnError(err, "Failed to declare an exchange")
	//2.试探性创建队列，这里注意队列名称不要写
	q, err := r.channel.QueueDeclare(
		"", //随机生产队列名称
		false,
		false,
		true,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	//绑定队列到 exchange 中
	err = r.channel.QueueBind(
		q.Name,
		"", // 在pub/sub模式下，这里的key要为空
		r.ExchangeName,
		false,
		nil)

	//消费消息
	messges, err := r.channel.Consume(
		q.Name,
		"",    // 消费者名称
		true,  // 自动确认
		false, // 唯一消费者
		false,
		false, // 无需确认
		nil,
	)

	forever := make(chan bool)

	go func() {
		for d := range messges {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Println("退出请按 CTRL+C\n")
	<-forever
}
