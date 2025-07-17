package api

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/thebuh/barn/internal/weather"
)

// WeatherAPI handles all weather/observing conditions related API endpoints
type WeatherAPI struct {
	*ApiServer
}

// NewWeatherAPI creates a new weather API handler
func NewWeatherAPI(apiServer *ApiServer) *WeatherAPI {
	return &WeatherAPI{
		ApiServer: apiServer,
	}
}

// ConfigureRoutes sets up all weather API routes
func (w *WeatherAPI) ConfigureRoutes(router *gin.Engine) {
	// Platform 7 asynchronous connect/disconnect endpoints
	observingConditions := router.Group("/api/v1/observingconditions/:device_id")
	// Apply validation middleware to all routes
	observingConditions.Use(alpacaValidationMiddleware())
	observingConditions.Use(deviceValidationMiddleware("observingconditions"))
	observingConditions.Use(alpacaResponseMiddleware())
	{
		observingConditions.PUT("/connect", w.handleConnect)
		observingConditions.PUT("/disconnect", w.handleDisconnect)
		observingConditions.PUT("/connected", w.handleConnected)
		observingConditions.PUT("/action", w.handleAction)
		observingConditions.PUT("/refresh", w.handleRefresh)
		observingConditions.PUT("/averageperiod", w.handleAveragePeriod)

		// Individual GET routes for each observing conditions property
		observingConditions.GET("/connected", w.handleConnectedGet)
		observingConditions.GET("/connecting", w.handleConnecting)
		observingConditions.GET("/name", w.handleName)
		observingConditions.GET("/description", w.handleDescription)
		observingConditions.GET("/driverinfo", w.handleDriverInfo)
		observingConditions.GET("/driverversion", w.handleDriverVersion)
		observingConditions.GET("/supportedactions", w.handleSupportedActions)
		observingConditions.GET("/interfaceversion", w.handleInterfaceVersion)
		observingConditions.GET("/averageperiod", w.handleAveragePeriodGet)
		observingConditions.GET("/cloudcover", w.handleCloudCover)
		observingConditions.GET("/dewpoint", w.handleDewPoint)
		observingConditions.GET("/humidity", w.handleHumidity)
		observingConditions.GET("/pressure", w.handlePressure)
		observingConditions.GET("/rainrate", w.handleRainRate)
		observingConditions.GET("/skybrightness", w.handleSkyBrightness)
		observingConditions.GET("/skyquality", w.handleSkyQuality)
		observingConditions.GET("/skytemperature", w.handleSkyTemperature)
		observingConditions.GET("/starfwhm", w.handleStarFWHM)
		observingConditions.GET("/temperature", w.handleTemperature)
		observingConditions.GET("/winddirection", w.handleWindDirection)
		observingConditions.GET("/windgust", w.handleWindGust)
		observingConditions.GET("/windspeed", w.handleWindSpeed)
		observingConditions.GET("/sensordescription", w.handleSensorDescription)
		observingConditions.GET("/timesincelastupdate", w.handleTimeSinceLastUpdate)
		observingConditions.GET("/devicestate", w.handleDeviceState)
	}
}

// isRequestConnected checks if the current request is from a connected client
func (w *WeatherAPI) isRequestConnected(c *gin.Context) bool {
	deviceId, _ := strconv.Atoi(c.Param("device_id"))
	device := w.Devices["observingconditions"][deviceId]
	return device.IsConnected(getFullClientId(c))
}

// handleConnectedGet handles GET requests for connected property
func (w *WeatherAPI) handleConnectedGet(c *gin.Context) {
	result := w.isRequestConnected(c)
	resp := boolResponse{
		Value: result,
	}
	w.prepareAlpacaResponse(c, &resp.alpacaResponse)
	c.IndentedJSON(http.StatusOK, resp)
}

// handleConnecting handles GET requests for connecting property
func (w *WeatherAPI) handleConnecting(c *gin.Context) {
	resp := boolResponse{
		Value: false,
	}
	w.prepareAlpacaResponse(c, &resp.alpacaResponse)
	c.IndentedJSON(http.StatusOK, resp)
}

