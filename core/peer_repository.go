package core

import (
	"sync"
	"util"
)

type PeerRepository struct {
	mutex     *sync.RWMutex
	tidIndex  map[uint64]*[]*Peer
	addrIndex map[byte]*Peer
	pubIndex  map[string]*Peer
}

func (r *PeerRepository) New() *PeerRepository {
	r.mutex = &sync.RWMutex{}
	r.tidIndex = map[uint64]*[]*Peer{}
	r.addrIndex = map[byte]*Peer{}
	r.pubIndex = map[string]*Peer{}
	return r
}
func (r *PeerRepository) SetMemento(mem *ConfigMemento) *PeerRepository {
	for _, p := range mem.Peers {
		p.Init()
		r.Add(*p)
	}
	return r
}

func (r *PeerRepository) GetAll() []*Peer {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.toList()
}
func (r *PeerRepository) GetByTid(tid uint64) ([]*Peer, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	peers, ok := r.tidIndex[tid]
	if !ok {
		return nil, false
	}
	return *peers, true
}
func (r *PeerRepository) GetByAddr(addr byte) (*Peer, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	peer, ok := r.addrIndex[addr]
	return peer, ok
}
func (r *PeerRepository) GetByPub(pub string) (*Peer, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	peer, ok := r.pubIndex[pub]
	return peer, ok
}

func (r *PeerRepository) NextAddr() (byte, error) {
	ok := false
	addr := byte(2)
	for addr <= 254 {
		_, found := r.addrIndex[addr]
		if !found {
			ok = true
			break
		}
		addr++
	}
	if !ok {
		return 0, util.DErr(util.Conflict,
			"No free addresses to allocate").
			SetMessage("No free addresses to allocate")
	}
	return addr, nil
}

func (r *PeerRepository) Add(peer Peer) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	_, ok := r.addrIndex[peer.Addr]
	if ok {
		return util.DErr(util.Conflict,
			"Peer with specified address already exists").
			SetMessage("Peer with specified address already exists")
	}

	p := &Peer{}
	*p = peer

	peers, ok := r.tidIndex[peer.TelegramId]
	if !ok {
		peers = new([]*Peer)
		r.tidIndex[peer.TelegramId] = peers
	}
	*peers = append(*peers, p)

	r.addrIndex[peer.Addr] = p
	r.pubIndex[peer.PublicKey] = p

	return nil
}
func (r *PeerRepository) Delete(peer *Peer) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	peers, ok := r.tidIndex[peer.TelegramId]
	if ok {
		for i, p := range *peers {
			if p.PublicKey == peer.PublicKey {
				if i == len(*peers)-1 {
					*peers = (*peers)[:i]
				} else {
					*peers = append((*peers)[:i], (*peers)[i+1:]...)
				}
				if len(*peers) == 0 {
					delete(r.tidIndex, peer.TelegramId)
				}
				break
			}
		}
	}
	delete(r.addrIndex, peer.Addr)
	delete(r.pubIndex, peer.PublicKey)
}

func (r *PeerRepository) toList() []*Peer {
	peers := make([]*Peer, len(r.pubIndex), len(r.pubIndex))
	i := 0
	for _, p := range r.pubIndex {
		peers[i] = p
		i++
	}
	return peers
}
