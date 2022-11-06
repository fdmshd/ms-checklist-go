package main

import (
	"checklist/internal/models"
	"database/sql"
	"encoding/json"
	"flag"
	"log"

	_ "github.com/go-sql-driver/mysql"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	dsn := flag.String("dsn", "root:password@tcp(mysql_checklist:3306)/checklist", "MySQL data source name")
	flag.Parse()
	db, err := openDB(*dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	model := models.TaskModel{DB: db}
	exchange := "user_deletion"
	exchangeDLX := "user_deletion_timeout"
	queue := "tasks_deletion_queue"
	tasksKey := "tasks_key"

	conn, err := amqp.Dial("amqp://user:password@rabbitmq_checklist:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(exchange, "direct", true, false, false, false, nil)
	failOnError(err, "Failed to declare Exchange")

	err = ch.ExchangeDeclare(exchangeDLX, "fanout", true, false, false, false, nil)
	failOnError(err, "Failed to declare Exchange")

	_, err = ch.QueueDeclare(queue, false, false, false, false, nil)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(queue, tasksKey, exchange, false, nil)
	failOnError(err, "Failed to bind Queue")
	msgs, err := ch.Consume(queue, "", true, false, false, false, nil)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			var message struct {
				Id int
			}
			err := json.Unmarshal(d.Body, &message)
			if err != nil {
				log.Printf("Error:unmarshal message:%v", err)
			}
			model.DeleteByUser(message.Id)
			if err != nil {
				log.Printf("Error: deleting user %d:%v", message.Id, err)
			}
		}
	}()

	<-forever
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
