// John Jovi Sidabutar - 1304212134
package main

import (
	"encoding/json" // untuk mengencoding dan decoding JSON
	"log"           // untuk logging
	"time"          // untuk memanipulasi waktu

	amqp "github.com/rabbitmq/amqp091-go" // untuk berinteraksi dengan RabbitMQ
)

type SensorData struct {
	SensorID    string    `json:"sensor_id"`
	Temperature string    `json:"temperature"`
	Humidity    string    `json:"humidity"`
	AirPressure string    `json:"air_pressure"`
	Timestamp   time.Time `json:"timestamp"`
}

// error handling
func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	// membuat koneksi ke rabbitmq
	// dengan guest sebagai username & password
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// membuat channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// mendeclare queue untuk menyimpan pesan yg akan diproses
	q, err := ch.QueueDeclare(
		"sensor_data", // name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// consume digunakan untuk menerima pesan dari queue
	// yang telah dideklarasikan
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool) // membuat go channel untuk menunggu signal

	// memulai goroutine (mini thread di go) untuk memproses pesan yg diterima
	// pada anonymous function ini, dilakukan proses unmarshall/decoding isi
	// pesan JSON ke dalam struct SensorData
	go func() {
		for d := range msgs {
			var sensorData SensorData
			err := json.Unmarshal(d.Body, &sensorData)
			failOnError(err, "Failed to unmarshall JSON")

			log.Printf("[*] Data sent: %+v\n", sensorData)
		}
	}()

	log.Printf("[*] Waiting for logs. To exit press CTRL+C\n")
	<-forever // menunggu signal dari channel forever (akan menunggu selamanya)
}
