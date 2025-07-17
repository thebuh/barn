package api

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// SafetyMonitorAPI handles all safety monitor related API endpoints
type SafetyMonitorAPI struct {
	*ApiServer
}

// NewSafetyMonitorAPI creates a new safety monitor API handler
func NewSafetyMonitorAPI(apiServer *ApiServer) *SafetyMonitorAPI {
	return &SafetyMonitorAPI{
		ApiServer: apiServer,
	}
}

// ConfigureRoutes sets up all safety monitor API routes
func (sm *SafetyMonitorAPI) ConfigureRoutes(router *gin.Engine) {
	safetyMonitor := router.Group("/api/v1/safetymonitor/:device_id")
	// Apply validation middleware to all routes
	safetyMonitor.Use(alpacaValidationMiddleware())
	safetyMonitor.Use(deviceValidationMiddleware("safetymonitor"))
	safetyMonitor.Use(alpacaResponseMiddleware())
	{
		safetyMonitor.PUT("/connected", sm.handleSafetyMonitorConnect)
		safetyMonitor.PUT("/action", sm.handleSafetyMonitorAction)

		// Platform 7 asynchronous connect/disconnect endpoints
		safetyMonitor.PUT("/connect", sm.handleConnect)
		safetyMonitor.PUT("/disconnect", sm.handleDisconnect)

		// Individual GET routes for each action
		safetyMonitor.GET("/issafe", sm.handleIsSafe)
		safetyMonitor.GET("/connected", sm.handleConnected)
		safetyMonitor.GET("/connecting", sm.handleConnecting)
		safetyMonitor.GET("/name", sm.handleName)
		safetyMonitor.GET("/description", sm.handleDescription)
		safetyMonitor.GET("/driverinfo", sm.handleDriverInfo)
		safetyMonitor.GET("/driverversion", sm.handleDriverVersion)
		safetyMonitor.GET("/supportedactions", sm.handleSupportedActions)
		safetyMonitor.GET("/interfaceversion", sm.handleInterfaceVersion)
		safetyMonitor.GET("/devicestate", sm.handleDeviceState)
	}
}

// isRequestConnected checks if the current request is from a connected client
func (sm *SafetyMonitorAPI) isRequestConnected(c *gin.Context) bool {
	validationCtx := GetValidationContext(c)
	if validationCtx == nil {
		return false
	}

	device, exists := c.Get("device")
	if !exists {
		return false
	}

	d := device.(*Device)
	return d.IsConnected(validationCtx.FullClientID)
}

// handleIsSafe handles GET requests for safety monitor isSafe property
func (sm *SafetyMonitorAPI) handleIsSafe(c *gin.Context) {
	device, _ := c.Get("device")
	d := device.(*Device)

	val := false
	if sm.isRequestConnected(c) {
		// Get the actual device from barn
		deviceId := GetValidationContext(c).DeviceID
		monitor, err := sm.Barn.GetMonitorByIndex(deviceId)
		if err == nil {
			val = monitor.IsSafe()
		}
	}

	resp := boolResponse{
		Value: val,
	}

	// Use the middleware-prepared response function
	if prepareResp, exists := c.Get("prepareResponse"); exists {
		if prepareFunc, ok := prepareResp.(func(*alpacaResponse)); ok {
			prepareFunc(&resp.alpacaResponse)
		}
	}

	log.WithFields(log.Fields{
		"deviceid": GetValidationContext(c).DeviceID,
		"monitor":  d.Id,
		"state":    val,
	}).Info(fmt.Sprintf("[BARN] Api [%s]. Returning state: [%t]", d.Id, val))

	c.IndentedJSON(http.StatusOK, resp)
}

// handleConnected handles GET requests for safety monitor connected property
func (sm *SafetyMonitorAPI) handleConnected(c *gin.Context) {
	result := sm.isRequestConnected(c)
	resp := boolResponse{
		Value: result,
	}
	sm.prepareAlpacaResponse(c, &resp.alpacaResponse)
	c.IndentedJSON(http.StatusOK, resp)
}