// handleName handles GET requests for name property
func (w *WeatherAPI) handleName(c *gin.Context) {
	deviceId, _ := strconv.Atoi(c.Param("device_id"))
	device, err := w.Barn.GetWeatherByIndex(deviceId)
	if err != nil {
		c.String(400, "Device not found")
		return
	}

	resp := stringResponse{
		Value: device.GetName(),
	}
	w.prepareAlpacaResponse(c, &resp.alpacaResponse)
	c.IndentedJSON(http.StatusOK, resp)
}

// handleDescription handles GET requests for description property
func (w *WeatherAPI) handleDescription(c *gin.Context) {
	deviceId, _ := strconv.Atoi(c.Param("device_id"))
	device, err := w.Barn.GetWeatherByIndex(deviceId)
	if err != nil {
		c.String(400, "Device not found")
		return
	}

	resp := stringResponse{
		Value: device.GetDescription(),
	}
	w.prepareAlpacaResponse(c, &resp.alpacaResponse)
	c.IndentedJSON(http.StatusOK, resp)
}

// handleDriverInfo handles GET requests for driverinfo property
func (w *WeatherAPI) handleDriverInfo(c *gin.Context) {
	resp := stringResponse{
		Value: "Alpaca Barn observing conditions",
	}
	w.prepareAlpacaResponse(c, &resp.alpacaResponse)
	c.IndentedJSON(http.StatusOK, resp)
}

// handleDriverVersion handles GET requests for driverversion property
func (w *WeatherAPI) handleDriverVersion(c *gin.Context) {
	bi, _ := debug.ReadBuildInfo()
	resp := stringResponse{
		Value: bi.Main.Version,
	}
	w.prepareAlpacaResponse(c, &resp.alpacaResponse)
	c.IndentedJSON(http.StatusOK, resp)
}

// handleSupportedActions handles GET requests for supportedactions property
func (w *WeatherAPI) handleSupportedActions(c *gin.Context) {
	resp := stringlistResponse{
		Value: []string{"Refresh"},
	}
	w.prepareAlpacaResponse(c, &resp.alpacaResponse)
	c.IndentedJSON(http.StatusOK, resp)
}

// handleInterfaceVersion handles GET requests for interfaceversion property
func (w *WeatherAPI) handleInterfaceVersion(c *gin.Context) {
	resp := int32Response{
		Value: 2,
	}
	w.prepareAlpacaResponse(c, &resp.alpacaResponse)
	c.IndentedJSON(http.StatusOK, resp)
}

// handleAveragePeriodGet handles GET requests for averageperiod property
func (w *WeatherAPI) handleAveragePeriodGet(c *gin.Context) {
	deviceId, _ := strconv.Atoi(c.Param("device_id"))
	device, err := w.Barn.GetWeatherByIndex(deviceId)
	if err != nil {
		c.String(400, "Device not found")
		return
	}

	if !w.isRequestConnected(c) {
		c.String(400, "Not connected")
		return
	}

	// Get client-specific average period
	clientId := getFullClientId(c)
	dev := w.Devices["observingconditions"][deviceId]
	averagePeriod, exists := dev.GetWeatherAveragePeriod(clientId)

	if !exists {
		// Fallback to device default if no client-specific value
		averagePeriod = device.GetAveragePeriod()
	}

	resp := float64Response{
		Value: averagePeriod,
	}
	w.prepareAlpacaResponse(c, &resp.alpacaResponse)
	c.IndentedJSON(http.StatusOK, resp)
}

// handleCloudCover handles GET requests for cloudcover property
func (w *WeatherAPI) handleCloudCover(c *gin.Context) {
	deviceId, _ := strconv.Atoi(c.Param("device_id"))
	device, err := w.Barn.GetWeatherByIndex(deviceId)
	if err != nil {
		c.String(400, "Device not found")
		return
	}

	if !w.isRequestConnected(c) {
		c.String(400, "Not connected")
		return
	}

	resp := percentDoubleResponse{
		Value: device.GetCloudCover(),
	}
	w.prepareAlpacaResponse(c, &resp.alpacaResponse)
	c.IndentedJSON(http.StatusOK, resp)
}

