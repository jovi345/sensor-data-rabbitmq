// John Jovi Sidabutar - 1304212134

package main

import (
	"encoding/json" // untuk mengencoding dan decoding JSON
	"log"           // untuk logging
	"math/rand"     // untuk generate random number
	"strconv"       // untuk konversi string
	"time"          // untuk memanipulasi waktu

	"github.com/google/uuid"              // untuk generate UUID
	amqp "github.com/rabbitmq/amqp091-go" // untuk berinteraksi dengan RabbitMQ
)

// mendefinisikan tipe bentukan data
type SensorData struct {
	SensorID    string    `json:"sensor_id"`
	Temperature string    `json:"temperature"`
	Humidity    string    `json:"humidity"`
	AirPressure string    `json:"air_pressure"`
	Timestamp   time.Time `json:"timestamp"`
}

// generate sensor data dummy
func generateSensorData() SensorData {
	// membuat id sensor unik dengan awalan SD- dan UUID
	id := "SD-" + uuid.New().String()[:8]
	// generate nilai suhu, kelembapan, dan tekanan udara
	temp := (20 + rand.Float64()) * (1.5)
	humidity := (30 + rand.Float64()) * (1.6)
	air_pressure := (40 + rand.Float64()) * (1.7)

	// mereturn objek SensorData dengan nilai-nilai yang dihasilkan
	return SensorData{
		SensorID:    id,
		Temperature: strconv.FormatFloat(temp, 'f', 2, 64) + " Celcius",
		Humidity:    strconv.FormatFloat(humidity, 'f', 2, 64) + "%",
		AirPressure: strconv.FormatFloat(air_pressure, 'f', 2, 64) + "%",
		Timestamp:   time.Now().Local(),
	}
}

// error handling
func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	// membuat koneksi ke rabbitmq
	// dengan username & password: guest
	// defer digunakan sebagai pernyataan
	// yang digunakan untuk menunda eksekusi
	// fungsi sampai fungsi yang memanggilnya selesai
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// membuat channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// membuat antrian/queue
	q, err := ch.QueueDeclare(
		"sensor_data", // name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	failOnError(err, "Failed to declare a queue")

	for {
		sensorData := generateSensorData()    // generate data sensor baru
		body, err := json.Marshal(sensorData) // konversi data sensor ke format JSON
		failOnError(err, "Failed to marshal JSON")

		ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate

			amqp.Publishing{
				ContentType: "application/json", // tipe konten data yg dikirim
				Body:        body,               // isi pesan yang dikirm
			})
		failOnError(err, "Failed to publish a message")

		log.Printf("[*] Data sent: %s\n", body)

		time.Sleep(1 * time.Minute) // memastikan agar data dikirim tiap 1 menit
	}

}
