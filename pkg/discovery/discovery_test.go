package discovery

import (
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
	"time"
)

func TestDiscoveryServer_DiscoveryReply(t *testing.T) {
	var server = NewDiscoverySever(32227, 12345)
	go server.Start()
	time.Sleep(1 * time.Second)
	defer server.Close()
	addr, _ := net.ResolveUDPAddr("udp4", "127.0.0.1:32227")
	conn, _ := net.ListenPacket("udp", ":0")
	conn.WriteTo([]byte("alpacadiscovery1"), addr)
	buf := make([]byte, 22)
	conn.ReadFrom(buf)
	conn.Close()
	msg := string(buf)
	assert.Equal(t, "{\n\"AlpacaPort\":12345\n}", msg, "they should be equal")
}

func TestDiscoveryServer_BindPortDefault(t *testing.T) {
	var server = NewDiscoverySever(32227, 11111)
	assert.Equal(t, ":32227", server.ListenString, "they should be equal")
}

func TestDiscoveryServer_BindPortCustom(t *testing.T) {
	var server = NewDiscoverySever(33227, 11111)
	assert.Equal(t, ":33227", server.ListenString, "they should be equal")
}

func TestDiscoveryServer_BindPortOutOfRange(t *testing.T) {
	var server = NewDiscoverySever(70000, 11111)
	assert.Equal(t, ":32227", server.ListenString, "they should be equal")
}

func TestDiscoveryServer_BindPortOutOfRange2(t *testing.T) {
	var server = NewDiscoverySever(0, 11111)
	assert.Equal(t, ":32227", server.ListenString, "they should be equal")
}

func TestDiscoveryServer_AlpacaPortReply(t *testing.T) {
	var server = NewDiscoverySever(32227, 11111)
	assert.Equal(t, "{\n\"AlpacaPort\":11111\n}", server.composeDiscoveryReply(), "they should be equal")
}

func TestDiscoveryServer_AlpacaPortReplyCustom(t *testing.T) {
	var server = NewDiscoverySever(32227, 80)
	assert.Equal(t, "{\n\"AlpacaPort\":80\n}", server.composeDiscoveryReply(), "they should be equal")
}

func TestDiscoveryServer_AlpacaPortOutOfRange(t *testing.T) {
	var server = NewDiscoverySever(32227, 70000)
	assert.Equal(t, "{\n\"AlpacaPort\":11111\n}", server.composeDiscoveryReply(), "they should be equal")
}

func TestDiscoveryServer_AlpacaPortOutOfRange2(t *testing.T) {
	var server = NewDiscoverySever(32227, 0)
	assert.Equal(t, "{\n\"AlpacaPort\":11111\n}", server.composeDiscoveryReply(), "they should be equal")
}
