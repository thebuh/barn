package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
	"strconv"
	"strings"
)

//Implementation of ASCOM Alpaca discovery protocol
//https://raw.githubusercontent.com/ASCOMInitiative/ASCOMRemote/main/Documentation/ASCOM%20Alpaca%20API%20Reference.pdf

const DefaultListenPort = 32227
const DefaultAlpacaApiPort = 11111

type DiscoveryServer struct {
	Conn         net.PacketConn
	ApiPort      uint32
	ListenString string
}

func NewDiscoverySever(listenPort uint32, apiPort uint32) *DiscoveryServer {
	if listenPort > 65535 || listenPort < 1 {
		listenPort = DefaultListenPort
	}
	if apiPort > 65535 || apiPort < 1 {
		apiPort = DefaultAlpacaApiPort
	}
	return &DiscoveryServer{
		ApiPort:      apiPort,
		ListenString: fmt.Sprintf(":%d", listenPort),
	}
}

// Start listening on all interfaces
func (s *DiscoveryServer) Start() {
	udpServer, err := net.ListenPacket("udp", s.ListenString)

	if err != nil {
		log.Fatal(err)
	}
	s.Conn = udpServer
	defer s.Close()
	//Listen for discovery packets on all interfaces
	for {
		buf := make([]byte, 1024)
		_, addr, err := udpServer.ReadFrom(buf)
		if err != nil {
			continue
		}
		log.Debug(fmt.Sprintf("GOT Discovery packet From %s", addr))
		msg := string(buf)
		//Only handle and reply to discovery packets 1st version
		if strings.HasPrefix(msg, "alpacadiscovery1") {
			go s.handleDiscoveryPacket(addr)
		}
	}
}
func (s *DiscoveryServer) composeDiscoveryReply() string {
	return fmt.Sprintf("{\n\"AlpacaPort\":%s\n}", strconv.Itoa(int(s.ApiPort)))
}

// Reply with our alpaca port
func (s *DiscoveryServer) handleDiscoveryPacket(addr net.Addr) {
	log.Debug(fmt.Sprintf("Sending discovery alpacaport packet to %s", addr))
	s.Conn.WriteTo([]byte(s.composeDiscoveryReply()), addr)
}
func (s *DiscoveryServer) Close() {
	s.Conn.Close()
}
