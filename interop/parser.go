package interop

import (
	"errors"
	"strconv"
	"strings"
)

type TransferWrapper struct {
	data string
}

func (tw *TransferWrapper) Parse() ([]*Transfer, error) {
	return parseShowTransfer(tw.data)
}

func parseShowTransfer(input string) ([]*Transfer, error) {
	lines := strings.Split(input, "\n")
	result := []*Transfer{}
	for _, line := range lines[:len(lines)-1] {
		split := strings.Split(line, "\t")
		transfer := Transfer{}
		transfer.PublicKey = split[0]
		transfer.Rx, _ = strconv.ParseUint(split[1], 10, 64)
		transfer.Tx, _ = strconv.ParseUint(split[2], 10, 64)
		result = append(result, &transfer)
	}
	return result, nil
}
func parseShow(input string) (*ShowLog, error) {
	lines := strings.Split(input, "\n")
	result := &ShowLog{}
	var currentPeer *Peer

	for _, line := range lines {
		if strings.HasPrefix(line, "interface:") {
			result.Interface = &Interface{}
		} else if strings.HasPrefix(line, "peer:") {
			currentPeer = &Peer{}
			parts := strings.Split(line, ": ")
			if len(parts) == 2 {
				currentPeer.PublicKey = strings.TrimSpace(parts[1])
			} else {
				return nil, errors.New("Invalid peer output, expected format: \"peer: PublicKey\"")
			}
			result.Peers = append(result.Peers, currentPeer)
		} else if currentPeer != nil {
			parts := strings.Split(line, ": ")
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])

				switch key {
				case "preshared key":
					currentPeer.PresharedKey = value
				case "endpoint":
					currentPeer.Endpoint = value
				case "allowed ips":
					currentPeer.AllowedIPs = value
				case "latest handshake":
					currentPeer.LatestHandshake = value
				case "transfer":
					currentPeer.Transfer = value
				}
			}
		} else if result.Interface != nil {
			parts := strings.Split(line, ": ")
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				switch key {
				case "public key":
					result.Interface.PublicKey = value
				case "listening port":
					port, err := strconv.Atoi(value)
					if err != nil {
						return nil, err
					}
					result.Interface.ListeningPort = port
				}
			}
		}
	}

	return result, nil
}
