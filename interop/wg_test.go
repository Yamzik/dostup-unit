package interop_test

import (
	"interop"
	"testing"
)

func TestShowTransfer(t *testing.T) {
	output, _ := interop.ShowTransfer()
	_ = output

	transfer, _ := output.Parse()
	_ = transfer
}

func TestSyncConf(t *testing.T) {
	interop.SyncConf()
}

func TestMap(t *testing.T) {
	m := make(map[int32]*interop.Transfer)
	transfer := &interop.Transfer{}
	m[15] = transfer
}

func TestGetKey(t *testing.T) {
	p, err := interop.GetKeypair()
	_ = err
	_ = p
}
