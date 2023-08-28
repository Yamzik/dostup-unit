package core

type ConfigMemento struct {
	Interface *Interface
	Peers     []*Peer
}

func (cm *ConfigMemento) New(port uint16) *ConfigMemento {
	i := &Interface{}
	i.New(port)
	cm.Interface = i
	cm.Peers = make([]*Peer, 0)
	return cm
}
