package controller

// import (
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"net"
// 	"net/http"
// 	"os"
// 	"time"

// 	"github.com/confluentinc/confluent-kafka-go/kafka"
// )

// func ProducerMessage(ip_address, request string) {
// 	broker := "localhost"
// 	topic := "quickstart-events"

// 	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker})

// 	if err != nil {
// 		fmt.Printf("Failed to create producer: %s\n", err)
// 		os.Exit(1)
// 	}

// 	fmt.Printf("Created Producer %v\n", p)

// 	// Optional delivery channel, if not specified the Producer object's
// 	// .Events channel is used.
// 	deliveryChan := make(chan kafka.Event, 1000)

// 	value := request
// 	err = p.Produce(&kafka.Message{
// 		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
// 		Value:          []byte(value),
// 		Headers:        []kafka.Header{{Key: "ip_address", Value: []byte("header values are binary")}},
// 	}, deliveryChan)

// 	e := <-deliveryChan
// 	m := e.(*kafka.Message)

// 	if m.TopicPartition.Error != nil {
// 		fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
// 		panic(m.TopicPartition.Error)
// 	} else {
// 		fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
// 			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
// 	}

// 	p.Flush(10e3)
// 	close(deliveryChan)
// }

// func createMessage(ip_address, request string) {
// 	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
// 	if err != nil {
// 		panic(err)
// 	}

// 	defer p.Close()

// 	// Delivery report handler for produced messages
// 	go func() {
// 		for e := range p.Events() {
// 			switch ev := e.(type) {
// 			case *kafka.Message:
// 				if ev.TopicPartition.Error != nil {
// 					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
// 				} else {
// 					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
// 				}
// 			}
// 		}
// 	}()

// 	// Produce messages to topic (asynchronously)
// 	topic := "quickstart-events"
// 	ejlog := request
// 	// p.Produce()
// 	p.Produce(&kafka.Message{
// 		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
// 		Value:          []byte(ejlog),
// 		Key:            []byte(ip_address),
// 	}, nil)

// 	// Wait for message deliveries before shutting down
// 	// p.Flush(15 * 1000)
// 	p.Flush(10e3)

// 	p.Close()
// }

// func KafkaMultilineWincor(w http.ResponseWriter, r *http.Request) {
// 	start := time.Now()
// 	requestBody, _ := ioutil.ReadAll(r.Body)
// 	ip_address, _, err := net.SplitHostPort(r.RemoteAddr)
// 	if err != nil {
// 		fmt.Println("Error : ", err)
// 	}
// 	// createMessage(string(ip_address), string(requestBody))
// 	ProducerMessage(string(ip_address), string(requestBody))
// 	elapsed := time.Since(start)
// 	log.Printf("This request took %s", elapsed)
// 	fmt.Fprintf(w, "Berhasil dipanggil dong")
// }
