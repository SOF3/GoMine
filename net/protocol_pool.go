package net

import (
	"github.com/irmine/gomine/interfaces"
	"github.com/irmine/gomine/net/proto"
)

type ProtocolPool struct {
	protocols map[int32]interfaces.IProtocol
}

func NewProtocolPool() *ProtocolPool {
	return &ProtocolPool{make(map[int32]interfaces.IProtocol)}
}

// GetProtocol returns a protocol by its protocol number.
func (pool *ProtocolPool) GetProtocol(protocolNumber int32) interfaces.IProtocol {
	if !pool.IsProtocolRegistered(protocolNumber) {
		return pool.protocols[200]
	}
	return pool.protocols[protocolNumber]
}

// RegisterProtocol registers the given protocol.
func (pool *ProtocolPool) RegisterProtocol(protocol interfaces.IProtocol) {
	pool.protocols[protocol.GetProtocolNumber()] = protocol
}

// IsProtocolRegistered checks if a protocol with the given protocol number is registered.
func (pool *ProtocolPool) IsProtocolRegistered(protocolNumber int32) bool {
	var _, ok = pool.protocols[protocolNumber]
	return ok
}

// DeregisterProtocol deregisterss a protocol from the pool.
func (pool *ProtocolPool) DeregisterProtocol(protocolNumber int32) {
	delete(pool.protocols, protocolNumber)
}

// RegisterDefaults registers all default protocols.
func (pool *ProtocolPool) RegisterDefaults() {
	pool.RegisterProtocol(proto.NewProtocol201())
	pool.RegisterProtocol(proto.NewProtocol200())
	pool.RegisterProtocol(proto.NewProtocol160())
}
