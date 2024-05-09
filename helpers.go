package main

import (
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"mime"
	"net/http"
	"time"
)

func CheckHTTPRequest(w http.ResponseWriter, req *http.Request) {
	// Enforce a JSON Content-Type.
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}
}

func JSONWrapper(obj interface{}) []byte {
	js, err := json.Marshal(obj)
	if err != nil {
		return nil
	}
	return js
}

func RabbitPrefix(conn *amqp.Connection, channel *amqp.Channel, queueName string) (*amqp.Connection, *amqp.Channel, amqp.Queue, error) {
	var err error
	if conn == nil {
		conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
		if err != nil {
			log.Fatalf("%s: %s", "RabbitMQ failed", err)
			return nil, nil, amqp.Queue{}, err
		}
	}

	if channel == nil {
		channel, err = conn.Channel()
		if err != nil {
			log.Fatalf("%s: %s", "Failed to open a channel", err)
			return nil, nil, amqp.Queue{}, err
		}
	}

	queue, err := channel.QueueDeclare(
		queueName, // Name of the queue
		false,     // Durable
		false,     // Delete when unused
		false,     // Exclusive
		false,     // No-wait
		nil,       // Arguments
	)

	return conn, channel, queue, nil
}

func sendToRabbit(obj interface{}, queueName string, conn *amqp.Connection, channel *amqp.Channel) {
	conn_ready := true
	channel_ready := true
	if conn == nil {
		conn_ready = false
	}
	if channel == nil {
		channel_ready = false
	}
	conn, channel, queue, err := RabbitPrefix(conn, channel, queueName)
	if !conn_ready {
		defer conn.Close()
	}
	if !channel_ready {
		defer channel.Close()
	}

	if err != nil {
		log.Fatalf("%s: %s", "Failed to preset RabbitMQ", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = channel.PublishWithContext(ctx,
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        JSONWrapper(obj),
		})
}

func getFromRabbit(queueName string, conn *amqp.Connection, channel *amqp.Channel) ([]byte, error) {
	conn_ready := true
	channel_ready := true
	if conn == nil {
		conn_ready = false
	}
	if channel == nil {
		channel_ready = false
	}
	conn, channel, queue, err := RabbitPrefix(conn, channel, queueName)
	if !conn_ready {
		//defer conn.Close()
	}
	if !channel_ready {
		//defer channel.Close()
	}

	if err != nil {
		log.Fatalf("%s: %s", "Failed to preset RabbitMQ", err)
		return []byte{}, nil
	}

	if err := channel.Qos(1, 0, false); err != nil {
		log.Fatal("Qos Setting was unsuccessful")
		return []byte{}, nil
	}

	messages, err := channel.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	msg := <-messages

	if !conn_ready {
		conn.Close()
	}
	if !channel_ready {
		channel.Close()
	}
	return msg.Body, nil
}
