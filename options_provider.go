package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"strconv"
)

type OptionsProvider struct {
	MementoPath string
	ConfigPath  string
	Port        uint16
	UnitPort    uint16
}
type RequiredOptionsProvider struct {
	PwdHash [32]byte
	Host    string
}

func (op *OptionsProvider) Load() *OptionsProvider {
	ok := false

	op.MementoPath, ok = os.LookupEnv("MEMENTO_PATH")
	if !ok {
		op.MementoPath = "/etc/wireguard/"
	}
	op.ConfigPath, ok = os.LookupEnv("CONFIG_PATH")
	if !ok {
		op.MementoPath = "/etc/wireguard/"
	}
	_port := ""
	_port, ok = os.LookupEnv("WG_PORT")
	if !ok {
		op.Port = 51820
	} else {
		port, err := strconv.ParseUint(_port, 10, 16)
		if err != nil {
			op.Port = 51820
		} else {
			op.Port = uint16(port)
		}
	}
	_port, ok = os.LookupEnv("UNIT_PORT")
	if !ok {
		op.UnitPort = 8080
	} else {
		port, err := strconv.ParseUint(_port, 10, 16)
		if err != nil {
			op.UnitPort = 8080
		} else {
			op.UnitPort = uint16(port)
		}
	}
	return op
}
func (rop *RequiredOptionsProvider) Load() (*RequiredOptionsProvider, []error) {
	ok := false
	result := []error{}

	pwd, ok := os.LookupEnv("PASSWORD")
	if !ok {
		result = append(result, errEnvNotFound("PASSWORD"))
	}

	rop.PwdHash = sha256.Sum256([]byte(pwd))

	rop.Host, ok = os.LookupEnv("WG_HOST")
	if !ok {
		result = append(result, errEnvNotFound("WG_HOST"))
	}

	return rop, result
}

func errEnvNotFound(name string) error {
	return errors.New(fmt.Sprintf("Env variable was not found: %v", name))
}