// handleDewPoint handles GET requests for dewpoint property
func (w *WeatherAPI) handleDewPoint(c *gin.Context) {
	deviceId, _ := strconv.Atoi(c.Param("device_id"))
	device, err := w.Barn.GetWeatherByIndex(deviceId)
	if err != nil {
		c.String(400, "Device not found")
		return
	}

	if !w.isRequestConnected(c) {
		c.String(400, "Not connected")
		return
	}

	resp := float64Response{
		Value: device.GetDewPoint(),
	}
	w.prepareAlpacaResponse(c, &resp.alpacaResponse)
	c.IndentedJSON(http.StatusOK, resp)
}

// handleHumidity handles GET requests for humidity property
func (w *WeatherAPI) handleHumidity(c *gin.Context) {
	deviceId, _ := strconv.Atoi(c.Param("device_id"))
	device, err := w.Barn.GetWeatherByIndex(deviceId)
	if err != nil {
		c.String(400, "Device not found")
		return
	}

	if !w.isRequestConnected(c) {
		c.String(400, "Not connected")
		return
	}

	resp := float64Response{
		Value: device.GetHumidity(),
	}
	w.prepareAlpacaResponse(c, &resp.alpacaResponse)
	c.IndentedJSON(http.StatusOK, resp)
}

// handlePressure handles GET requests for pressure property
func (w *WeatherAPI) handlePressure(c *gin.Context) {
	deviceId, _ := strconv.Atoi(c.Param("device_id"))
	device, err := w.Barn.GetWeatherByIndex(deviceId)
	if err != nil {
		c.String(400, "Device not found")
		return
	}

	if !w.isRequestConnected(c) {
		c.String(400, "Not connected")
		return
	}

	resp := float64Response{
		Value: device.GetPressure(),
	}
	w.prepareAlpacaResponse(c, &resp.alpacaResponse)
	c.IndentedJSON(http.StatusOK, resp)
}

// handleRainRate handles GET requests for rainrate property
func (w *WeatherAPI) handleRainRate(c *gin.Context) {
	deviceId, _ := strconv.Atoi(c.Param("device_id"))
	device, err := w.Barn.GetWeatherByIndex(deviceId)
	if err != nil {
		c.String(400, "Device not found")
		return
	}

	if !w.isRequestConnected(c) {
		c.String(400, "Not connected")
		return
	}

	resp := float64Response{
		Value: device.GetRainRate(),
	}
	w.prepareAlpacaResponse(c, &resp.alpacaResponse)
	c.IndentedJSON(http.StatusOK, resp)
}

// handleSkyBrightness handles GET requests for skybrightness property
func (w *WeatherAPI) handleSkyBrightness(c *gin.Context) {
	deviceId, _ := strconv.Atoi(c.Param("device_id"))
	device, err := w.Barn.GetWeatherByIndex(deviceId)
	if err != nil {
		c.String(400, "Device not found")
		return
	}

	if !w.isRequestConnected(c) {
		c.String(400, "Not connected")
		return
	}

	// Check if sensor is available
	if !weather.IsSensorAvailable(weather.SensorSkyBrightness) {
		resp := alpacaResponse{
			ClientTransactionID: 0,
			ServerTransactionID: 0,
			ErrorNumber:         0x400, // NotImplemented
			ErrorMessage:        "Sensor SkyBrightness is not supported by this device",
		}
		w.prepareAlpacaResponse(c, &resp)
		c.IndentedJSON(http.StatusOK, resp)
		return
	}

	resp := float64Response{
		Value: device.GetSkyBrightness(),
	}
	w.prepareAlpacaResponse(c, &resp.alpacaResponse)
	c.IndentedJSON(http.StatusOK, resp)
}

// handleSkyQuality handles GET requests for skyquality property
func (w *WeatherAPI) handleSkyQuality(c *gin.Context) {
	deviceId, _ := strconv.Atoi(c.Param("device_id"))
	device, err := w.Barn.GetWeatherByIndex(deviceId)
	if err != nil {
		c.String(400, "Device not found")
		return
	}

	if !w.isRequestConnected(c) {
		c.String(400, "Not connected")
		return
	}

	// Check if sensor is available
	if !weather.IsSensorAvailable(weather.SensorSkyQuality) {
		resp := alpacaResponse{
			ClientTransactionID: 0,
			ServerTransactionID: 0,
			ErrorNumber:         0x400, // NotImplemented
			ErrorMessage:        "Sensor SkyQuality is not supported by this device",
		}
		w.prepareAlpacaResponse(c, &resp)
		c.IndentedJSON(http.StatusOK, resp)
		return
	}

	resp := float64Response{
		Value: device.GetSkyQuality(),
	}
	w.prepareAlpacaResponse(c, &resp.alpacaResponse)
	c.IndentedJSON(http.StatusOK, resp)
}

