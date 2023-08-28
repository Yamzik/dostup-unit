package core

import (
	"fmt"
	"interop"
	"time"
)

type Interface struct {
	PrivateKey string
	PublicKey  string
	Address    string
	PreUp      string
	PostUp     string
	PreDown    string
	PostDown   string
}

type Peer struct {
	PrivateKey    string
	PublicKey     string
	PresharedKey  string
	AllowedIPs    string
	TelegramId    uint64
	SRx           uint64
	STx           uint64
	Rx            uint64
	Tx            uint64
	HandshakeTime time.Duration
	IsEnabled     bool
	Addr          byte
}

func (i *Interface) New(port uint16) error {
	keypair, err := interop.GetKeypair()
	if err != nil {
		return err
	}

	i.PrivateKey = keypair.PrivateKey
	i.PublicKey = keypair.PublicKey
	i.Address = "10.8.0.1/24"
	i.PostUp = defaultPostUp(port)

	return nil
}
func (i *Interface) ToConfig(port uint16) string {
	return fmt.Sprintf(
		"[Interface]\n"+
			"PrivateKey = %v\n"+
			"Address = %v\n"+
			"ListenPort = %v\n"+
			"PostUp = %v\n", i.PrivateKey, i.Address, port, i.PostUp)
}

func (p *Peer) New(tid uint64, addr byte) (*Peer, error) {
	keypair, err := interop.GetKeypair()
	if err != nil {
		return nil, err
	}
	psk, err := interop.GetPSK()
	if err != nil {
		return nil, err
	}

	p.TelegramId = tid
	p.PrivateKey = keypair.PrivateKey
	p.PublicKey = keypair.PublicKey
	p.PresharedKey = psk
	p.Addr = addr
	p.AllowedIPs = getPeerIP(addr)
	p.IsEnabled = true

	return p, nil
}
func (p *Peer) Init() *Peer {
	p.SRx = p.Rx
	p.STx = p.Tx
	return p
}

func (p *Peer) ToConfig() string {
	if !p.IsEnabled {
		return ""
	}
	return fmt.Sprintf(
		"# TelegramId = %v\n"+
			"[Peer]\n"+
			"PublicKey = %v\n"+
			"PresharedKey = %v\n"+
			"AllowedIPs = %v\n", p.TelegramId, p.PublicKey, p.PresharedKey, p.AllowedIPs)
}

func (p *Peer) ToInterface(i *Interface, endpoint string) string {
	return fmt.Sprintf(
		"[Interface]\n"+
			"PrivateKey = %v\n"+
			"Address = %v\n"+
			"DNS = 1.1.1.1\n"+
			"\n"+
			"[Peer]\n"+
			"PublicKey = %v\n"+
			"PresharedKey = %v\n"+
			"AllowedIPs = 0.0.0.0/0, ::/0\n"+
			"PersistentKeepalive = 0\n"+
			"Endpoint = %v\n", p.PrivateKey, p.AllowedIPs, i.PublicKey, p.PresharedKey, endpoint)
}

func (p *Peer) Disable() {
	p.IsEnabled = false
}
func (p *Peer) Enable() {
	p.IsEnabled = true
}

func (p *Peer) AddTx(tx uint64) {
	p.Tx = p.STx + tx
}
func (p *Peer) AddRx(rx uint64) {
	p.Rx = p.SRx + rx
}

func defaultPostUp(port uint16) string {
	form := "iptables -t nat -A POSTROUTING -s 10.8.0.0/24 -o eth0 -j MASQUERADE; " +
		"iptables -A INPUT -p udp -m udp --dport %v -j ACCEPT; " +
		"iptables -A FORWARD -i wg0 -j ACCEPT; " +
		"iptables -A FORWARD -o wg0 -j ACCEPT;"
	return fmt.Sprintf(form, port)
}
func getPeerIP(addr byte) string {
	return fmt.Sprintf("10.8.0.%v/32", addr)
}
