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
	PublishMessage(qu amqp.Queue, conn *amqp.Connection, body string) error
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

func (q *QueueService) PublishMessage(qu amqp.Queue, conn *amqp.Connection, body string) error {
	ch, err := conn.Channel()
	if err != nil {
		log.Println("error while creating a channel: ", err)
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = ch.PublishWithContext(ctx,
		"",
		qu.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
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
		log.Println("Error while creating a channel:", err)
		return err
	}
	defer ch.Close()

	// Ensure prefetch is set to avoid overwhelming the consumer
	err = ch.Qos(1, 0, false)
	if err != nil {
		log.Println("Error setting QoS:", err)
		return err
	}

	msgs, err := ch.Consume(
		qu.Name,
		"",
		false, // Auto-Ack set to false to manually acknowledge after processing
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println("Error while consuming the msgs:", err)
		return err
	}

	if isTest {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			_ = d.Ack(false) // Acknowledge message after processing
			break
		}
		fmt.Println("Consumed Message Successfully!")
	} else {
		for d := range msgs {
			go func(msg amqp.Delivery) {
				log.Printf("Processing message: %s", msg.Body)
				err := msg.Ack(false) // Acknowledge message after processing
				if err != nil {
					log.Println("Failed to acknowledge message:", err)
				}
			}(d)
		}
	}

	return nil
}
