package net

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"pioneer-server/io"
	"pioneer-server/net/protocol"

	"github.com/google/uuid"
)

type TcpClient struct {
	Id   uuid.UUID
	Conn net.Conn
}

type TcpServer struct {
	Addr     string
	Listener *net.Listener
	Clients  map[string]TcpClient
	handlers map[int]tcpPacketHandlerWrapper
}

type TcpPacketHandler func(client *TcpClient, packet protocol.Packet)
type tcpPacketHandlerWrapper func(client *TcpClient, buffer protocol.PacketBuffer)
type PacketFactory func() protocol.Packet

func CreateTcpServer(onlyLocal bool, port int) TcpServer {
	var host string = "0.0.0.0"
	if onlyLocal {
		host = "127.0.0.1"
	}

	return TcpServer{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Clients:  make(map[string]TcpClient),
		handlers: make(map[int]tcpPacketHandlerWrapper),
	}
}

func (server *TcpServer) Start() {
	io.Logf(io.Info, "Starting TCP server on %s...", server.Addr)

	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		io.Logf(io.Error, "Failed to start TCP server on %s: %s", server.Addr, err)
		os.Exit(1)
	}

	server.Listener = &listener
	io.Logf(io.Info, "TCP server listening on %s! Awaiting connections...", server.Addr)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				io.Log(io.Error, "Failed to accept connection")
				continue
			}

			client := TcpClient{
				Id:   uuid.New(),
				Conn: conn,
			}

			server.Clients[client.Id.String()] = client
			io.Logf(io.Debug, "Accepted connection from %s", client.Conn.RemoteAddr().String())

			go server.handleConnection(client)
		}
	}()
}

func (server *TcpServer) handleConnection(client TcpClient) {
	defer client.Conn.Close()

	for {
		data, err := bufio.NewReader(client.Conn).ReadBytes('\n')
		if err != nil {
			io.Logf(io.Error, "Failed to read data from client #%s", client.Id)
			break
		}

		io.Logf(io.Debug, "Received data from client #%s: %s", client.Id, string(data))

		id, buffer := protocol.ParsePacketBuffer(data)
		if server.handlers[id] == nil {
			io.Logf(io.Warn, "Received unknown packet from client #%s: %d", client.Id, id)
			continue
		}

		server.handlers[id](&client, buffer)
	}

	delete(server.Clients, client.Id.String())

	io.Logf(io.Error, "Closed client connection #%s from %s", client.Id, client.Conn.RemoteAddr().String())
}

func (server *TcpServer) Stop() {
	io.Logf(io.Info, "Stopping TCP server on %s...", server.Addr)

	if server.Listener != nil {
		(*server.Listener).Close()
	}
}

func (server *TcpServer) RegisterHandler(handler TcpPacketHandler, factory PacketFactory) {
	empty := factory()
	id := empty.GetId()

	server.handlers[id] = func(client *TcpClient, buffer protocol.PacketBuffer) {
		packet := factory()

		protocol.Deserialize(buffer, packet)
		handler(client, packet)
	}
}

func (client *TcpClient) Send(packet protocol.Packet) {
	buffer := protocol.Serialize(packet)
	client.Conn.Write(buffer.Bytes)
}

func (client *TcpClient) Close() {
	client.Conn.Close()
}
