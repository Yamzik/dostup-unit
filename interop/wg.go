package interop

import (
	"fmt"
	"strings"
)

type Keypair struct {
	PrivateKey string
	PublicKey  string
}

type Transfer struct {
	PublicKey string
	Tx        uint64
	Rx        uint64
}

type LatestHandshake struct {
	PublicKey string
	Timestamp int64
}

type Interface struct {
	PublicKey     string `json:"public_key"`
	ListeningPort int    `json:"listening_port"`
}

type Peer struct {
	PublicKey       string `json:"public_key"`
	PresharedKey    string `json:"preshared_key,omitempty"`
	Endpoint        string `json:"endpoint"`
	AllowedIPs      string `json:"allowed_ips"`
	LatestHandshake string `json:"latest_handshake,omitempty"`
	Transfer        string `json:"transfer,omitempty"`
}

type ShowLog struct {
	Interface *Interface `json:"interface"`
	Peers     []*Peer    `json:"peers"`
}

func ShowTransfer() (*TransferWrapper, error) {
	result, err := execInternal("wg", "show", "wg0", "transfer")
	tw := &TransferWrapper{data: result}
	return tw, err
}

func ShowLatestHandshakes() (string, error) {
	result, err := execInternal("wg", "show", "wg0", "latest-handshakes")
	return result, err
}

func GetKeypair() (*Keypair, error) {
	keypair := Keypair{}
	var err error

	keypair.PrivateKey, err = execInternal("wg", "genkey")
	if err != nil {
		return nil, err
	}
	keypair.PrivateKey = strings.Trim(keypair.PrivateKey, "\n")
	keypair.PublicKey, err = execInternal("bash", "-c", fmt.Sprintf("echo %v | wg pubkey", keypair.PrivateKey))
	if err != nil {
		return nil, err
	}
	keypair.PublicKey = strings.Trim(keypair.PublicKey, "\n")

	return &keypair, nil
}

func GetPSK() (string, error) {
	key, err := execInternal("wg", "genpsk")
	key = strings.Trim(key, "\n")
	if err != nil {
		return "", err
	}
	return key, nil
}

func SyncConf() error {
	_, err := execInternal("bash", "-c", "wg syncconf wg0 <(wg-quick strip wg0)")
	return err
}

func Up() error {
	_, err := execInternal("bash", "-c", "wg-quick up wg0")
	return err
}

func Down() error {
	_, err := execInternal("bash", "-c", "wg-quick down wg0")
	return err
}