// handleSkyTemperature handles GET requests for skytemperature property
func (w *WeatherAPI) handleSkyTemperature(c *gin.Context) {
	deviceId, _ := strconv.Atoi(c.Param("device_id"))
	device, err := w.Barn.GetWeatherByIndex(deviceId)
	if err != nil {
		c.String(400, "Device not found")
		return
	}

	if !w.isRequestConnected(c) {
		c.String(400, "Not connected")
		return
	}

	// Check if sensor is available
	if !weather.IsSensorAvailable(weather.SensorSkyTemperature) {
		resp := alpacaResponse{
			ClientTransactionID: 0,
			ServerTransactionID: 0,
			ErrorNumber:         0x400, // NotImplemented
			ErrorMessage:        "Sensor SkyTemperature is not supported by this device",
		}
		w.prepareAlpacaResponse(c, &resp)
		c.IndentedJSON(http.StatusOK, resp)
		return
	}

	resp := float64Response{
		Value: device.GetSkyTemperature(),
	}
	w.prepareAlpacaResponse(c, &resp.alpacaResponse)
	c.IndentedJSON(http.StatusOK, resp)
}

// handleStarFWHM handles GET requests for starfwhm property
func (w *WeatherAPI) handleStarFWHM(c *gin.Context) {
	deviceId, _ := strconv.Atoi(c.Param("device_id"))
	device, err := w.Barn.GetWeatherByIndex(deviceId)
	if err != nil {
		c.String(400, "Device not found")
		return
	}

	if !w.isRequestConnected(c) {
		c.String(400, "Not connected")
		return
	}

	// Check if sensor is available
	if !weather.IsSensorAvailable(weather.SensorStarFWHM) {
		resp := alpacaResponse{
			ClientTransactionID: 0,
			ServerTransactionID: 0,
			ErrorNumber:         0x400, // NotImplemented
			ErrorMessage:        "Sensor StarFWHM is not supported by this device",
		}
		w.prepareAlpacaResponse(c, &resp)
		c.IndentedJSON(http.StatusOK, resp)
		return
	}

	resp := float64Response{
		Value: device.GetStarFWHM(),
	}
	w.prepareAlpacaResponse(c, &resp.alpacaResponse)
	c.IndentedJSON(http.StatusOK, resp)
}

// handleTemperature handles GET requests for temperature property
func (w *WeatherAPI) handleTemperature(c *gin.Context) {
	deviceId, _ := strconv.Atoi(c.Param("device_id"))
	device, err := w.Barn.GetWeatherByIndex(deviceId)
	if err != nil {
		c.String(400, "Device not found")
		return
	}

	if !w.isRequestConnected(c) {
		c.String(400, "Not connected")
		return
	}

	// Check if sensor is available
	if !weather.IsSensorAvailable(weather.SensorTemperature) {
		resp := alpacaResponse{
			ClientTransactionID: 0,
			ServerTransactionID: 0,
			ErrorNumber:         0x400, // NotImplemented
			ErrorMessage:        "Sensor Temperature is not supported by this device",
		}
		w.prepareAlpacaResponse(c, &resp)
		c.IndentedJSON(http.StatusOK, resp)
		return
	}

	resp := float64Response{
		Value: device.GetTemperature(),
	}
	w.prepareAlpacaResponse(c, &resp.alpacaResponse)
	c.IndentedJSON(http.StatusOK, resp)
}

