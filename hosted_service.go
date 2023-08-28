package main

import (
	"core"
	"interop"
	"time"
)

func SyncTransfer(ps *core.PeerService, mp *MementoProvider, cancel chan struct{}) {
	for {
		select {
		case <-cancel:
			return
		default:
			transfer, err := interop.ShowTransfer()
			if err != nil {

			}

			data, err := transfer.Parse()
			if err != nil {

			}

			for _, t := range data {
				ps.UpdateTransfer(t)
			}

			mem := ps.GetMemento()
			mp.Save(mem)
			time.Sleep(5 * time.Second)
		}
	}
}
