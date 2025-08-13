package weather

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Common errors
var (
	ErrInvalidPeriod = errors.New("average period must be greater than 0")
	ErrInvalidURL    = errors.New("invalid URL provided")
)

// Sensor definitions
const (
	SensorAveragePeriod  = "AveragePeriod"
	SensorCloudCover     = "CloudCover"
	SensorDewPoint       = "DewPoint"
	SensorHumidity       = "Humidity"
	SensorPressure       = "Pressure"
	SensorRainRate       = "RainRate"
	SensorSkyBrightness  = "SkyBrightness"
	SensorSkyQuality     = "SkyQuality"
	SensorSkyTemperature = "SkyTemperature"
	SensorStarFWHM       = "StarFWHM"
	SensorTemperature    = "Temperature"
	SensorWindDirection  = "WindDirection"
	SensorWindGust       = "WindGust"
	SensorWindSpeed      = "WindSpeed"
	SensorTimeStamp      = "TimeStamp"
)

// AvailableSensors contains all supported sensor names
var AvailableSensors = map[string]bool{
	SensorAveragePeriod:  true,
	SensorCloudCover:     false,
	SensorDewPoint:       true,
	SensorHumidity:       true,
	SensorPressure:       true,
	SensorRainRate:       true,
	SensorSkyBrightness:  false,
	SensorSkyQuality:     false,
	SensorSkyTemperature: false,
	SensorStarFWHM:       false,
	SensorTemperature:    true,
	SensorWindDirection:  true,
	SensorWindGust:       true,
	SensorWindSpeed:      true,
}

// SensorDescriptions maps sensor names to their descriptions
var SensorDescriptions = map[string]string{
	SensorAveragePeriod:  "Average period for weather measurements",
	SensorCloudCover:     "Cloud cover percentage",
	SensorDewPoint:       "Dew point temperature",
	SensorHumidity:       "Relative humidity percentage",
	SensorPressure:       "Atmospheric pressure",
	SensorRainRate:       "Rain rate measurement",
	SensorSkyBrightness:  "Sky brightness measurement",
	SensorSkyQuality:     "Sky quality measurement",
	SensorSkyTemperature: "Sky temperature measurement",
	SensorStarFWHM:       "Star full width at half maximum",
	SensorTemperature:    "Ambient temperature",
	SensorWindDirection:  "Wind direction in degrees",
	SensorWindGust:       "Wind gust speed",
	SensorWindSpeed:      "Wind speed measurement",
}

// IsValidSensor checks if a sensor name is valid (case-insensitive)
func IsValidSensor(sensorName string) bool {
	for key := range AvailableSensors {
		if strings.EqualFold(key, sensorName) {
			return true
		}
	}
	return false
}

// IsSensorAvailable checks if a sensor is both valid and available (enabled)
func IsSensorAvailable(sensorName string) bool {
	if !IsValidSensor(sensorName) {
		return false
	}
	return AvailableSensors[sensorName]
}

// GetSensorDescription returns the description for a given sensor name
func GetSensorDescription(sensorName string) (string, bool) {
	description, exists := SensorDescriptions[sensorName]
	return description, exists
}

// GetAvailableSensors returns a list of all available sensor names
func GetAvailableSensors() []string {
	sensors := make([]string, 0, len(AvailableSensors))
	for sensor, enabled := range AvailableSensors {
		if enabled {
			sensors = append(sensors, sensor)
		}
	}
	return sensors
}

// ObservingConditions defines the interface for weather observing conditions
// following the ASCOM Alpaca standard.
type ObservingConditions interface {
	// Refresh updates the current weather conditions
	Refresh() error

	// GetId returns the unique identifier of the weather station
	GetId() string

	// GetName returns the name of the weather station
	GetName() string

	// GetDescription returns a description of the weather station
	GetDescription() string

	GetState() string

	// ASCOM Alpaca observing conditions methods
	GetAveragePeriod() float64
	SetAveragePeriod(period float64) error
	GetCloudCover() float64
	GetDewPoint() float64
	GetHumidity() float64
	GetPressure() float64
	GetRainRate() float64
	GetSkyBrightness() float64
	GetSkyQuality() float64
	GetSkyTemperature() float64
	GetStarFWHM() float64
	GetTemperature() float64
	GetWindDirection() float64
	GetWindGust() float64
	GetWindSpeed() float64
	GetTimeSinceLastUpdate() float64
}

