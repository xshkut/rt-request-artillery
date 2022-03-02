package main

import (
	"gra/m/v2/internal"
	"log"

	"github.com/pkg/errors"
	"github.com/xshkut/roletalk-go"
)

func createPeer() *roletalk.Peer {
	peer := roletalk.NewPeer(roletalk.PeerOptions{Name: "Attacker peer", Friendly: false})

	peer.Role(internal.ATTACKER_ROLE_NAME)
	peer.AddKey("superkey", "superpassword")

	coordinator := peer.Destination(internal.COORDINATOR_ROLE_NAME)

	coordinator.OnClose(func() {
		log.Println("Coordinator disconnected")
	})

	coordinator.OnUnit(func(u *roletalk.Unit) {
		log.Println("Connected coordinator:", u.Name())
	})

	peer.Role(internal.ATTACKER_ROLE_NAME).OnRequest(internal.INTRODUCE_EVENT, func(im *roletalk.RequestContext) {
		myInfo, error := internal.GetMyIp()
		if error != nil {
			err := errors.Wrap(error, "cannot get my ip")
			im.Reject(err)
			log.Fatalln(err.Error())
			return
		}

		log.Println("My info:", myInfo)

		im.Reply(myInfo)
	})

	return peer
}

func connectPeer(peer *roletalk.Peer, address string) {
	_, err := peer.Connect(address, roletalk.ConnectOptions{})
	if err != nil {
		log.Println("Cannot connect to coordinator. Will periodically try... Error:", err.Error())
	}
}