// handleWindDirection handles GET requests for winddirection property
func (w *WeatherAPI) handleWindDirection(c *gin.Context) {
	deviceId, _ := strconv.Atoi(c.Param("device_id"))
	device, err := w.Barn.GetWeatherByIndex(deviceId)
	if err != nil {
		c.String(400, "Device not found")
		return
	}

	if !w.isRequestConnected(c) {
		c.String(400, "Not connected")
		return
	}

	resp := float64Response{
		Value: device.GetWindDirection(),
	}
	w.prepareAlpacaResponse(c, &resp.alpacaResponse)
	c.IndentedJSON(http.StatusOK, resp)
}

// handleWindGust handles GET requests for windgust property
func (w *WeatherAPI) handleWindGust(c *gin.Context) {
	deviceId, _ := strconv.Atoi(c.Param("device_id"))
	device, err := w.Barn.GetWeatherByIndex(deviceId)
	if err != nil {
		c.String(400, "Device not found")
		return
	}

	if !w.isRequestConnected(c) {
		c.String(400, "Not connected")
		return
	}

	resp := float64Response{
		Value: device.GetWindGust(),
	}
	w.prepareAlpacaResponse(c, &resp.alpacaResponse)
	c.IndentedJSON(http.StatusOK, resp)
}

// handleWindSpeed handles GET requests for windspeed property
func (w *WeatherAPI) handleWindSpeed(c *gin.Context) {
	deviceId, _ := strconv.Atoi(c.Param("device_id"))
	device, err := w.Barn.GetWeatherByIndex(deviceId)
	if err != nil {
		c.String(400, "Device not found")
		return
	}

	if !w.isRequestConnected(c) {
		c.String(400, "Not connected")
		return
	}

	resp := float64Response{
		Value: device.GetWindSpeed(),
	}
	w.prepareAlpacaResponse(c, &resp.alpacaResponse)
	c.IndentedJSON(http.StatusOK, resp)
}

// handleTimeSinceLastUpdate handles GET requests for timesincelastupdate property
func (w *WeatherAPI) handleTimeSinceLastUpdate(c *gin.Context) {
	deviceId, _ := strconv.Atoi(c.Param("device_id"))
	device, err := w.Barn.GetWeatherByIndex(deviceId)
	if err != nil {
		c.String(400, "Device not found")
		return
	}

	if !w.isRequestConnected(c) {
		c.String(400, "Not connected")
		return
	}

	// Get sensor name from query parameter as required by spec
	sensorName := getQuery(c, "SensorName")

	// If sensor name is provided, check if it's available
	if sensorName != "" && !weather.IsSensorAvailable(sensorName) {
		resp := alpacaResponse{
			ClientTransactionID: 0,
			ServerTransactionID: 0,
			ErrorNumber:         0x400, // NotImplemented
			ErrorMessage:        fmt.Sprintf("Sensor '%s' is not supported by this device", sensorName),
		}
		w.prepareAlpacaResponse(c, &resp)
		c.IndentedJSON(http.StatusOK, resp)
		return
	}

	resp := float64Response{
		Value: device.GetTimeSinceLastUpdate(),
	}
	w.prepareAlpacaResponse(c, &resp.alpacaResponse)
	c.IndentedJSON(http.StatusOK, resp)
}

// handleSensorDescription handles GET requests for sensordescription property
func (w *WeatherAPI) handleSensorDescription(c *gin.Context) {
	deviceId, _ := strconv.Atoi(c.Param("device_id"))
	_, err := w.Barn.GetWeatherByIndex(deviceId)
	if err != nil {
		c.String(400, "Device not found")
		return
	}

	if !w.isRequestConnected(c) {
		c.String(400, "Not connected")
		return
	}

	// Get sensor name from query parameter as required by spec
	sensorName := getQuery(c, "SensorName")
	if sensorName == "" {
		c.String(400, "SensorName parameter is required")
		return
	}

	// Check if the sensor is available (not just valid)
	if !weather.IsSensorAvailable(sensorName) {
		resp := alpacaResponse{
			ClientTransactionID: 0,
			ServerTransactionID: 0,
			ErrorNumber:         0x400, // NotImplemented
			ErrorMessage:        fmt.Sprintf("Sensor '%s' is not supported by this device", sensorName),
		}
		w.prepareAlpacaResponse(c, &resp)
		c.IndentedJSON(http.StatusOK, resp)
		return
	}

	// Get sensor description from weather package
	description, exists := weather.GetSensorDescription(sensorName)
	if !exists {
		c.String(400, fmt.Sprintf("Sensor '%s' not found", sensorName))
		return
	}

	resp := stringResponse{
		Value: description,
	}
	w.prepareAlpacaResponse(c, &resp.alpacaResponse)
	c.IndentedJSON(http.StatusOK, resp)
}