// WeatherCondition represents the current weather conditions
type WeatherCondition struct {
	AveragePeriod  float64 `json:"average_period"`
	CloudCover     float64 `json:"cloud_cover"`
	DewPoint       float64 `json:"dew_point"`
	Humidity       float64 `json:"humidity"`
	Pressure       float64 `json:"pressure"`
	RainRate       float64 `json:"rain_rate"`
	SkyBrightness  float64 `json:"sky_brightness"`
	SkyQuality     float64 `json:"sky_quality"`
	SkyTemperature float64 `json:"sky_temperature"`
	StarFWHM       float64 `json:"star_fwhm"`
	Temperature    float64 `json:"temperature"`
	WindDirection  float64 `json:"wind_direction"`
	WindGust       float64 `json:"wind_gust"`
	WindSpeed      float64 `json:"wind_speed"`
}

// BaseObservingConditions contains common fields for all observing conditions implementations
type BaseObservingConditions struct {
	id              string
	name            string
	description     string
	lastRefreshTime time.Time
	condition       WeatherCondition
}

// ObservingConditionsDummy implements ObservingConditions with static values
type ObservingConditionsDummy struct {
	BaseObservingConditions
}

// NewObservingConditionsDummy creates a new dummy weather station
func NewObservingConditionsDummy(id string, name string, description string) *ObservingConditionsDummy {
	return &ObservingConditionsDummy{
		BaseObservingConditions: BaseObservingConditions{
			id:          id,
			name:        name,
			description: description,
		},
	}
}

func (o *ObservingConditionsDummy) Refresh() error {
	// Dummy implementation doesn't need to do anything
	return nil
}

func (o *ObservingConditionsDummy) GetId() string {
	return o.id
}

func (o *ObservingConditionsDummy) GetName() string {
	return o.name
}

func (o *ObservingConditionsDummy) GetDescription() string {
	return o.description
}

func (o *ObservingConditionsDummy) GetAveragePeriod() float64 {
	return o.condition.AveragePeriod
}

func (o *ObservingConditionsDummy) SetAveragePeriod(period float64) error {
	if period < 0 {
		return ErrInvalidPeriod
	}
	o.condition.AveragePeriod = period
	return nil
}

func (o *ObservingConditionsDummy) GetCloudCover() float64 {
	return o.condition.CloudCover
}

func (o *ObservingConditionsDummy) GetDewPoint() float64 {
	return o.condition.DewPoint
}

func (o *ObservingConditionsDummy) GetHumidity() float64 {
	return o.condition.Humidity
}

func (o *ObservingConditionsDummy) GetPressure() float64 {
	return o.condition.Pressure
}

func (o *ObservingConditionsDummy) GetRainRate() float64 {
	return o.condition.RainRate
}

func (o *ObservingConditionsDummy) GetSkyBrightness() float64 {
	return o.condition.SkyBrightness
}

func (o *ObservingConditionsDummy) GetSkyQuality() float64 {
	return o.condition.SkyQuality
}

func (o *ObservingConditionsDummy) GetSkyTemperature() float64 {
	return o.condition.SkyTemperature
}

func (o *ObservingConditionsDummy) GetStarFWHM() float64 {
	return o.condition.StarFWHM
}

func (o *ObservingConditionsDummy) GetTemperature() float64 {
	return o.condition.Temperature
}

func (o *ObservingConditionsDummy) GetWindDirection() float64 {
	return o.condition.WindDirection
}

func (o *ObservingConditionsDummy) GetWindGust() float64 {
	return o.condition.WindGust
}

func (o *ObservingConditionsDummy) GetWindSpeed() float64 {
	return o.condition.WindSpeed
}

func (o *ObservingConditionsDummy) GetTimeSinceLastUpdate() float64 {
	return time.Since(o.lastRefreshTime).Seconds()
}

func (o *ObservingConditionsDummy) GetState() string {
	json, _ := json.Marshal(o.condition)
	return string(json)
}

// ObservingConditionsHttp implements ObservingConditions by fetching data from an HTTP endpoint
type ObservingConditionsHttp struct {
	BaseObservingConditions
	url    string
	client *http.Client
}

// NewObservingConditionsHttp creates a new HTTP-based weather station
func NewObservingConditionsHttp(id string, name string, description string, url string) (*ObservingConditionsHttp, error) {
	if url == "" {
		return nil, ErrInvalidURL
	}

	cond := &ObservingConditionsHttp{
		BaseObservingConditions: BaseObservingConditions{
			id:          id,
			name:        name,
			description: description,
		},
		url: url,
	}
	cond.client = &http.Client{
		Timeout: 5 * time.Second,
	}
	cond.SetAveragePeriod(0)
	cond.Refresh()
	return cond, nil
}

