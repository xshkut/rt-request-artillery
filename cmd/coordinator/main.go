package main

import (
	"flag"
	"gra/m/v2/internal"
	"os"
	"sync"

	mapset "github.com/deckarep/golang-set"
	"github.com/pkg/errors"
	"github.com/withmandala/go-log"
	"github.com/xshkut/roletalk-go"
)

var help = flag.Bool("help", false, "Print help")
var bindAddress = flag.String("bind", "0.0.0.0:9000", "Bind address")

var logger *log.Logger = log.New(os.Stdout)

func init() {
	flag.BoolVar(help, "h", *help, "alias for --help")
	flag.StringVar(bindAddress, "b", *bindAddress, "alias for --bind")

	flag.Parse()

	if *help {
		internal.ArgsUsage()
		os.Exit(0)
	}
}

func main() {
	ipMap := make(map[string]mapset.Set)

	config, err := getConfig("targets.yml")
	if err != nil {
		logger.Fatal(errors.Wrap(err, "cannot read config"))
		return
	}

	peer, err := createPeer(ipMap)
	if err != nil {
		logger.Fatal(errors.Wrap(err, "Cannot listen"))
	}

	for _, ac := range config {
		go coordinate(peer, ipMap, ac)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}

func coordinate(peer *roletalk.Peer, ipMap map[string]mapset.Set, ac internal.AttackConfig) {
	addrStateCh := make(chan addressState)
	rateCh := make(chan float64, 1)

	go startCheckingAddress(ac.Address, addrStateCh)

	go updateRate(addrStateCh, rateCh, peer)

	err := downstreamAttackVectors(ac, rateCh, peer, ipMap)
	if err != nil {
		logger.Fatal(errors.Wrap(err, "cannot comsume rate"))
	}
}
