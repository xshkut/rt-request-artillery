package main

import (
	"encoding/json"
	"fmt"
	"gra/m/v2/internal"
	"net"
	"sync"

	mapset "github.com/deckarep/golang-set"
	"github.com/pkg/errors"
	"github.com/xshkut/roletalk-go"
)

func createPeer(ipMap map[string]mapset.Set) (*roletalk.Peer, error) {
	var err error

	peer := roletalk.NewPeer(roletalk.PeerOptions{Name: "Coordinator peer", Friendly: false})
	ipMapMx := sync.RWMutex{}

	peer.Role(internal.COORDINATOR_ROLE_NAME)
	peer.AddKey("superkey", "superpassword")

	peer.Destination(internal.ATTACKER_ROLE_NAME).OnUnit(func(u *roletalk.Unit) {
		logger.Info("Connected new attacker:", u.Name(), ". Initiating...")

		res, err := peer.Destination(internal.ATTACKER_ROLE_NAME).Request(internal.INTRODUCE_EVENT, roletalk.EmitOptions{Data: "Hello!!!!", Unit: u})
		if err != nil {
			logger.Info(errors.Wrap(err, "cannot congratulate unit"))
			return
		}

		if !(res.OriginData().T == roletalk.DatatypeJSON) {
			return
		}

		result := internal.IntroduceType{}

		err = json.Unmarshal(res.OriginData().Data, &result)
		if err != nil {
			return
		}

		logger.Info("Attacker intiated:")
		logger.Info(fmt.Sprintf("%+v", result))

		ip := result.Query

		ipMapMx.Lock()
		defer ipMapMx.Unlock()

		_, ok := ipMap[ip]
		if !ok {
			ipMap[ip] = mapset.NewSet()
		}

		ipMap[ip].Add(u)

		u.OnClose(func(err error) {
			ipMapMx.Lock()
			defer ipMapMx.Unlock()

			logger.Info("Attacker disconnected:", u.Name())

			ipMap[ip].Remove(u)

			if ipMap[ip].Cardinality() == 0 {
				delete(ipMap, ip)
			}
		})
	})

	var addr net.Addr

	if addr, err = peer.Listen(*bindAddress); err != nil {
		return nil, err
	}

	logger.Info(fmt.Sprintf("Listening on %v", addr))

	return peer, nil
}
