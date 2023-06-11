package clientChat

import (
	"context"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go" // Делаем удобное имя для импорта в нашем коде
)

func RegisterQueue() (*amqp.Channel, *amqp.Queue) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/") // Создаем подключение к RabbitMQ
	if err != nil {
		log.Fatalf("unable to open connect to RabbitMQ server. Error: %s", err)
	}

	defer func() {
		_ = conn.Close() // Закрываем подключение в случае удачной попытки
	}()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open channel. Error: %s", err)
	}

	defer func() {
		_ = ch.Close() // Закрываем канал в случае удачной попытки открытия
	}()

	q, err := ch.QueueDeclare(
		"hello2", // name
		false,    // durable
		true,     // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		log.Fatalf("failed to declare a queue. Error: %s", err)
	}
	return ch, &q
}

func main() {
	ch, q := RegisterQueue()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for i := 0; i < 3; i++ {
		var body string
		fmt.Scanf("%s\n", &body)
		//body := "Hello World!" + fmt.Sprint(i)
		err := ch.PublishWithContext(ctx,
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		if err != nil {
			log.Fatalf("failed to publish a message. Error: %s", err)
		}
		log.Printf(" [x] Sent %s\n", body)
	}

}
