
package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/hashicorp/consul/api"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
	"transactions/cmd/app"
	"transactions/pkg/transactions"
)

const (
	defaultPort = "9999"
	defaultHost      = "transactions"
	defaultConsulURL = "http://consul:8500"
	defaultIP        = "127.0.0.1"
	defaultDSN  = "postgres://app:pass@transactionsdb:5432/db"
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

	consulURL, ok := os.LookupEnv("APP_CONSUL_URL")
	if !ok {
		consulURL = defaultConsulURL
	}

	ip, ok := os.LookupEnv("APP_IP")
	if !ok {
		ip = defaultIP
	}
	dsn, ok := os.LookupEnv("APP_DSN")
	if !ok {
		dsn = defaultDSN
	}

	if err := execute(host, port, consulURL, ip, dsn); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func execute(host string, port string, consulURL string, ip string, dsn string) error {
	parsedConsulURL, err := url.Parse(consulURL)
	if err != nil {
		return err
	}

	client, err := api.NewClient(&api.Config{
		Address: parsedConsulURL.Host,
		Scheme:  parsedConsulURL.Scheme,
	})
	if err != nil {
		return err
	}

	parsedPort, err := strconv.Atoi(port)
	if err != nil {
		return err
	}

	err = waitForConsulRegistration(client, &api.AgentServiceRegistration{
		ID:      "transactions",
		Name:    "transactions",
		Address: ip,
		Port:    parsedPort,
		Check: &api.AgentServiceCheck{
			Interval:                       "5s",
			Timeout:                        "1s",
			HTTP:                           fmt.Sprintf("http://%s:%s/api/health", host, port),
			Method:                         "GET",
			DeregisterCriticalServiceAfter: "1m",
		},
	})
	if err != nil {
		return err
	}
	ctx := context.Background()
	pool, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		log.Print(err)
		return err
	}

	transactionsSvc := transactions.NewService(pool)
	mux := chi.NewRouter()

	application := app.NewServer(transactionsSvc, mux)
	err = application.Init()
	if err != nil {
		log.Print(err)
		return err
	}

	server := &http.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: application,
	}
	return server.ListenAndServe()
}
func waitForConsulRegistration(client *api.Client, opts *api.AgentServiceRegistration) error {
	for {
		select {
		case <-time.After(time.Minute):
			return errors.New("can't connect to consul")
		default:

		}

		err := client.Agent().ServiceRegister(opts)
		if err != nil {
			log.Print(err)
			time.Sleep(time.Second * 5)
			continue
		}

		return nil
	}
}