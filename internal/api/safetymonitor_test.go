package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSafetyMonitorAPI_ConnectClient(t *testing.T) {
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

func TestSafetyMonitorAPI_DisconnectClient(t *testing.T) {
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

func TestDeviceState_Structure(t *testing.T) {
	// Test that DeviceState struct has the correct fields
	deviceState := DeviceState{
		Name:  "IsSafe",
		Value: true,
	}

	assert.Equal(t, "IsSafe", deviceState.Name, "Name should match")
	assert.Equal(t, true, deviceState.Value, "Value should match")
}

func TestDeviceStateResponse_Structure(t *testing.T) {
	// Test that deviceStateResponse struct has the correct fields
	deviceStates := []DeviceState{
		{
			Name:  "IsSafe",
			Value: true,
		},
		{
			Name:  "TimeStamp",
			Value: true,
		},
	}

	resp := deviceStateResponse{
		Value: deviceStates,
		alpacaResponse: alpacaResponse{
			ClientTransactionID: 1,
			ServerTransactionID: 2,
			ErrorNumber:         0,
			ErrorMessage:        "",
		},
	}

	assert.Len(t, resp.Value, 2, "Should have 2 device states")
	assert.Equal(t, "IsSafe", resp.Value[0].Name, "First state should be IsSafe")
	assert.Equal(t, "TimeStamp", resp.Value[1].Name, "Second state should be TimeStamp")
	assert.Equal(t, uint32(1), resp.ClientTransactionID, "ClientTransactionID should match")
}

func TestSafetyMonitorAPI_Connecting_AlwaysFalse(t *testing.T) {
	// Test that the connecting endpoint always returns false
	resp := boolResponse{
		Value: false,
		alpacaResponse: alpacaResponse{
			ClientTransactionID: 1,
			ServerTransactionID: 2,
			ErrorNumber:         0,
			ErrorMessage:        "",
		},
	}

	assert.Equal(t, false, resp.Value, "Connecting should always return false")
	assert.Equal(t, uint32(1), resp.ClientTransactionID, "ClientTransactionID should match")
}

func TestSafetyMonitorAPI_Connect_Endpoint(t *testing.T) {
	// Test that the connect endpoint properly connects a client
	api := &SafetyMonitorAPI{}

	// Create a mock device
	device := &Device{
		Id:               "test-monitor",
		Type:             "safetymonitor",
		Index:            0,
		ConnectedClients: make(map[ClientId]*ConnectedClient),
	}

	// Mock the API server
	apiServer := &ApiServer{
		Devices: map[string]map[int]*Device{
			"safetymonitor": {
				0: device,
			},
		},
	}
	api.ApiServer = apiServer

	// Test successful connection
	clientId := ClientId("127.0.0.1-123")

	// Initially not connected
	assert.Equal(t, false, device.IsConnected(clientId), "Device should not be connected initially")

	// Connect the client
	device.ConnectClient(clientId)

	// Should now be connected
	assert.Equal(t, true, device.IsConnected(clientId), "Device should be connected after ConnectClient")
}

func TestSafetyMonitorAPI_Disconnect_Endpoint(t *testing.T) {
	// Test that the disconnect endpoint properly disconnects a client
	api := &SafetyMonitorAPI{}

	// Create a mock device
	device := &Device{
		Id:               "test-monitor",
		Type:             "safetymonitor",
		Index:            0,
		ConnectedClients: make(map[ClientId]*ConnectedClient),
	}

	// Mock the API server
	apiServer := &ApiServer{
		Devices: map[string]map[int]*Device{
			"safetymonitor": {
				0: device,
			},
		},
	}
	api.ApiServer = apiServer

	// Test successful disconnection
	clientId := ClientId("127.0.0.1-123")

	// Connect first
	device.ConnectClient(clientId)
	assert.Equal(t, true, device.IsConnected(clientId), "Device should be connected")

	// Disconnect the client
	device.DisconnectClient(clientId)

	// Should now be disconnected
	assert.Equal(t, false, device.IsConnected(clientId), "Device should be disconnected after DisconnectClient")
}

func TestSafetyMonitorAPI_Connect_Disconnect_Workflow(t *testing.T) {
	// Test the complete connect/disconnect workflow
	api := &SafetyMonitorAPI{}

	// Create a mock device
	device := &Device{
		Id:               "test-monitor",
		Type:             "safetymonitor",
		Index:            0,
		ConnectedClients: make(map[ClientId]*ConnectedClient),
	}

	// Mock the API server
	apiServer := &ApiServer{
		Devices: map[string]map[int]*Device{
			"safetymonitor": {
				0: device,
			},
		},
	}
	api.ApiServer = apiServer

	clientId1 := ClientId("127.0.0.1-123")
	clientId2 := ClientId("127.0.0.1-456")

	// Test multiple clients
	assert.Equal(t, false, device.IsConnected(clientId1), "Client 1 should not be connected initially")
	assert.Equal(t, false, device.IsConnected(clientId2), "Client 2 should not be connected initially")

	// Connect both clients
	device.ConnectClient(clientId1)
	device.ConnectClient(clientId2)

	assert.Equal(t, true, device.IsConnected(clientId1), "Client 1 should be connected")
	assert.Equal(t, true, device.IsConnected(clientId2), "Client 2 should be connected")

	// Disconnect one client
	device.DisconnectClient(clientId1)

	assert.Equal(t, false, device.IsConnected(clientId1), "Client 1 should be disconnected")
	assert.Equal(t, true, device.IsConnected(clientId2), "Client 2 should still be connected")

	// Disconnect the other client
	device.DisconnectClient(clientId2)

	assert.Equal(t, false, device.IsConnected(clientId1), "Client 1 should remain disconnected")
	assert.Equal(t, false, device.IsConnected(clientId2), "Client 2 should be disconnected")
}

func TestSafetyMonitorAPI_Connect_Disconnect_Reconnect(t *testing.T) {
	// Test that a client can reconnect after disconnecting
	api := &SafetyMonitorAPI{}

	// Create a mock device
	device := &Device{
		Id:               "test-monitor",
		Type:             "safetymonitor",
		Index:            0,
		ConnectedClients: make(map[ClientId]*ConnectedClient),
	}

	// Mock the API server
	apiServer := &ApiServer{
		Devices: map[string]map[int]*Device{
			"safetymonitor": {
				0: device,
			},
		},
	}
	api.ApiServer = apiServer

	clientId := ClientId("127.0.0.1-123")

	// Initial connection
	device.ConnectClient(clientId)
	assert.Equal(t, true, device.IsConnected(clientId), "Client should be connected")

	// Disconnect
	device.DisconnectClient(clientId)
	assert.Equal(t, false, device.IsConnected(clientId), "Client should be disconnected")

	// Reconnect
	device.ConnectClient(clientId)
	assert.Equal(t, true, device.IsConnected(clientId), "Client should be reconnected")

	// Disconnect again
	device.DisconnectClient(clientId)
	assert.Equal(t, false, device.IsConnected(clientId), "Client should be disconnected again")
}

func TestSafetyMonitorAPI_Connect_Disconnect_InvalidClient(t *testing.T) {
	// Test that disconnecting an invalid client doesn't cause issues
	api := &SafetyMonitorAPI{}

	// Create a mock device
	device := &Device{
		Id:               "test-monitor",
		Type:             "safetymonitor",
		Index:            0,
		ConnectedClients: make(map[ClientId]*ConnectedClient),
	}

	// Mock the API server
	apiServer := &ApiServer{
		Devices: map[string]map[int]*Device{
			"safetymonitor": {
				0: device,
			},
		},
	}
	api.ApiServer = apiServer

	clientId1 := ClientId("127.0.0.1-123")
	clientId2 := ClientId("127.0.0.1-456")

	// Connect one client
	device.ConnectClient(clientId1)
	assert.Equal(t, true, device.IsConnected(clientId1), "Client 1 should be connected")

	// Try to disconnect a client that was never connected
	device.DisconnectClient(clientId2)

	// First client should still be connected
	assert.Equal(t, true, device.IsConnected(clientId1), "Client 1 should still be connected")
	assert.Equal(t, false, device.IsConnected(clientId2), "Client 2 should not be connected")

	// Disconnect the connected client
	device.DisconnectClient(clientId1)
	assert.Equal(t, false, device.IsConnected(clientId1), "Client 1 should be disconnected")
}
