package main

import (
	"fmt"
	"net"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	"performance-metrics/rpc"
)

type Config struct {
	Port int `envconfig:"PORT" default:"15001"`
}

func main() {
	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "2006-01-02T15:04:05.000Z07:00"}).With().Timestamp().Logger()

	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal().Msgf("can't process environment variables: %v", err)
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Port))
	if err != nil {
		log.Fatal().Msgf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	rpc.RegisterLockboxServiceServer(server, NewService())
	_ = server.Serve(listener)
}
