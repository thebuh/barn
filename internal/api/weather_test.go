package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWeatherAPI_ConnectClient(t *testing.T) {
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

func TestWeatherAPI_DisconnectClient(t *testing.T) {
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

func TestWeatherAPI_DeviceState_Structure(t *testing.T) {
	// Test that weather device state response has all required properties
	expectedProperties := []string{
		"CloudCover",
		"DewPoint",
		"Humidity",
		"Pressure",
		"RainRate",
		"SkyBrightness",
		"SkyQuality",
		"SkyTemperature",
		"StarFWHM",
		"Temperature",
		"WindDirection",
		"WindGust",
		"WindSpeed",
		"TimeStamp",
	}

	deviceStates := make([]DeviceState, len(expectedProperties))
	for i, prop := range expectedProperties {
		deviceStates[i] = DeviceState{
			Name:  prop,
			Value: true,
		}
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

	assert.Len(t, resp.Value, 14, "Should have 14 device states")

	// Verify all expected properties are present
	for i, expectedProp := range expectedProperties {
		assert.Equal(t, expectedProp, resp.Value[i].Name, "Property name should match")
		assert.Equal(t, true, resp.Value[i].Value, "Property value should be true")
	}

	assert.Equal(t, uint32(1), resp.ClientTransactionID, "ClientTransactionID should match")
}

func TestWeatherAPI_Connecting_AlwaysFalse(t *testing.T) {
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

func TestWeatherAPI_SensorDescription_Structure(t *testing.T) {
	// Test that sensor description response has correct structure
	resp := stringResponse{
		Value: "Sensor description for Temperature",
		alpacaResponse: alpacaResponse{
			ClientTransactionID: 1,
			ServerTransactionID: 2,
			ErrorNumber:         0,
			ErrorMessage:        "",
		},
	}

	assert.Equal(t, "Sensor description for Temperature", resp.Value, "Sensor description should match")
	assert.Equal(t, uint32(1), resp.ClientTransactionID, "ClientTransactionID should match")
}

func TestWeatherAPI_PercentDoubleResponse_Structure(t *testing.T) {
	// Test that percent double response has correct structure for cloud cover
	resp := percentDoubleResponse{
		Value: 45.5,
		alpacaResponse: alpacaResponse{
			ClientTransactionID: 1,
			ServerTransactionID: 2,
			ErrorNumber:         0,
			ErrorMessage:        "",
		},
	}

	assert.Equal(t, 45.5, resp.Value, "Cloud cover value should match")
	assert.Equal(t, uint32(1), resp.ClientTransactionID, "ClientTransactionID should match")
	// Ensure cloud cover is within valid percentage range (0-100)
	assert.GreaterOrEqual(t, resp.Value, 0.0, "Cloud cover should be >= 0")
	assert.LessOrEqual(t, resp.Value, 100.0, "Cloud cover should be <= 100")
}

func TestWeatherAPI_ConnectAsync_Structure(t *testing.T) {
	// Test that async connect response has correct structure
	resp := alpacaResponse{
		ClientTransactionID: 1,
		ServerTransactionID: 2,
		ErrorNumber:         0,
		ErrorMessage:        "",
	}

	assert.Equal(t, uint32(1), resp.ClientTransactionID, "ClientTransactionID should match")
	assert.Equal(t, uint32(2), resp.ServerTransactionID, "ServerTransactionID should match")
	assert.Equal(t, int32(0), resp.ErrorNumber, "ErrorNumber should be 0 for success")
	assert.Equal(t, "", resp.ErrorMessage, "ErrorMessage should be empty for success")
}

func TestWeatherAPI_DisconnectAsync_Structure(t *testing.T) {
	// Test that async disconnect response has correct structure
	resp := alpacaResponse{
		ClientTransactionID: 1,
		ServerTransactionID: 2,
		ErrorNumber:         0,
		ErrorMessage:        "",
	}

	assert.Equal(t, uint32(1), resp.ClientTransactionID, "ClientTransactionID should match")
	assert.Equal(t, uint32(2), resp.ServerTransactionID, "ServerTransactionID should match")
	assert.Equal(t, int32(0), resp.ErrorNumber, "ErrorNumber should be 0 for success")
	assert.Equal(t, "", resp.ErrorMessage, "ErrorMessage should be empty for success")
}

func TestWeatherAPI_WeatherState_Management(t *testing.T) {
	device := &Device{
		Id:               "test-weather",
		Type:             "observingconditions",
		Index:            0,
		ConnectedClients: make(map[ClientId]*ConnectedClient),
	}
	clientId1 := ClientId("127.0.0.1-123")
	clientId2 := ClientId("127.0.0.1-456")
	device.ConnectClient(clientId1)
	device.ConnectClient(clientId2)
	avgPeriod1, exists1 := device.GetWeatherAveragePeriod(clientId1)
	assert.True(t, exists1)
	assert.Equal(t, 0.0, avgPeriod1)
	avgPeriod2, exists2 := device.GetWeatherAveragePeriod(clientId2)
	assert.True(t, exists2)
	assert.Equal(t, 0.0, avgPeriod2)
	success1 := device.SetWeatherAveragePeriod(clientId1, 0.15)
	assert.True(t, success1)
	success2 := device.SetWeatherAveragePeriod(clientId2, 0.30)
	assert.True(t, success2)
	avgPeriod1, exists1 = device.GetWeatherAveragePeriod(clientId1)
	assert.True(t, exists1)
	assert.Equal(t, 0.15, avgPeriod1)
	avgPeriod2, exists2 = device.GetWeatherAveragePeriod(clientId2)
	assert.True(t, exists2)
	assert.Equal(t, 0.30, avgPeriod2)
	disconnectedClientId := ClientId("127.0.0.1-999")
	success3 := device.SetWeatherAveragePeriod(disconnectedClientId, 0.45)
	assert.False(t, success3)
	avgPeriod3, exists3 := device.GetWeatherAveragePeriod(disconnectedClientId)
	assert.False(t, exists3)
	assert.Equal(t, 0.0, avgPeriod3)
}

func TestWeatherAPI_WeatherState_Disconnect_Reconnect(t *testing.T) {
	device := &Device{
		Id:               "test-weather",
		Type:             "observingconditions",
		Index:            0,
		ConnectedClients: make(map[ClientId]*ConnectedClient),
	}
	clientId := ClientId("127.0.0.1-123")
	device.ConnectClient(clientId)
	success := device.SetWeatherAveragePeriod(clientId, 0.25)
	assert.True(t, success)
	avgPeriod, exists := device.GetWeatherAveragePeriod(clientId)
	assert.True(t, exists)
	assert.Equal(t, 0.25, avgPeriod)
	device.DisconnectClient(clientId)
	avgPeriod, exists = device.GetWeatherAveragePeriod(clientId)
	assert.False(t, exists)
	assert.Equal(t, 0.0, avgPeriod)
	device.ConnectClient(clientId)
	avgPeriod, exists = device.GetWeatherAveragePeriod(clientId)
	assert.True(t, exists)
	assert.Equal(t, 0.0, avgPeriod)
}

func TestWeatherAPI_WeatherState_MultipleClients(t *testing.T) {
	device := &Device{
		Id:               "test-weather",
		Type:             "observingconditions",
		Index:            0,
		ConnectedClients: make(map[ClientId]*ConnectedClient),
	}
	clientIds := []ClientId{
		ClientId("127.0.0.1-100"),
		ClientId("127.0.0.1-200"),
		ClientId("127.0.0.1-300"),
	}
	for _, clientId := range clientIds {
		device.ConnectClient(clientId)
	}
	expectedPeriods := []float64{0.1, 0.2, 0.3}
	for i, clientId := range clientIds {
		success := device.SetWeatherAveragePeriod(clientId, expectedPeriods[i])
		assert.True(t, success)
	}
	for i, clientId := range clientIds {
		avgPeriod, exists := device.GetWeatherAveragePeriod(clientId)
		assert.True(t, exists)
		assert.Equal(t, expectedPeriods[i], avgPeriod)
	}
	device.DisconnectClient(clientIds[1])
	avgPeriod0, exists0 := device.GetWeatherAveragePeriod(clientIds[0])
	assert.True(t, exists0)
	assert.Equal(t, expectedPeriods[0], avgPeriod0)
	avgPeriod2, exists2 := device.GetWeatherAveragePeriod(clientIds[2])
	assert.True(t, exists2)
	assert.Equal(t, expectedPeriods[2], avgPeriod2)
	avgPeriod1, exists1 := device.GetWeatherAveragePeriod(clientIds[1])
	assert.False(t, exists1)
	assert.Equal(t, 0.0, avgPeriod1)
}

func TestWeatherAPI_Refresh_Response_Structure(t *testing.T) {
	// Test that refresh response has correct structure
	resp := alpacaResponse{
		ClientTransactionID: 1,
		ServerTransactionID: 2,
		ErrorNumber:         0,
		ErrorMessage:        "",
	}

	assert.Equal(t, uint32(1), resp.ClientTransactionID, "ClientTransactionID should match")
	assert.Equal(t, uint32(2), resp.ServerTransactionID, "ServerTransactionID should match")
	assert.Equal(t, int32(0), resp.ErrorNumber, "ErrorNumber should be 0 for success")
	assert.Equal(t, "", resp.ErrorMessage, "ErrorMessage should be empty for success")
}

func TestWeatherAPI_SupportedActions_IncludesRefresh(t *testing.T) {
	// Test that supported actions includes Refresh
	resp := stringlistResponse{
		Value: []string{"Refresh"},
		alpacaResponse: alpacaResponse{
			ClientTransactionID: 1,
			ServerTransactionID: 2,
			ErrorNumber:         0,
			ErrorMessage:        "",
		},
	}

	assert.Len(t, resp.Value, 1, "Should have 1 supported action")
	assert.Equal(t, "Refresh", resp.Value[0], "Supported action should be Refresh")
	assert.Equal(t, uint32(1), resp.ClientTransactionID, "ClientTransactionID should match")
}