// handleConnecting handles GET requests for safety monitor connecting property
func (sm *SafetyMonitorAPI) handleConnecting(c *gin.Context) {
	resp := boolResponse{
		Value: false,
	}
	sm.prepareAlpacaResponse(c, &resp.alpacaResponse)
	c.IndentedJSON(http.StatusOK, resp)
}

// handleName handles GET requests for safety monitor name property
func (sm *SafetyMonitorAPI) handleName(c *gin.Context) {
	deviceId, _ := strconv.Atoi(c.Param("device_id"))
	device, err := sm.Barn.GetMonitorByIndex(deviceId)
	if err != nil {
		c.String(400, "Device not found")
		return
	}

	resp := stringResponse{
		Value: device.GetName(),
	}
	sm.prepareAlpacaResponse(c, &resp.alpacaResponse)
	c.IndentedJSON(http.StatusOK, resp)
}

// handleDescription handles GET requests for safety monitor description property
func (sm *SafetyMonitorAPI) handleDescription(c *gin.Context) {
	deviceId, _ := strconv.Atoi(c.Param("device_id"))
	device, err := sm.Barn.GetMonitorByIndex(deviceId)
	if err != nil {
		c.String(400, "Device not found")
		return
	}

	resp := stringResponse{
		Value: device.GetDescription(),
	}
	sm.prepareAlpacaResponse(c, &resp.alpacaResponse)
	c.IndentedJSON(http.StatusOK, resp)
}

// handleDriverInfo handles GET requests for safety monitor driverinfo property
func (sm *SafetyMonitorAPI) handleDriverInfo(c *gin.Context) {
	resp := stringResponse{
		Value: "Alpaca Barn safety monitor",
	}
	sm.prepareAlpacaResponse(c, &resp.alpacaResponse)
	c.IndentedJSON(http.StatusOK, resp)
}

// handleDriverVersion handles GET requests for safety monitor driverversion property
func (sm *SafetyMonitorAPI) handleDriverVersion(c *gin.Context) {
	bi, _ := debug.ReadBuildInfo()
	resp := stringResponse{
		Value: bi.Main.Version,
	}
	sm.prepareAlpacaResponse(c, &resp.alpacaResponse)
	c.IndentedJSON(http.StatusOK, resp)
}

// handleSupportedActions handles GET requests for safety monitor supportedactions property
func (sm *SafetyMonitorAPI) handleSupportedActions(c *gin.Context) {
	resp := stringlistResponse{
		Value: []string{"RawValue"},
	}
	sm.prepareAlpacaResponse(c, &resp.alpacaResponse)
	c.IndentedJSON(http.StatusOK, resp)
}

// handleInterfaceVersion handles GET requests for safety monitor interfaceversion property
func (sm *SafetyMonitorAPI) handleInterfaceVersion(c *gin.Context) {
	resp := int32Response{
		Value: 2,
	}
	sm.prepareAlpacaResponse(c, &resp.alpacaResponse)
	c.IndentedJSON(http.StatusOK, resp)
}

// handleDeviceState handles GET requests for safety monitor devicestate property
func (sm *SafetyMonitorAPI) handleDeviceState(c *gin.Context) {
	deviceId, _ := strconv.Atoi(c.Param("device_id"))
	device, err := sm.Barn.GetMonitorByIndex(deviceId)
	if err != nil {
		c.String(400, "Device not found")
		return
	}

	if !sm.isRequestConnected(c) {
		c.String(400, "Not connected")
		return
	}

	// Create device state array with IsSafe and TimeStamp properties
	deviceStates := []DeviceState{
		{
			Name:  "IsSafe",
			Value: device.IsSafe(),
		},
		{
			Name:  "TimeStamp",
			Value: device.GetTimeStamp(),
		},
	}

	resp := deviceStateResponse{
		Value: deviceStates,
	}
	sm.prepareAlpacaResponse(c, &resp.alpacaResponse)
	c.IndentedJSON(http.StatusOK, resp)
}

