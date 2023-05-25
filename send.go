package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/Kibuns/UserService/Models"
	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func send(msgName string, user *Models.User) {
    conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
    failOnError(err, "Failed to connect to RabbitMQ")
    defer conn.Close()

    ch, err := conn.Channel()
    failOnError(err, "Failed to open a channel")
    defer ch.Close()

    q, err := ch.QueueDeclare(
        msgName, // name
        false,   // durable
        false,   // delete when unused
        false,   // exclusive
        false,   // no-wait
        nil,     // arguments
    )
    failOnError(err, "Failed to declare a queue")

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    body, err := json.Marshal(user)
    if err != nil {
        log.Panicf("Failed to marshal Twoot: %s", err)
    }

    err = ch.PublishWithContext(ctx,
        "",     // exchange
        q.Name, // routing key
        false,  // mandatory
        false,  // immediate
        amqp.Publishing{
            ContentType: "application/json",
            Body:        body,
        })
    failOnError(err, "Failed to publish a message")

    log.Printf(" [x] Sent %s\n", body)
}