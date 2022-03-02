package internal

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

var ArgsUsage = func() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

func CreateRateLimiter(rate float64, ch <-chan interface{}) {
	var waitTime time.Duration = time.Duration(int64(1.0 / rate * 1000000))

	log.Println(rate, waitTime)

	for {
		time.Sleep(time.Microsecond * waitTime)
		<-ch
	}
}
