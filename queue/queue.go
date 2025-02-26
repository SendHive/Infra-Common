package queue

import (
	"context"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type IQueueService interface {
	Connect() (conn *amqp.Connection, err error)
	DeclareQueue(conn *amqp.Connection) (qu amqp.Queue, err error)
	PublishMessage(qu amqp.Queue, conn *amqp.Connection) error
	ConsumeMessage(qu amqp.Queue, conn *amqp.Connection, isTest bool) error
}

type QueueService struct{}

func NewQueueRequest() (IQueueService, error) {
	return &QueueService{}, nil
}

func (q *QueueService) Connect() (conn *amqp.Connection, err error) {
	// Connect to rabbitmq
	conn, err = amqp.Dial("amqp://user:password@localhost:5672/")
	if err != nil {
		log.Println("error while connecting to RabbitMQ: ", err)
		panic(err)
	}
	fmt.Println("Connected to RabbitMQ Sucessfully!")
	return conn, err
}

func (q *QueueService) DeclareQueue(conn *amqp.Connection) (qu amqp.Queue, err error) {
	ch, err := conn.Channel()
	if err != nil {
		log.Println("error while creating a channel: ", err)
		return amqp.Queue{}, err
	}
	queue, err := ch.QueueDeclare(
		"task-queue",
		false,
		false,
		false,
		false,
		nil,
	)
	defer ch.Close()
	fmt.Println("Declared Queue Sucessfully!")
	return queue, nil
}

func (q *QueueService) PublishMessage(qu amqp.Queue, conn *amqp.Connection) error {
	ch, err := conn.Channel()
	if err != nil {
		log.Println("error while creating a channel: ", err)
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	body := "Hello World"
	err = ch.PublishWithContext(ctx,
		"",
		qu.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	if err != nil {
		log.Println("error while publishing the queue text: ", err)
		return err
	}
	fmt.Println("Published Message Sucessfully!")
	return nil
}

func (q *QueueService) ConsumeMessage(qu amqp.Queue, conn *amqp.Connection, isTest bool) error {
	ch, err := conn.Channel()
	if err != nil {
		log.Println("error while creating a channel: ", err)
		return err
	}
	defer ch.Close()
	msgs, err := ch.Consume(
		qu.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println("error while consuming the msgs : ", err)
		return err
	}

	var forever chan struct{}

	if isTest {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			break
		}
		fmt.Println("Consumed Message Sucessfully!")
	} else {
		go func() {
			for d := range msgs {
				log.Printf("Received a message: %s", d.Body)
				if len(d.Body) == 1 {
					break
				}
			}
		}()
		<-forever
	}
	return nil
}
