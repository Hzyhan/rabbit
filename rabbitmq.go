package rabbit

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
	"strconv"
	"sync"
)

// MQRabbit rabbitmq对象
type MQRabbit struct {
	conn     *amqp.Connection
	connOnce sync.Once
	channel  *amqp.Channel
	//队列名称
	QueueName string
	//交换机
	ExchangeName string
	//key
	Key string
}

// GetMqURL 获取mq链接
func GetMqURL() string {
	name := os.Getenv("RABBITMQ_NAME")
	pwd := os.Getenv("RABBITMQ_PASSWORD")
	ip := os.Getenv("RABBITMQ_IP")
	vhost := os.Getenv("RABBITMQ_VHOST")
	port, _ := strconv.Atoi(os.Getenv("RABBITMQ_PORT"))
	return getMqURL(name, pwd, ip, vhost, port)
}

func getMqURL(name, pwd, ip, vhost string, port int) string {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", name, pwd, ip, port, vhost)
	return url
}

// failOnError 错误处理函数
func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

// NewRabbitMQ 创建rabbitmq实例
func NewRabbitMQ(queueName, exchange, key, mqURL string) *MQRabbit {
	rabbitmq := &MQRabbit{QueueName: queueName, ExchangeName: exchange, Key: key}
	var err error
	//创建rabbitmq连接
	rabbitmq.conn, err = amqp.Dial(mqURL)
	failOnError(err, "创建连接错误")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	failOnError(err, "获取channel失败")
	return rabbitmq
}