// handleSafetyMonitorConnect handles PUT requests to connect/disconnect safety monitors
func (sm *SafetyMonitorAPI) handleSafetyMonitorConnect(c *gin.Context) {
	deviceId, _ := strconv.Atoi(c.Param("device_id"))
	connected := c.PostForm("Connected")
	_, err := sm.Barn.GetMonitorByIndex(deviceId)
	if err != nil {
		c.String(400, "Device not found")
		return
	}
	dev := sm.Devices["safetymonitor"][deviceId]
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
	sm.prepareAlpacaResponse(c, &resp)
	c.IndentedJSON(http.StatusOK, resp)
}

// handleSafetyMonitorAction handles PUT requests for safety monitor actions
func (sm *SafetyMonitorAPI) handleSafetyMonitorAction(c *gin.Context) {
	deviceId, _ := strconv.Atoi(c.Param("device_id"))
	device, err := sm.Barn.GetMonitorByIndex(deviceId)
	if err != nil {
		c.String(400, "Device not found")
		return
	}

	action := c.PostForm("Action")
	_, err = sm.Barn.GetMonitorByIndex(deviceId)
	if err != nil {
		c.String(400, "Device not found")
		return
	}
	if !sm.isRequestConnected(c) {
		c.String(400, "Not connected")
		return
	}
	if action == "RawValue" {
		resp := stringResponse{
			Value: device.GetRawValue(),
		}
		sm.prepareAlpacaResponse(c, &resp.alpacaResponse)
		c.IndentedJSON(http.StatusOK, resp)
	} else {
		c.String(400, "The device did not understand which operation was being requested or insufficient information was given to complete the operation.")
		return
	}
}

// handleConnect handles PUT requests to start an asynchronous connection to the device
func (sm *SafetyMonitorAPI) handleConnect(c *gin.Context) {
	deviceId, _ := strconv.Atoi(c.Param("device_id"))
	_, err := sm.Barn.GetMonitorByIndex(deviceId)
	if err != nil {
		c.String(400, "Device not found")
		return
	}
	dev := sm.Devices["safetymonitor"][deviceId]
	id := getFullClientId(c)

	// Start asynchronous connection process
	// In a real implementation, this would typically start a goroutine
	// For now, we'll connect immediately but log it as an async operation
	log.WithFields(log.Fields{
		"deviceid": deviceId,
		"clientid": id,
	}).Info("[BARN] Starting asynchronous connection to safety monitor")

	// Connect the client
	dev.ConnectClient(id)

	resp := alpacaResponse{
		ClientTransactionID: 0,
		ServerTransactionID: 0,
		ErrorNumber:         0,
		ErrorMessage:        "",
	}
	sm.prepareAlpacaResponse(c, &resp)
	c.IndentedJSON(http.StatusOK, resp)
}

// handleDisconnect handles PUT requests to start an asynchronous disconnection from the device
func (sm *SafetyMonitorAPI) handleDisconnect(c *gin.Context) {
	deviceId, _ := strconv.Atoi(c.Param("device_id"))
	_, err := sm.Barn.GetMonitorByIndex(deviceId)
	if err != nil {
		c.String(400, "Device not found")
		return
	}
	dev := sm.Devices["safetymonitor"][deviceId]
	id := getFullClientId(c)

	// Start asynchronous disconnection process
	// In a real implementation, this would typically start a goroutine
	// For now, we'll disconnect immediately but log it as an async operation
	log.WithFields(log.Fields{
		"deviceid": deviceId,
		"clientid": id,
	}).Info("[BARN] Starting asynchronous disconnection from safety monitor")

	// Disconnect the client
	dev.DisconnectClient(id)

	resp := alpacaResponse{
		ClientTransactionID: 0,
		ServerTransactionID: 0,
		ErrorNumber:         0,
		ErrorMessage:        "",
	}
	sm.prepareAlpacaResponse(c, &resp)
	c.IndentedJSON(http.StatusOK, resp)
}
