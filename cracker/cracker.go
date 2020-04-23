package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
)

const (
	passwordLetters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-+"
	iterationCount  = 1e6
)

type Cracker struct {
	workerCount         int
	getOriginalPassword func() (string, error)
	passwordCracked     func()
	letterIndexes       map[byte]byte
}

func NewCracker(workerCount int, getOriginalPassword func() (string, error), passwordCracked func()) (*Cracker, error) {
	if len(passwordLetters)%workerCount != 0 {
		return nil, fmt.Errorf("len(passwordLetters)[%d] %% workerCount[%d] != 0", len(passwordLetters), workerCount)
	}

	letterIndexes := make(map[byte]byte)
	for i, letter := range []byte(passwordLetters) {
		letterIndexes[letter] = byte(i)
	}

	return &Cracker{
		workerCount:         workerCount,
		getOriginalPassword: getOriginalPassword,
		passwordCracked:     passwordCracked,
		letterIndexes:       letterIndexes,
	}, nil
}

func (c *Cracker) Start(ctx context.Context) <-chan error {
	done := make(chan error, 1)
	go func() {
		defer close(done)
		c.run(ctx, done)
	}()
	return done
}

func (c *Cracker) run(ctx context.Context, done chan<- error) {
	for {
		password, err := c.getOriginalPassword()
		if err != nil {
			if err == io.EOF {
				done <- nil
			} else {
				done <- err
			}
			return
		}

		if c.guessPassword(ctx, password) {
			c.passwordCracked()
		}
	}
}

func (c *Cracker) guessPassword(ctx context.Context, originalPassword string) (ok bool) {
	indexedPassword := []byte(originalPassword)
	for i, letter := range indexedPassword {
		indexedPassword[i] = c.letterIndexes[letter]
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	results := make(chan bool, c.workerCount)
	defer close(results)

	chunkSizePerWorker := byte(len(passwordLetters) / c.workerCount)
	for i := 0; i < c.workerCount; i++ {
		minFirstLetterIndex := chunkSizePerWorker * byte(i)
		maxFirstLetterIndex := minFirstLetterIndex + chunkSizePerWorker - 1
		go c.bruteForcePassword(ctx, indexedPassword, minFirstLetterIndex, maxFirstLetterIndex, results)
	}

	for i := 0; i < c.workerCount; i++ {
		guessed := <-results
		if guessed {
			ok = true
			cancel()
		}
	}

	return ok
}

func (c *Cracker) bruteForcePassword(ctx context.Context, originalPassword []byte, minFirstLetterIndex, maxFirstLetterIndex byte, results chan<- bool) {
	originalPasswordClone := make([]byte, len(originalPassword))
	copy(originalPasswordClone, originalPassword)
	originalPassword = originalPasswordClone

	startPassword := make([]byte, len(originalPassword))
	startPassword[0] = minFirstLetterIndex

	minLetterIndex := byte(len(passwordLetters) - 1)
	password := NewIndexedPassword(startPassword, minLetterIndex, maxFirstLetterIndex)
	for {
		if ctx.Err() != nil {
			results <- false
			break
		}

		for i := 0; i < iterationCount; i++ {
			if bytes.Equal(password.Data, originalPassword) {
				results <- true
				return
			}

			if !password.Increment() {
				results <- false
				return
			}
		}
	}
}
