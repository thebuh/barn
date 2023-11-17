package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
)

type ClientId string

type ConnectedClient struct {
	ClientId            ClientId
	ClientTransactionID uint32
	Connected           bool
}

type ApiServer struct {
	ApiPort             uint32
	Barn                Server
	ServerTransactionID uint32
	Devices             map[string]map[int]*Device
}

type Device struct {
	Id               string
	Type             string
	Index            int
	ConnectedClients map[ClientId]*ConnectedClient
}

func (d *Device) IsConnected(id ClientId) bool {
	for clientId, client := range d.ConnectedClients {
		if clientId == id {
			return client.Connected
		}
	}
	return false
}
func (d *Device) ConnectClient(id ClientId) {
	for clientId, client := range d.ConnectedClients {
		if clientId != id {
			continue
		}
		if client.Connected == false {
			client.Connected = true
			return
		}
		return
	}
	cc := &ConnectedClient{
		ClientId:            id,
		ClientTransactionID: 0,
		Connected:           true,
	}
	d.ConnectedClients[id] = cc
}

func (d *Device) DisconnectClient(id ClientId) {
	for clientId := range d.ConnectedClients {
		if clientId != id {
			continue
		}
		delete(d.ConnectedClients, id)
	}
}

func NewApiServer(barn Server, apiPort uint32) *ApiServer {
	return &ApiServer{
		ApiPort: apiPort,
		Barn:    barn,
		Devices: make(map[string]map[int]*Device),
	}
}

func (srv *ApiServer) Start() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Alpaca Barn server")
	})
	srv.Devices["safetymonitor"] = make(map[int]*Device)
	ids := srv.Barn.GetMonitorIds()
	for i, id := range ids {
		srv.Devices["safetymonitor"][i] = &Device{
			Id:               id,
			Type:             "safetymonitor",
			Index:            i,
			ConnectedClients: make(map[ClientId]*ConnectedClient),
		}
	}
	srv.configureManagementAPI(router)
	srv.configureSafetyMonitorAPI(router)
	err := router.Run(fmt.Sprintf("0.0.0.0:%d", srv.ApiPort))
	if err != nil {
		return
	}
}

func (srv *ApiServer) configureManagementAPI(router *gin.Engine) {
	router.GET("/management/apiversions", func(c *gin.Context) {
		resp := uint32listResponse{
			Value: []uint32{1},
			alpacaResponse: alpacaResponse{
				ClientTransactionID: 0,
				ServerTransactionID: 0,
				ErrorNumber:         0,
				ErrorMessage:        "",
			},
		}
		srv.prepareAlpacaResponse(c, &resp.alpacaResponse)
		c.IndentedJSON(http.StatusOK, resp)
	})
	router.GET("/management/v1/description", func(c *gin.Context) {
		bi, _ := debug.ReadBuildInfo()
		resp := managementDescriptionResponse{
			Value: ServerDescription{
				ServerName:          "Alpaca Barn",
				Manufacturer:        "https://github.com/thebuh/barn",
				ManufacturerVersion: bi.Main.Version,
				Location:            "Location string",
			},
			alpacaResponse: alpacaResponse{
				ClientTransactionID: 0,
				ServerTransactionID: 0,
				ErrorNumber:         0,
				ErrorMessage:        "",
			},
		}
		srv.prepareAlpacaResponse(c, &resp.alpacaResponse)
		c.IndentedJSON(http.StatusOK, resp)
	})
	router.GET("/management/v1/configureddevices", func(c *gin.Context) {
		var val []DeviceConfiguration

		ids := srv.Barn.GetMonitorIds()
		for i, id := range ids {
			val = append(val, DeviceConfiguration{
				DeviceName:   srv.Barn.GetMonitor(id).GetName(),
				DeviceType:   "safetymonitor",
				DeviceNumber: i,
				UniqueID:     id,
			})
		}
		resp := managementDevicesListResponse{
			Value: val,
			alpacaResponse: alpacaResponse{
				ClientTransactionID: 0,
				ServerTransactionID: 0,
				ErrorNumber:         0,
				ErrorMessage:        "",
			},
		}
		srv.prepareAlpacaResponse(c, &resp.alpacaResponse)
		c.IndentedJSON(http.StatusOK, resp)
	})

}

func (srv *ApiServer) configureSafetyMonitorAPI(router *gin.Engine) {
	router.PUT("/api/v1/safetymonitor/:device_id/connected", srv.handleSafetyMonitorConnect)
	router.PUT("/api/v1/safetymonitor/:device_id/action", srv.handleSafetyMonitorAction)
	router.GET("/api/v1/safetymonitor/:device_id/:action", srv.handleSafetyMonitorRequest)

}

func (srv *ApiServer) prepareAlpacaResponse(c *gin.Context, resp *alpacaResponse) {
	ctid := getClientTransactionId(c)
	if ctid < 0 {
		ctid = 0
	}
	srv.ServerTransactionID += 1
	resp.ClientTransactionID = uint32(ctid)
	resp.ServerTransactionID = srv.ServerTransactionID

}

func getClientId(c *gin.Context) int {
	cidv := ""
	if c.Request.Method == "GET" {
		cidv = c.DefaultQuery("ClientID", "")
		if cidv == "" {
			cidv = c.DefaultQuery("clientid", "")
			if cidv == "" {
				return -1
			}
		}
	} else {
		cidv = c.PostForm("ClientID")
		if cidv == "" {
			cidv = c.PostForm("clientid")
			if cidv == "" {
				return -1
			}
		}
	}
	cid, err := strconv.Atoi(cidv)
	if err != nil {
		return -1
	}
	if cid < 0 {
		return -1
	}
	return cid
}

