package main

type alpacaResponse struct {
	ClientTransactionID uint32 `json:"ClientTransactionID"`
	ServerTransactionID uint32 `json:"ServerTransactionID"`
	ErrorNumber         int32  `json:"ErrorNumber"`
	ErrorMessage        string `json:"ErrorMessage"`
}

type stringResponse struct {
	Value string `json:"Value"`
	alpacaResponse
}

type stringlistResponse struct {
	Value []string `json:"Value"`
	alpacaResponse
}

type booleanResponse struct {
	Value bool `json:"Value"`
	alpacaResponse
}

type float64Response struct {
	Value float64 `json:"Value"`
	alpacaResponse
}

type int32Response struct {
	Value int32 `json:"Value"`
	alpacaResponse
}

type uint32listResponse struct {
	Value []uint32 `json:"Value"`
	alpacaResponse
}

type uint32Rank2ArrayResponse struct {
	Value [][]uint32 `json:"Value"`
	Rank  uint32     `json:"Rank"`
	alpacaResponse
}

type putResponse struct {
	alpacaResponse
}

type managementDevicesListResponse struct {
	Value []DeviceConfiguration `json:"Value"`
	alpacaResponse
}

type DeviceConfiguration struct {
	DeviceName   string `json:"DeviceName"`
	DeviceType   string `json:"DeviceType"`
	DeviceNumber int    `json:"DeviceNumber"`
	UniqueID     string `json:"DeviceId"`
}

type managementDescriptionResponse struct {
	Value ServerDescription `json:"Value"`
	alpacaResponse
}

type ServerDescription struct {
	ServerName          string `json:"ServerName"`
	Manufacturer        string `json:"Manufacturer"`
	ManufacturerVersion string `json:"ManufacturerVersion"`
	Location            string `json:"Location"`
}
