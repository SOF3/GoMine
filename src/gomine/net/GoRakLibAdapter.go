package net

import (
	"gomine/interfaces"
	server2 "goraklib/server"
	"gomine/net/info"
	"goraklib/protocol"
)

type GoRakLibAdapter struct {
	server interfaces.IServer
	rakLibServer *server2.GoRakLibServer
}

/**
 * Returns a new GoRakLib adapter to adapt to the RakNet server.
 */
func NewGoRakLibAdapter(server interfaces.IServer) *GoRakLibAdapter {
	var rakServer = server2.NewGoRakLibServer(server.GetName(), server.GetAddress(), server.GetPort())
	rakServer.SetMinecraftProtocol(info.LatestProtocol)
	rakServer.SetMinecraftVersion(info.GameVersionNetwork)
	rakServer.SetServerName(server.GetName())
	rakServer.SetMaxConnectedSessions(server.GetMaximumPlayers())
	rakServer.SetConnectedSessionCount(0)
	rakServer.SetDefaultGameMode("Creative")
	rakServer.SetMotd(server.GetMotd())

	InitPacketPool()
	InitHandlerPool()

	return &GoRakLibAdapter{server, rakServer}
}

/**
 * Ticks the adapter
 */
func (adapter *GoRakLibAdapter) Tick() {
	go adapter.rakLibServer.Tick()

	go func() {
		for _, session := range adapter.rakLibServer.GetSessionManager().GetSessions() {
			for _, encapsulatedPacket := range session.GetReadyEncapsulatedPackets() {

				batch := NewMinecraftPacketBatch()
				batch.stream.Buffer = encapsulatedPacket.Buffer
				batch.Decode()

				for _, packet := range batch.GetPackets() {
					packet.DecodeHeader()
					packet.Decode()

					var player, _ = adapter.server.GetPlayerFactory().GetPlayerBySession(session.GetAddress(), session.GetPort())

					handlers := GetPacketHandlers(packet.GetId())
					for _, handler := range handlers {
						handler.Handle(packet, player, session, adapter.server)
					}
				}
			}
		}
	}()
}

func (adapter *GoRakLibAdapter) SendPacket(pk interfaces.IPacket, session *server2.Session) {
	pk.EncodeHeader()
	pk.Encode()
	var b = NewMinecraftPacketBatch()
	b.AddPacket(pk)

	adapter.SendBatch(&b, session)
}

func (adapter *GoRakLibAdapter) SendBatch(batch interfaces.IMinecraftPacketBatch, session *server2.Session) {
	batch.Encode()

	var encPacket = protocol.NewEncapsulatedPacket()
	encPacket.SetBuffer(batch.GetStream().GetBuffer())

	var datagram = protocol.NewDatagram()
	datagram.AddPacket(&encPacket)
	datagram.Encode()

	adapter.rakLibServer.GetSessionManager().SendPacket(datagram, session)
}