package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/streadway/amqp"

	"github.com/k8-proxy/k8-go-comm/pkg/rabbitmq"
)

var (
	AdpatationReuquestExchange   = "adaptation-request-exchange"
	AdpatationReuquestRoutingKey = "adaptation-request"
	AdpatationReuquestQueueName  = "adaptation-request"

	AdaptationOutcomeExchange   = "adaptation-outcome-exchange"
	AdaptationOutcomeRoutingKey = "adaptation-outcome"
	AdaptationOutcomeQueueName  = "adaptation-outcome"
)

var (
	conn      *amqp.Connection
	mqtimeout error = errors.New("consumer queue timeout")
)

func main() {

	var err error

	MqHost := os.Getenv("MQHOST")
	MqPort := os.Getenv("MQPORT")

	MqUser := os.Getenv("MQUSER")
	MqPass := os.Getenv("MQPASS")

	if len(os.Args) < 2 {
		log.Println("not enough args")
		return
	}

	fn := os.Args[1]
	fullpath := fmt.Sprintf("/tmp/%s", fn)
	fnrebuild := fmt.Sprintf("/tmp/rebuild-%s", fn)
	//path := path.Dir(fullpath)

	conn, err = rabbitmq.NewInstance(MqHost, MqPort, MqUser, MqPass)
	if err != nil {
		log.Fatal("rabbitmq server not found\n", err)
	}

	ctable := amqp.Table{}

	// Start a consumer
	msgs, ch, err := rabbitmq.NewQueueConsumer(conn, AdaptationOutcomeQueueName, AdaptationOutcomeExchange, amqp.ExchangeDirect, AdaptationOutcomeRoutingKey, ctable)
	if err != nil {
		log.Println(err)

	}

	defer ch.Close()

	table := amqp.Table{
		"file-id":               fn,
		"source-file-location":  fullpath,
		"rebuilt-file-location": fnrebuild,
	}

	publisher, err := rabbitmq.NewQueuePublisher(conn, AdpatationReuquestExchange, amqp.ExchangeDirect)
	if err != nil {
		log.Println(err)
		return
	}

	defer publisher.Close()

	err = rabbitmq.PublishMessage(publisher, AdpatationReuquestExchange, AdpatationReuquestRoutingKey, table, []byte{})
	if err != nil {
		log.Println("PublishMessage", err)

		return
	}

	var urlf, urlr string

	const duration = 30 * time.Second
	timer := time.NewTimer(duration)
	brk := make(chan bool, 1)

	go func() {
		for d := range msgs {
			ok := false

			urlf, ok = d.Headers["rebuilt-file-location"].(string)
			if !ok {
				log.Println("error converting rebuilt-file-location to string")

			}

			urlr, ok = d.Headers["report-presigned-url"].(string)
			if !ok {
				log.Println("error converting clean-presigned-url to string")

			}

			brk <- true

		}

	}()
	select {
	case <-timer.C:
		log.Println(mqtimeout)
	case <-brk:
		log.Println(("file processed"))
	}

	log.Println(urlf, urlr)

}
