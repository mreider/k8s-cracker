package main

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	defaultMemLeakChunkSizeMB = 10
)

func startMemLeakIfNeed(ctx context.Context, memLeak string) {
	if memLeak == "" {
		return
	}

	memLeak = strings.ToUpper(memLeak)
	memLeakChunkSizeMB, err := strconv.Atoi(memLeak)
	if err != nil {
		if memLeak != "TRUE" {
			return
		}
		memLeakChunkSizeMB = defaultMemLeakChunkSizeMB
	}

	log.Info().Msgf("Memory Leak: %d MB/sec", memLeakChunkSizeMB)
	go doMemoryLeak(ctx, memLeakChunkSizeMB*1024*1024, time.Second*1)
}

func doMemoryLeak(ctx context.Context, leakSize int, interval time.Duration) {
	var chunks [][]byte
	chunk := make([]byte, leakSize)
	for i := range chunk {
		chunk[i] = byte(i)
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			clone := make([]byte, len(chunk))
			copy(clone, chunk)
			chunks = append(chunks, clone)
		}
	}
}
