package main

import (
	"context"
	"errors"
	"github.com/Shopify/sarama"
	"google.golang.org/grpc"
	"log"
	"microservice/services/payments-api/cmd/app"
	"microservice/services/payments-api/pkg/payments"
	"net"
	serverPb "microservice/services/payments-api/pkg/server"
	"os"
	"time"
)

const (
	defaultPort = "9999"
	defaultHost = "0.0.0.0"
	defaultBrokerURL = "kafka:9092"
	defaultTopic = "payments"
)

func main() {
	port, ok := os.LookupEnv("APP_PORT")
	if !ok {
		port = defaultPort
	}

	host, ok := os.LookupEnv("APP_HOST")
	if !ok {
		host = defaultHost
	}

	brokerURL, ok := os.LookupEnv("APP_BROKER_URL")
	if !ok {
		brokerURL = defaultBrokerURL
	}

	topic, ok := os.LookupEnv("APP_TOPIC")
	if !ok {
		topic = defaultTopic
	}

	if err := execute(net.JoinHostPort(host, port), brokerURL, topic); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func execute(addr string, brokerURL string, topic string) error {
	log.Println("Start execution of Payment Api")
	producer, err := waitForKafka(brokerURL)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := producer.Close(); cerr != nil {
			if err == nil {
				err = cerr
				return
			}
			log.Print(cerr)
		}
	}()

	paymentsSvc := payments.NewService(producer, topic)
	listener, err := net.Listen("tcp", addr)
	if err != nil{
		return nil
	}
	ctx := context.Background()
	grpcServer := grpc.NewServer()
	server := app.NewServer(paymentsSvc, ctx)
	serverPb.RegisterPaymentsServerServer(grpcServer, server)
	return grpcServer.Serve(listener)
}

func waitForKafka(brokerURL string) (sarama.SyncProducer, error) {
	for {
		select {
		case <- time.After(time.Minute):
			return nil, errors.New("can't connect to kafka")
		default:

		}
		sarama.Logger = log.New(os.Stdout, "", log.Ltime)
		config := sarama.NewConfig()
		config.ClientID = "payments-api"
		config.Producer.Return.Successes = true
		config.Version = sarama.V2_6_0_0

		producer, err := sarama.NewSyncProducer([]string{brokerURL}, config)
		if err != nil {
			log.Print(err)
			time.Sleep(time.Second)
			continue
		}
		return producer, nil
	}
}