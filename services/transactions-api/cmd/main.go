package main

import (
	"context"
	"contrib.go.opencensus.io/exporter/jaeger"
	"go.opencensus.io/trace"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"transactions-api/cmd/app"
	serverPb "transactions-api/pkg/server"
	"transactions-api/pkg/transactions"
)

const (
	defaultPort = "9999"
	defaultHost = "0.0.0.0"
	defaultTransactionsURL = "http://transactions:9999"
)
func InitJaeger(serviceName string) error{
	exporter, err := jaeger.NewExporter(jaeger.Options{
		AgentEndpoint: "localhost:6831",
		Process: jaeger.Process{
			ServiceName: serviceName,
			Tags: []jaeger.Tag{
				jaeger.StringTag("hostname", "localhost"),
			},
		},
	})
	if err != nil {
		return err
	}
	trace.RegisterExporter(exporter)
	trace.ApplyConfig(trace.Config{
		DefaultSampler: trace.AlwaysSample(),
	})
	return nil
}

func main() {
	port, ok := os.LookupEnv("APP_PORT")
	if !ok {
		port = defaultPort
	}

	host, ok := os.LookupEnv("APP_HOST")
	if !ok {
		host = defaultHost
	}

	transactionsURL, ok := os.LookupEnv("APP_TRANSACTIONS_URL")
	if !ok {
		transactionsURL = defaultTransactionsURL
	}
	err := InitJaeger("transactions")
	if err != nil{
		log.Println(err)
		os.Exit(1)
	}
	if err := execute(net.JoinHostPort(host, port), transactionsURL); err != nil {
		os.Exit(1)
	}
}

func execute(addr string, transactionsURL string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil{
		return nil
	}
	transactionsSvc := transactions.NewService(&http.Client{}, transactionsURL)
	ctx := context.Background()
	grpcServer := grpc.NewServer()
	server := app.NewServer(transactionsSvc, ctx)
	serverPb.RegisterTransactionsServerServer(grpcServer, server)
	return grpcServer.Serve(listener)
}