func getFullClientId(c *gin.Context) ClientId {
	return ClientId(c.RemoteIP() + "-" + strconv.Itoa(getClientId(c)))
}

func getClientTransactionId(c *gin.Context) int {
	ctidv := ""
	if c.Request.Method == "GET" {
		ctidv = c.DefaultQuery("ClientTransactionID", "")
		if ctidv == "" {
			ctidv = c.DefaultQuery("clienttransactionid", "")
			if ctidv == "" {
				return -1
			}
		}
	} else {
		ctidv = c.PostForm("ClientTransactionID")
		if ctidv == "" {
			ctidv = c.PostForm("ClientTransactionID")
			if ctidv == "" {
				return -1
			}
		}
	}
	ctid, err := strconv.Atoi(ctidv)
	if err != nil {
		return -1
	}
	if ctid < 0 {
		return -1
	}
	return ctid
}

func (srv *ApiServer) validAlpacaRequest(c *gin.Context) bool {
	if _, err := strconv.Atoi(c.Param("device_id")); err != nil {
		return false
	}
	cid := getClientId(c)
	if cid < 0 {
		return false
	}

	ctidv := getClientTransactionId(c)
	if ctidv < 0 {
		return false
	}
	return true
}

func (srv *ApiServer) isRequestConnected(c *gin.Context) bool {
	deviceId, _ := strconv.Atoi(c.Param("device_id"))
	device := srv.Devices["safetymonitor"][deviceId]
	return device.IsConnected(getFullClientId(c))
}

func (srv *ApiServer) handleSafetyMonitorRequest(c *gin.Context) {
	if !srv.validAlpacaRequest(c) {
		c.String(400, "Invalid request")
		return
	}
	deviceId, _ := strconv.Atoi(c.Param("device_id"))
	device, err := srv.Barn.GetMonitorByIndex(deviceId)
	if err != nil {
		c.String(400, "Device not found")
		return
	}
	action := strings.ToLower(c.Param("action"))
	if action == "issafe" {
		val := false
		if srv.isRequestConnected(c) {
			val = device.IsSafe()
		}
		resp := booleanResponse{
			Value: val,
		}
		srv.prepareAlpacaResponse(c, &resp.alpacaResponse)
		c.IndentedJSON(http.StatusOK, resp)
	} else if action == "connected" {
		result := srv.isRequestConnected(c)
		resp := booleanResponse{
			Value: result,
		}
		srv.prepareAlpacaResponse(c, &resp.alpacaResponse)
		c.IndentedJSON(http.StatusOK, resp)
	} else if action == "name" {
		resp := stringResponse{
			Value: device.GetName(),
		}
		srv.prepareAlpacaResponse(c, &resp.alpacaResponse)
		c.IndentedJSON(http.StatusOK, resp)
	} else if action == "description" {
		resp := stringResponse{
			Value: device.GetDescription(),
		}
		srv.prepareAlpacaResponse(c, &resp.alpacaResponse)
		c.IndentedJSON(http.StatusOK, resp)
	} else if action == "driverinfo" {
		resp := stringResponse{
			Value: "Alpaca Barn safety monitor",
		}
		srv.prepareAlpacaResponse(c, &resp.alpacaResponse)
		c.IndentedJSON(http.StatusOK, resp)
	} else if action == "driverversion" {
		bi, _ := debug.ReadBuildInfo()
		resp := stringResponse{
			Value: bi.Main.Version,
		}
		srv.prepareAlpacaResponse(c, &resp.alpacaResponse)
		c.IndentedJSON(http.StatusOK, resp)
	} else if action == "supportedactions" {
		resp := stringlistResponse{
			Value: []string{"RawValue"},
		}
		srv.prepareAlpacaResponse(c, &resp.alpacaResponse)
		c.IndentedJSON(http.StatusOK, resp)
	} else if action == "interfaceversion" {
		resp := int32Response{
			Value: 2,
		}
		srv.prepareAlpacaResponse(c, &resp.alpacaResponse)
		c.IndentedJSON(http.StatusOK, resp)
	} else {
		c.String(400, "The device did not understand which operation was being requested or insufficient information was given to complete the operation.")
		return
	}
}

func (srv *ApiServer) handleSafetyMonitorConnect(c *gin.Context) {
	if !srv.validAlpacaRequest(c) {
		c.String(400, "Invalid request")
		return
	}
	deviceId, _ := strconv.Atoi(c.Param("device_id"))
	connected := c.PostForm("Connected")
	_, err := srv.Barn.GetMonitorByIndex(deviceId)
	if err != nil {
		c.String(400, "Device not found")
		return
	}
	dev := srv.Devices["safetymonitor"][deviceId]
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
	srv.prepareAlpacaResponse(c, &resp)
	c.IndentedJSON(http.StatusOK, resp)
}

func (srv *ApiServer) handleSafetyMonitorAction(c *gin.Context) {
	if !srv.validAlpacaRequest(c) {
		c.String(400, "Invalid request")
		return
	}
	deviceId, _ := strconv.Atoi(c.Param("device_id"))
	device, err := srv.Barn.GetMonitorByIndex(deviceId)
	if err != nil {
		c.String(400, "Device not found")
		return
	}

	action := c.PostForm("Action")
	_, err = srv.Barn.GetMonitorByIndex(deviceId)
	if err != nil {
		c.String(400, "Device not found")
		return
	}
	if !srv.isRequestConnected(c) {
		c.String(400, "Not connected")
		return
	}
	if action == "RawValue" {
		resp := stringResponse{
			Value: device.GetRawValue(),
		}
		srv.prepareAlpacaResponse(c, &resp.alpacaResponse)
		c.IndentedJSON(http.StatusOK, resp)
	} else {
		c.String(400, "The device did not understand which operation was being requested or insufficient information was given to complete the operation.")
		return
	}
}
