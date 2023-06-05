package rabbitmq

import (
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/okayes/gopkgs/logger"
	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	url           string
	conn          *amqp.Connection
	ch            *amqp.Channel
	ReconnectFlag int32
}

func NewRabbitMQ(url string) *RabbitMQ {
	rabbitMQ := &RabbitMQ{
		url: url,
	}
	rabbitMQ.connect(false)
	return rabbitMQ
}

func (th *RabbitMQ) connect(isReconnect bool) {
	for {
		if isReconnect {
			time.Sleep(time.Second)
		}

		conn, err := amqp.Dial(th.url)
		if err != nil {
			em := fmt.Sprintf("RabbitMQ dial error: %s, url: %s", err, th.url)
			logger.ErrorMsg(em)
			time.Sleep(time.Second)
			continue
		}

		ch, err := conn.Channel()
		if err != nil {
			em := fmt.Sprintf("RabbitMQ open channel error: %s", err)
			logger.ErrorMsg(em)
			time.Sleep(time.Second)
			continue
		}

		err = ch.Confirm(false)
		if err != nil {
			em := fmt.Sprintf("RabbitMQ set channel confirm error: %s", err)
			logger.ErrorMsg(em)
			time.Sleep(time.Second)
			continue
		}

		if conn == nil || conn.IsClosed() {
			time.Sleep(time.Second)
			continue
		} else {
			atomic.StoreInt32(&th.ReconnectFlag, 0)
			// todo: atomic
			th.conn = conn
			th.ch = ch

			em := fmt.Sprintf("RabbitMQ conn finish: %s", time.Now().Format("2006-01-02 15:04:05"))
			logger.ErrorMsg(em)
			return
		}
	}
}

func (th *RabbitMQ) Publish(exchange string, routingKey string, message []byte, expiration string) error {
	err := th.ch.Publish(exchange, routingKey, false, false,
		amqp.Publishing{
			Body:         message,
			DeliveryMode: 2,
			Expiration:   expiration,
		})

	// reconnect
	if err != nil && th.conn.IsClosed() {
		if atomic.CompareAndSwapInt32(&th.ReconnectFlag, 0, 1) {
			em := fmt.Sprintf("RabbitMQ conn is closed, error: %s", err)
			logger.ErrorMsg(em)
			go th.connect(true)
		}
	}

	return err
}

func (th *RabbitMQ) Subscribe(queue string, prefetchCount int, handler func([]byte) bool) {
	err := th.ch.Qos(prefetchCount, 0, false)
	if err != nil {
		em := fmt.Sprintf("RabbitMQ channel set Qos error: %s", err)
		logger.ErrorMsg(em)
		return
	}

	messages, err := th.ch.Consume(queue, "", false, false, false, false, nil)
	if err != nil {
		em := fmt.Sprintf("RabbitMQ consume message error: %s", err)
		logger.ErrorMsg(em)
		return
	}

	log.Println("consumer waiting to receive message from queue: " + queue)
	for item := range messages {
		if handler(item.Body) {
			err := th.ch.Ack(item.DeliveryTag, false)
			if err != nil {
				em := fmt.Sprintf("RabbitMQ channel ack error: %s, data: %s", err, item.Body)
				logger.ErrorMsg(em)
			}
		} else {
			err := th.ch.Nack(item.DeliveryTag, false, true)
			if err != nil {
				em := fmt.Sprintf("RabbitMQ channel nack error: %s, data: %s", err, item.Body)
				logger.ErrorMsg(em)
			}
		}
	}
	em := fmt.Sprintf("RabbitMQ conn is closed, consumer exit")
	logger.ErrorMsg(em)

	// reconnect
	if th.conn.IsClosed() {
		em := fmt.Sprintf("RabbitMQ conn is closed, begin reconnect")
		logger.ErrorMsg(em)
		th.connect(true)
		go th.Subscribe(queue, prefetchCount, handler)
	}
}

func (th *RabbitMQ) Close() {
	if th.conn == nil {
		return
	}

	err := th.conn.Close()
	if err != nil {
		em := fmt.Sprintf("RabbitMQ close conn error: %s", err)
		logger.ErrorMsg(em)
	}
}