// handleDeviceState handles GET requests for devicestate property
func (w *WeatherAPI) handleDeviceState(c *gin.Context) {
	deviceId, _ := strconv.Atoi(c.Param("device_id"))
	_, err := w.Barn.GetWeatherByIndex(deviceId)
	if err != nil {
		c.String(400, "Device not found")
		return
	}

	if !w.isRequestConnected(c) {
		c.String(400, "Not connected")
		return
	}

	// Create device state array with only available sensors
	deviceStates := make([]DeviceState, 0)

	// Add only available sensors
	for sensorName, available := range weather.AvailableSensors {
		if available {
			deviceStates = append(deviceStates, DeviceState{
				Name:  sensorName,
				Value: true,
			})
		}
	}

	// Add TimeStamp sensor
	deviceStates = append(deviceStates, DeviceState{
		Name:  weather.SensorTimeStamp,
		Value: true,
	})

	resp := deviceStateResponse{
		Value: deviceStates,
	}
	w.prepareAlpacaResponse(c, &resp.alpacaResponse)
	c.IndentedJSON(http.StatusOK, resp)
}

// handleConnected handles PUT requests to connect/disconnect weather devices
func (w *WeatherAPI) handleConnected(c *gin.Context) {
	deviceId, _ := strconv.Atoi(c.Param("device_id"))
	connected := c.PostForm("Connected")
	_, err := w.Barn.GetWeatherByIndex(deviceId)
	if err != nil {
		c.String(400, "Device not found")
		return
	}
	dev := w.Devices["observingconditions"][deviceId]
	id := getFullClientId(c)

	if strings.ToLower(connected) == "true" { // Connect
		dev.ConnectClient(id)
	} else if strings.ToLower(connected) == "false" { //Disconnect
		dev.DisconnectClient(id)
	} else {
		c.String(400, "Invalid request")
		return
	}
	resp := alpacaResponse{
		ClientTransactionID: 0,
		ServerTransactionID: 0,
		ErrorNumber:         0,
		ErrorMessage:        "",
	}
	w.prepareAlpacaResponse(c, &resp)
	c.IndentedJSON(http.StatusOK, resp)
}

// handleAction handles PUT requests for weather device actions
func (w *WeatherAPI) handleAction(c *gin.Context) {
	deviceId, _ := strconv.Atoi(c.Param("device_id"))
	device, err := w.Barn.GetWeatherByIndex(deviceId)
	if err != nil {
		c.String(400, "Device not found")
		return
	}

	action := c.PostForm("Action")
	if !w.isRequestConnected(c) {
		c.String(400, "Not connected")
		return
	}

	switch action {
	case "Refresh":
		w.handleRefreshAction(deviceId, device, c)
	default:
		c.String(400, "The device did not understand which operation was being requested or insufficient information was given to complete the operation.")
		return
	}
}

// handleRefresh handles PUT requests specifically for refreshing weather sensor values
func (w *WeatherAPI) handleRefresh(c *gin.Context) {
	deviceId, _ := strconv.Atoi(c.Param("device_id"))
	device, err := w.Barn.GetWeatherByIndex(deviceId)
	if err != nil {
		c.String(400, "Device not found")
		return
	}

	if !w.isRequestConnected(c) {
		c.String(400, "Not connected")
		return
	}

	w.handleRefreshAction(deviceId, device, c)
}

// handleRefreshAction is a shared method that handles the refresh action logic
func (w *WeatherAPI) handleRefreshAction(deviceId int, device weather.ObservingConditions, c *gin.Context) {
	err := device.Refresh()
	if err != nil {
		log.WithFields(log.Fields{
			"deviceid": deviceId,
			"weather":  device.GetName(),
			"error":    err,
		}).Error(fmt.Sprintf("[BARN] Weather [%s]. Failed to refresh: %v", device.GetName(), err))
		c.String(500, "Failed to refresh weather data")
		return
	}
	resp := alpacaResponse{
		ClientTransactionID: 0,
		ServerTransactionID: 0,
		ErrorNumber:         0,
		ErrorMessage:        "",
	}
	w.prepareAlpacaResponse(c, &resp)
	c.IndentedJSON(http.StatusOK, resp)
}

