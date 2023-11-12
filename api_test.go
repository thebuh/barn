package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApiServer_ConnectClient(t *testing.T) {
	ip := ClientIp("127.0.0.2")
	var d = Device{
		Id:               "",
		Type:             "",
		Index:            0,
		ConnectedClients: nil,
	}
	d.ConnectClient(ip)
	assert.Equal(t, true, d.IsConnected(ip), "they should be equal")
	assert.Equal(t, false, d.IsConnected("127.0.0.1"), "they should be equal")
}

func TestApiServer_DisconnectClient(t *testing.T) {
	ip := ClientIp("127.0.0.2")
	var d = Device{
		Id:               "",
		Type:             "",
		Index:            0,
		ConnectedClients: nil,
	}
	d.ConnectClient(ip)
	assert.Equal(t, true, d.IsConnected(ip), "they should be equal")
	d.DisconnectClient(ip)
	assert.Equal(t, false, d.IsConnected(ip), "they should be equal")
}
