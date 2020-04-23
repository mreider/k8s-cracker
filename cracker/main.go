package main

import (
	"context"
	"os"
	"runtime"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	"performance-metrics/rpc"
	"performance-metrics/utils"
)

type Config struct {
	LockboxEndpoint  string `envconfig:"LOCKBOX_ENDPOINT" default:":15001"`
	FrontendEndpoint string `envconfig:"FRONTEND_ENDPOINT" default:":15002"`
	Workers          int    `envconfig:"WORKERS" default:"0"`
	MemLeak          string `envconfig:"MEMLEAK"`
	CPULeak          string `envconfig:"CPULEAK"`
}

func main() {
	ctx := context.Background()
	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "2006-01-02T15:04:05.000Z07:00"}).With().Timestamp().Logger()

	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal().Msgf("can't process environment variables: %v", err)
	}

	crackerID, err := utils.GenerateRandomString(10)
	if err != nil {
		log.Fatal().Msgf("can't generate cracker id: %v", err)
	}

	workerCount := config.Workers
	if workerCount <= 0 {
		workerCount = runtime.NumCPU()
	}

	startMemLeakIfNeed(ctx, config.MemLeak)
	startCPULeakIfNeed(config.CPULeak)

	lockboxConn, err := grpc.Dial(config.LockboxEndpoint, grpc.WithInsecure())
	if err != nil {
		log.Fatal().Msgf("can't connect to lockbox: %v", err)
	}
	defer func() { _ = lockboxConn.Close() }()

	frontendConn, err := grpc.Dial(config.FrontendEndpoint, grpc.WithInsecure())
	if err != nil {
		log.Fatal().Msgf("can't connect to frontend: %v", err)
	}
	defer func() { _ = frontendConn.Close() }()

	lockboxClient := rpc.NewLockboxServiceClient(lockboxConn)
	frontendClient := rpc.NewFrontendServiceClient(frontendConn)

	notifications := make(chan struct{}, 100)
	defer close(notifications)
	go sendNotifications(ctx, frontendClient, crackerID, notifications)

	cracker, err := NewCracker(workerCount, func() (string, error) {
		resp, err := lockboxClient.GeneratePassword(ctx, &rpc.GeneratePasswordRequest{})
		if err != nil {
			return "", err
		}
		return resp.Password, nil
	}, func() {
		notifications <- struct{}{}
	})
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	done := cracker.Start(ctx)
	err = <-done
	if err != nil {
		log.Fatal().Err(err).Send()
	}
}

func sendNotifications(ctx context.Context, frontendClient rpc.FrontendServiceClient, crackerID string, notifications <-chan struct{}) {
	for range notifications {
		_, err := frontendClient.NotifyCracked(ctx, &rpc.NotifyCrackedRequest{
			CrackerId: crackerID,
		})
		if err != nil {
			log.Error().Err(err).Send()
		}
	}
}
