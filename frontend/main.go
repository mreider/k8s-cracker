package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	"performance-metrics/rpc"
)

type CrackerID string

type Config struct {
	Port     int `envconfig:"PORT" default:"15002"`
	HTTPPort int `envconfig:"HTTP_PORT" default:"15003"`
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

	cracks := make(chan CrackerID, 100)

	server := grpc.NewServer()
	rpc.RegisterFrontendServiceServer(server, NewService(cracks))
	go func() {
		_ = server.Serve(listener)
		close(cracks)
	}()

	wsConnections := NewWebsocketConnections()

	go func() {
		scores := make(map[CrackerID]int)
		for crackerID := range cracks {
			scores[crackerID]++
			message := fmt.Sprintf(`{"cracker_id": "%s", "score": %d}`, crackerID, scores[crackerID])
			wsConnections.Broadcast(context.Background(), []byte(message))
		}
	}()

	httpService := NewHTTPService(config.HTTPPort, wsConnections)
	httpService.Run()
}
