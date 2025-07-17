package api

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/thebuh/barn/internal/app"
)

type ClientId string

type ConnectedClient struct {
	ClientId            ClientId
	ClientTransactionID uint32
	Connected           bool
	// Weather-specific state
	WeatherState *WeatherClientState
}

// WeatherClientState holds weather-specific state for each connected client
type WeatherClientState struct {
	AveragePeriod float64
}

type ApiServer struct {
	ApiPort             uint32
	Barn                app.Server
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
		WeatherState:        &WeatherClientState{AveragePeriod: 0.0}, // Default average period
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

// GetWeatherClientState returns the weather state for a specific client
func (d *Device) GetWeatherClientState(id ClientId) *WeatherClientState {
	if client, exists := d.ConnectedClients[id]; exists && client.Connected {
		return client.WeatherState
	}
	return nil
}

// SetWeatherAveragePeriod sets the average period for a specific client
func (d *Device) SetWeatherAveragePeriod(id ClientId, averagePeriod float64) bool {
	if client, exists := d.ConnectedClients[id]; exists && client.Connected {
		if client.WeatherState == nil {
			client.WeatherState = &WeatherClientState{}
		}
		client.WeatherState.AveragePeriod = averagePeriod
		return true
	}
	return false
}

// GetWeatherAveragePeriod gets the average period for a specific client
func (d *Device) GetWeatherAveragePeriod(id ClientId) (float64, bool) {
	if client, exists := d.ConnectedClients[id]; exists && client.Connected {
		if client.WeatherState != nil {
			return client.WeatherState.AveragePeriod, true
		}
	}
	return 0.0, false
}

func NewApiServer(barn app.Server, apiPort uint32) *ApiServer {
	return &ApiServer{
		ApiPort: apiPort,
		Barn:    barn,
		Devices: make(map[string]map[int]*Device),
	}
}

func (srv *ApiServer) Start() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	// Add global middleware to inject API server into context
	router.Use(func(c *gin.Context) {
		c.Set("apiServer", srv)
		c.Next()
	})
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Alpaca Barn server")
	})

	// Initialize safety monitor devices
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

	// Initialize weather devices
	srv.Devices["observingconditions"] = make(map[int]*Device)
	weatherIds := srv.Barn.GetWeatherIds()
	for i, id := range weatherIds {
		srv.Devices["observingconditions"][i] = &Device{
			Id:               id,
			Type:             "observingconditions",
			Index:            i,
			ConnectedClients: make(map[ClientId]*ConnectedClient),
		}
	}

	srv.configureManagementAPI(router)

	// Configure separate API handlers
	safetyMonitorAPI := NewSafetyMonitorAPI(srv)
	safetyMonitorAPI.ConfigureRoutes(router)

	weatherAPI := NewWeatherAPI(srv)
	weatherAPI.ConfigureRoutes(router)

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
				ManufacturerVersion: "Version:" + bi.Main.Version,
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

		// Add safety monitor devices
		ids := srv.Barn.GetMonitorIds()
		for i, id := range ids {
			val = append(val, DeviceConfiguration{
				DeviceName:   srv.Barn.GetMonitor(id).GetName(),
				DeviceType:   "safetymonitor",
				DeviceNumber: i,
				UniqueID:     id,
			})
		}

		// Add weather devices
		weatherIds := srv.Barn.GetWeatherIds()
		for i, id := range weatherIds {
			val = append(val, DeviceConfiguration{
				DeviceName:   srv.Barn.GetWeather(id).GetName(),
				DeviceType:   "observingconditions",
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
		cidv = getQuery(c, "ClientID")
		if cidv == "" {
			return -1
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

func getQuery(c *gin.Context, targetKey string) string {
	queryParams := c.Request.URL.Query()
	for key, values := range queryParams {
		if strings.EqualFold(targetKey, key) {
			return strings.Join(values, ",") // if there are multiple values keyed to an input parameter, they'll be returned comma-separated
		}
	}
	return ""
}

func getClientTransactionId(c *gin.Context) int {
	ctidv := ""
	if c.Request.Method == "GET" {
		ctidv = getQuery(c, "ClientTransactionID")
		if ctidv == "" {
			return -1
		}
	} else {
		ctidv = c.PostForm("ClientTransactionID")
		if ctidv == "" {
			ctidv = c.PostForm("clienttransactionid")
			if ctidv == "" {
				ctidv = "0"
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
