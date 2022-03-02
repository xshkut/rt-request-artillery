package main

import (
	"flag"
	"gra/m/v2/internal"
	"os"
	"sync"
)

var help = flag.Bool("help", false, "Print help")
var coordinatorAddress = flag.String("connect", internal.COORDINATOR_DEFAULT_ADDRESS, "Address of coordinator server (with protocol)")
var maxRate = flag.Int("max-rate", 100, "Maximal rate of outgoing requests")

func init() {
	flag.BoolVar(help, "h", *help, "alias for --help")
	flag.StringVar(coordinatorAddress, "c", *coordinatorAddress, "alias for --connect")
	flag.IntVar(maxRate, "m", *maxRate, "alias for --max-rate")

	flag.Parse()

	if *help {
		internal.ArgsUsage()
		os.Exit(0)
	}
}

func main() {
	wg := sync.WaitGroup{}

	attackVectorCh := make(chan internal.AttackVector)
	globalRateLimiter := make(chan interface{})

	peer := createPeer()

	go internal.CreateRateLimiter(float64(*maxRate), globalRateLimiter)

	go consumeAttackVectors(peer, attackVectorCh)

	go distributeAttackVectors(attackVectorCh, globalRateLimiter)

	go connectPeer(peer, *coordinatorAddress)

	wg.Add(1)

	wg.Wait()
}
