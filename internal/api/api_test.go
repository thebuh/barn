package api

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApiServer_ConnectClient(t *testing.T) {
	id := ClientId("127.0.0.2-1")
	var d = &Device{
		Id:               "",
		Type:             "",
		Index:            0,
		ConnectedClients: make(map[ClientId]*ConnectedClient),
	}
	d.ConnectClient(id)
	assert.Equal(t, true, d.IsConnected(id), "they should be equal")
	assert.Equal(t, false, d.IsConnected("127.0.0.1-1"), "they should be equal")
}

func TestApiServer_DisconnectClient(t *testing.T) {
	id := ClientId("127.0.0.2-1")
	var d = &Device{
		Id:               "",
		Type:             "",
		Index:            0,
		ConnectedClients: make(map[ClientId]*ConnectedClient),
	}
	d.ConnectClient(id)
	assert.Equal(t, true, d.IsConnected(id), "they should be equal")
	d.DisconnectClient(id)
	assert.Equal(t, false, d.IsConnected(id), "they should be equal")
}
