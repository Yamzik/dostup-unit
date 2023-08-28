package core

import (
	"fmt"
	"interop"
	"strings"
)

type PeerService struct {
	peerRepo     *PeerRepository
	userProvider *UserProvider

	_interface *Interface

	host string
	port uint16
}

type PeerSM struct {
	TelegramId uint64
	PublicKey  string
	Tx         uint64
	Rx         uint64
}

func (psm *PeerSM) New(p *Peer) *PeerSM {
	psm.TelegramId = p.TelegramId
	psm.PublicKey = p.PublicKey
	psm.Tx = p.Tx
	psm.Rx = p.Rx
	return psm
}

func (ps *PeerService) New(pr *PeerRepository, up *UserProvider, host string, port uint16) *PeerService {
	ps.peerRepo = pr
	ps.userProvider = up
	ps.host = host
	ps.port = port
	return ps
}
func (ps *PeerService) SetMemento(mem *ConfigMemento) *PeerService {
	ps._interface = mem.Interface
	return ps
}
func (ps *PeerService) GetMemento() *ConfigMemento {
	return &ConfigMemento{
		Interface: ps._interface,
		Peers:     ps.peerRepo.toList(),
	}
}

func (ps *PeerService) GetAllPeers() []*PeerSM {
	peers := ps.peerRepo.GetAll()
	return mapPeers(peers)
}
func (ps *PeerService) GetPeers(tid uint64) []*PeerSM {
	peers, ok := ps.peerRepo.GetByTid(tid)
	if !ok {
		return make([]*PeerSM, 0)
	}
	return mapPeers(peers)
}
func (ps *PeerService) GetPeer(tid uint64, pub string) (*PeerSM, error) {
	peer, err := ps.userProvider.GetPeerForUser(tid, pub)
	if err != nil {
		return nil, err
	}
	return new(PeerSM).New(peer), nil
}

func (ps *PeerService) GetPeerInterface(tid uint64, pub string) (string, error) {
	peer, err := ps.userProvider.GetPeerForUser(tid, pub)
	if err != nil {
		return "", err
	}
	endpoint := fmt.Sprintf("%v:%v", ps.host, ps.port)
	return peer.ToInterface(ps._interface, endpoint), nil
}

func (ps *PeerService) Add(tid uint64) (*PeerSM, error) {
	addr, err := ps.peerRepo.NextAddr()
	if err != nil {
		return nil, err
	}

	peer, err := (&Peer{}).New(tid, addr)
	if err != nil {
		return nil, err
	}

	err = ps.peerRepo.Add(*peer)
	if err != nil {
		return nil, err
	}

	return new(PeerSM).New(peer), nil
}
func (ps *PeerService) Delete(tid uint64, pub string) error {
	peer, err := ps.userProvider.GetPeerForUser(tid, pub)
	if err != nil {
		return err
	}

	ps.peerRepo.Delete(peer)
	return nil
}

func (ps *PeerService) DisablePeer(tid uint64, pub string) error {
	peer, err := ps.userProvider.GetPeerForUser(tid, pub)
	if err != nil {
		return err
	}
	peer.Disable()
	return nil
}
func (ps *PeerService) EnablePeer(tid uint64, pub string) error {
	peer, err := ps.userProvider.GetPeerForUser(tid, pub)
	if err != nil {
		return err
	}
	peer.Enable()
	peer.SRx = peer.Rx
	peer.STx = peer.Tx
	return nil
}

func (ps *PeerService) MakeConfig() (string, error) {
	builder := strings.Builder{}
	_, err := builder.WriteString(ps._interface.ToConfig(ps.port))
	for _, peer := range ps.peerRepo.GetAll() {
		_, err = builder.WriteString(peer.ToConfig())
	}
	if err != nil {
		return "", err
	}
	return builder.String(), nil
}
func (ps *PeerService) UpdateTransfer(t *interop.Transfer) {
	peer, ok := ps.peerRepo.GetByPub(t.PublicKey)
	if !ok {

	}
	peer.AddTx(t.Tx)
	peer.AddRx(t.Rx)
}

func mapPeers(peers []*Peer) []*PeerSM {
	r := make([]*PeerSM, len(peers), len(peers))
	for i, p := range peers {
		r[i] = new(PeerSM).New(p)
	}
	return r
}
