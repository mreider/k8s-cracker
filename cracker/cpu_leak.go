package main

import (
	"math/rand"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	defaultCPULeakSecondsTo100 = 60
)

func startCPULeakIfNeed(cpuLeak string) {
	if cpuLeak == "" {
		return
	}

	cpuLeak = strings.ToUpper(cpuLeak)
	cpuLeakSecondsTo100, err := strconv.Atoi(cpuLeak)
	if err != nil {
		if cpuLeak != "TRUE" {
			return
		}
		cpuLeakSecondsTo100 = defaultCPULeakSecondsTo100
	}

	log.Info().Msgf("CPU Leak: %d sec to 100%%", cpuLeakSecondsTo100)
	go runCPULoad(cpuLeakSecondsTo100)
}

// Modified version of https://github.com/vikyd/go-cpu-load/blob/master/cpu_load.go
func runCPULoad(seconds int) {
	percentageStep := 100.0 / float64(seconds)
	unitHundredsOfMicrosecond := 1000.0
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
			runtime.LockOSThread()
			percentage := 0.0
			var runMicrosecond, sleepMicrosecond int
			for {
				if percentage < 100 {
					percentage += percentageStep
					if percentage > 100 {
						percentage = 100
					}
					runMicrosecond = int(unitHundredsOfMicrosecond * percentage)
					sleepMicrosecond = int(unitHundredsOfMicrosecond*100 - float64(runMicrosecond))
				}

				if percentage >= 100 {
					j := 0
					for j >= 0 {
						j = j * j % 17
					}
				}

				for k := 0; k < 10; k++ {
					begin := time.Now()
					for {
						// run 100%
						if time.Since(begin) > time.Duration(runMicrosecond)*time.Microsecond {
							break
						}
					}
					// sleep
					if sleepMicrosecond > 0 {
						time.Sleep(time.Duration(sleepMicrosecond) * time.Microsecond)
					}
				}
			}
		}()
	}
}