// handleAveragePeriod handles PUT requests to set the average period
func (w *WeatherAPI) handleAveragePeriod(c *gin.Context) {
	deviceId, _ := strconv.Atoi(c.Param("device_id"))
	device, err := w.Barn.GetWeatherByIndex(deviceId)
	if err != nil {
		c.String(400, "Device not found")
		return
	}

	if !w.isRequestConnected(c) {
		c.String(400, "Not connected")
		return
	}

	// Parse the average period from the request body
	averagePeriodStr := c.PostForm("AveragePeriod")
	if averagePeriodStr == "" {
		c.String(400, "AveragePeriod parameter is required")
		return
	}

	averagePeriod, err := strconv.ParseFloat(averagePeriodStr, 64)
	if err != nil {
		c.String(400, "Invalid AveragePeriod value")
		return
	}

	// Store the average period in client-specific state
	clientId := getFullClientId(c)
	dev := w.Devices["observingconditions"][deviceId]
	success := dev.SetWeatherAveragePeriod(clientId, averagePeriod)

	if !success {
		log.WithFields(log.Fields{
			"deviceid": deviceId,
			"weather":  device.GetName(),
			"clientid": clientId,
		}).Error(fmt.Sprintf("[BARN] Weather [%s]. Failed to set client-specific average period for client %s", device.GetName(), clientId))
		c.String(500, "Failed to set average period")
		return
	}

	log.WithFields(log.Fields{
		"deviceid":      deviceId,
		"weather":       device.GetName(),
		"clientid":      clientId,
		"averageperiod": averagePeriod,
	}).Info(fmt.Sprintf("[BARN] Weather [%s]. Set client-specific average period to %f for client %s", device.GetName(), averagePeriod, clientId))

	resp := alpacaResponse{
		ClientTransactionID: 0,
		ServerTransactionID: 0,
		ErrorNumber:         0,
		ErrorMessage:        "",
	}
	w.prepareAlpacaResponse(c, &resp)
	c.IndentedJSON(http.StatusOK, resp)
}

// handleConnect handles PUT requests to start an asynchronous connect
func (w *WeatherAPI) handleConnect(c *gin.Context) {
	deviceId, _ := strconv.Atoi(c.Param("device_id"))
	_, err := w.Barn.GetWeatherByIndex(deviceId)
	if err != nil {
		c.String(400, "Device not found")
		return
	}

	// For this implementation, we'll make the connect synchronous
	// In a real implementation, this would start an asynchronous connection process
	dev := w.Devices["observingconditions"][deviceId]
	id := getFullClientId(c)
	dev.ConnectClient(id)

	resp := alpacaResponse{
		ClientTransactionID: 0,
		ServerTransactionID: 0,
		ErrorNumber:         0,
		ErrorMessage:        "",
	}
	w.prepareAlpacaResponse(c, &resp)
	c.IndentedJSON(http.StatusOK, resp)
}

// handleDisconnect handles PUT requests to start an asynchronous disconnect
func (w *WeatherAPI) handleDisconnect(c *gin.Context) {
	deviceId, _ := strconv.Atoi(c.Param("device_id"))
	_, err := w.Barn.GetWeatherByIndex(deviceId)
	if err != nil {
		c.String(400, "Device not found")
		return
	}

	// For this implementation, we'll make the disconnect synchronous
	// In a real implementation, this would start an asynchronous disconnection process
	dev := w.Devices["observingconditions"][deviceId]
	id := getFullClientId(c)
	dev.DisconnectClient(id)

	resp := alpacaResponse{
		ClientTransactionID: 0,
		ServerTransactionID: 0,
		ErrorNumber:         0,
		ErrorMessage:        "",
	}
	w.prepareAlpacaResponse(c, &resp)
	c.IndentedJSON(http.StatusOK, resp)
}
