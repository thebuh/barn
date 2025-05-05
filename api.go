package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type ClientIp string

type ConnectedClient struct {
	Address             ClientIp
	ClientTransactionID uint32
	Connected           bool
}

type ApiServer struct {
	ApiPort             int
	Barn                Server
	ServerTransactionID uint32
	Devices             map[string]map[int]Device
}

type Device struct {
	Id               string
	Type             string
	Index            int
	ConnectedClients map[ClientIp]ConnectedClient
}

func (d *Device) IsConnected(ip ClientIp) bool {
	for clientIp, client := range d.ConnectedClients {
		if clientIp == ip {
			return client.Connected
		}
	}
	return false
}
func (d *Device) ConnectClient(ip ClientIp) {
	if d.ConnectedClients == nil {
		d.ConnectedClients = make(map[ClientIp]ConnectedClient)
	}
	for clientIp, client := range d.ConnectedClients {
		if clientIp == ip {
			if client.Connected == false {
				client.Connected = true
				return
			}
			return
		}
	}
	cc := ConnectedClient{
		Address:             ip,
		ClientTransactionID: 0,
		Connected:           true,
	}
	d.ConnectedClients[ip] = cc
}

func (d *Device) DisconnectClient(ip ClientIp) {
	for clientIp, _ := range d.ConnectedClients {
		if clientIp == ip {
			delete(d.ConnectedClients, ip)
		}
	}
}

func NewApiServer(barn Server, apiPort int) *ApiServer {
	return &ApiServer{
		ApiPort: apiPort,
		Barn:    barn,
		Devices: make(map[string]map[int]Device),
	}
}

func (srv *ApiServer) Start() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Alpaca Barn server")
	})
	srv.configureManagementAPI(router)
	srv.configureSafetyMonitorAPI(router)
	router.Run(fmt.Sprintf("0.0.0.0:%d", srv.ApiPort))
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
		resp := managementDescriptionResponse{
			Value: ServerDescription{
				ServerName:          "Alpaca Barn",
				Manufacturer:        "https://github.com/thebuh/barn",
				ManufacturerVersion: Version,
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
		srv.Devices["safetymonitor"] = make(map[int]Device)
		ids := srv.Barn.GetMonitorIds()
		for i, id := range ids {
			val = append(val, DeviceConfiguration{
				DeviceName:   srv.Barn.GetMonitor(id).GetName(),
				DeviceType:   "safetymonitor",
				DeviceNumber: i,
				UniqueID:     id,
			})
			srv.Devices["safetymonitor"][i] = Device{
				Id:               id,
				Type:             "safetymonitor",
				Index:            i,
				ConnectedClients: make(map[ClientIp]ConnectedClient),
			}
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
	router.GET("/api/v1/safetymonitor/:device_id/:action", srv.handleSafetyMonitorRequest)
}

func (srv *ApiServer) prepareAlpacaResponse(c *gin.Context, resp *alpacaResponse) {
	cid, _ := strconv.Atoi(c.DefaultQuery("ClientTransactionID", "0"))
	resp.ClientTransactionID = uint32(cid)
	resp.ServerTransactionID = srv.ServerTransactionID
	srv.ServerTransactionID += 1
}

func (srv *ApiServer) handleSafetyMonitorRequest(c *gin.Context) {
	deviceId, _ := strconv.Atoi(c.Param("device_id"))
	device, err := srv.Barn.GetMonitorByIndex(deviceId)
	if err != nil {
		c.String(400, "Device not found")
		return
	}
	action := strings.ToLower(c.Param("action"))
	ip := ClientIp(c.RemoteIP())
	d := srv.Devices["safetymonitor"][deviceId]
	if action == "is_safe" || action == "issafe" {
		val := false
		if d.IsConnected(ip) {
			val = device.IsSafe()
		}
		resp := booleanResponse{
			Value: val,
		}
		srv.prepareAlpacaResponse(c, &resp.alpacaResponse)
		c.IndentedJSON(http.StatusOK, resp)
	} else if action == "connected" {
		ip := ClientIp(c.RemoteIP())
		result := d.IsConnected(ip)
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
		resp := stringResponse{
			Value: Version,
		}
		srv.prepareAlpacaResponse(c, &resp.alpacaResponse)
		c.IndentedJSON(http.StatusOK, resp)
	} else if action == "supportedactions" {
		resp := stringlistResponse{
			Value: []string{"issafe"},
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
	deviceId, _ := strconv.Atoi(c.Param("device_id"))
	connected := c.PostForm("Connected")
	_, err := srv.Barn.GetMonitorByIndex(deviceId)
	if err != nil {
		c.String(400, "Device not found")
		return
	}
	dev := srv.Devices["safetymonitor"][deviceId]
	ip := ClientIp(c.RemoteIP())

	if strings.ToLower(connected) == "true" { // Connect
		dev.ConnectClient(ip)
	} else { //Disconnect
		dev.DisconnectClient(ip)
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

func (srv *ApiServer) isConnected() {

}
