package core

import (
	"fmt"
	"util"
)

type UserProvider struct {
	peers *PeerRepository
}

func (up *UserProvider) New(repo *PeerRepository) *UserProvider {
	up.peers = repo
	return up
}
func (up *UserProvider) GetPeerForUser(tid uint64, pub string) (*Peer, error) {
	peer, ok := up.peers.GetByPub(pub)
	if !ok {
		return nil, util.DErr(util.NotFound, fmt.Sprintf("tid: %v, pub: %v", tid, pub))
	}
	if peer.TelegramId != tid {
		return nil, util.DErr(util.Conflict,
			fmt.Sprintf("tid: %v, pub: %v", tid, pub)).
			SetMessage("Peer with specified pub does not belong to this user")
	}

	return peer, nil
}