func (o *ObservingConditionsHttp) Refresh() error {
	resp, err := o.client.Get(o.url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP request failed with status code %d", resp.StatusCode)
	}

	buf := make([]byte, 4096)
	n, err := resp.Body.Read(buf)
	if err != nil && err != io.EOF {
		return err
	}
	buf = buf[:n]
	_ = resp.Body.Close()
	content := string(buf)
	// Parse the JSON content
	var weatherData struct {
		ID             int     `json:"id"`
		IndoorTemp     float64 `json:"indoortemp"`
		Temp           float64 `json:"temp"`
		DewPt          float64 `json:"dewpt"`
		WindChill      float64 `json:"windchill"`
		IndoorHumidity int     `json:"indoorhumidity"`
		Humidity       int     `json:"humidity"`
		WindSpeedMS    float64 `json:"windspeedms"`
		WindGustMS     float64 `json:"windgustms"`
		WindDir        int     `json:"winddir"`
		AbsBaroMin     float64 `json:"absbaromin"`
		BaroMin        float64 `json:"baromin"`
		RainIn         float64 `json:"rainin"`
		DailyRainIn    float64 `json:"dailyrainin"`
		WeeklyRainIn   float64 `json:"weeklyrainin"`
		MonthlyRainIn  float64 `json:"monthlyrainin"`
		SolarRadiation float64 `json:"solarradiation"`
		UV             int     `json:"UV"`
		DateUTC        string  `json:"dateutc"`
		SoftwareType   string  `json:"softwaretype"`
	}

	if err := json.Unmarshal([]byte(content), &weatherData); err != nil {
		return fmt.Errorf("failed to parse weather data: %w", err)
	}

	// Update the condition values
	o.condition.Temperature = weatherData.Temp
	o.condition.DewPoint = weatherData.DewPt
	o.condition.Humidity = float64(weatherData.Humidity)
	o.condition.Pressure = weatherData.BaroMin
	o.condition.WindSpeed = weatherData.WindSpeedMS
	o.condition.WindGust = weatherData.WindGustMS
	o.condition.WindDirection = float64(weatherData.WindDir)
	o.condition.RainRate = weatherData.RainIn
	o.lastRefreshTime = time.Now()

	fmt.Println("Refreshed weather conditions from", o.url)
	return nil
}

func (o *ObservingConditionsHttp) GetId() string {
	return o.id
}

func (o *ObservingConditionsHttp) GetName() string {
	return o.name
}

func (o *ObservingConditionsHttp) GetDescription() string {
	return o.description
}

func (o *ObservingConditionsHttp) GetAveragePeriod() float64 {
	return o.condition.AveragePeriod
}

func (o *ObservingConditionsHttp) SetAveragePeriod(period float64) error {
	if period < 0 {
		return ErrInvalidPeriod
	}
	o.condition.AveragePeriod = period
	return nil
}

func (o *ObservingConditionsHttp) GetCloudCover() float64 {
	return o.condition.CloudCover
}

func (o *ObservingConditionsHttp) GetDewPoint() float64 {
	return o.condition.DewPoint
}

func (o *ObservingConditionsHttp) GetHumidity() float64 {
	return o.condition.Humidity
}

func (o *ObservingConditionsHttp) GetPressure() float64 {
	return o.condition.Pressure
}

func (o *ObservingConditionsHttp) GetRainRate() float64 {
	return o.condition.RainRate
}

func (o *ObservingConditionsHttp) GetSkyBrightness() float64 {
	return o.condition.SkyBrightness
}

func (o *ObservingConditionsHttp) GetSkyQuality() float64 {
	return o.condition.SkyQuality
}

func (o *ObservingConditionsHttp) GetSkyTemperature() float64 {
	return o.condition.SkyTemperature
}

func (o *ObservingConditionsHttp) GetStarFWHM() float64 {
	return o.condition.StarFWHM
}

func (o *ObservingConditionsHttp) GetTemperature() float64 {
	return o.condition.Temperature
}

func (o *ObservingConditionsHttp) GetWindDirection() float64 {
	return o.condition.WindDirection
}

func (o *ObservingConditionsHttp) GetWindGust() float64 {
	return o.condition.WindGust
}

func (o *ObservingConditionsHttp) GetWindSpeed() float64 {
	return o.condition.WindSpeed
}

func (o *ObservingConditionsHttp) GetTimeSinceLastUpdate() float64 {
	return time.Since(o.lastRefreshTime).Seconds()
}

func (o *ObservingConditionsHttp) GetState() string {
	json, _ := json.Marshal(o.condition)
	return string(json)
}
