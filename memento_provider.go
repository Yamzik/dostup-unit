package main

import (
	"core"
	"encoding/json"
	"os"
)

type MementoProvider struct {
	port        uint16
	mementoPath string
	configPath  string
}

func (mp *MementoProvider) New(port uint16, mementoPath string, configPath string) *MementoProvider {
	mp.port = port
	mp.mementoPath = mementoPath
	mp.configPath = configPath
	return mp
}
func (mp *MementoProvider) Load() (*core.ConfigMemento, error) {
	b, err := os.ReadFile(mp.mementoPath)

	memento := &core.ConfigMemento{}
	if err != nil && os.IsNotExist(err) || len(b) == 0 {
		inner_err := os.WriteFile(mp.mementoPath, []byte{}, 0777)
		if inner_err != nil {
			return nil, inner_err
		}
		return memento.New(mp.port), nil
	} else if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, memento)
	if err != nil {
		return nil, err
	}

	return memento, nil
}
func (mp *MementoProvider) Save(mem *core.ConfigMemento) error {
	data, err := json.MarshalIndent(mem, "", "\t")
	if err != nil {
		return err
	}
	err = os.WriteFile(mp.mementoPath, data, os.ModeAppend)
	if err != nil {
		return err
	}

	return nil
}
func (mp *MementoProvider) SaveConfig(config string) error {
	err := os.WriteFile(mp.configPath, []byte(config), 0777)
	return err
}
