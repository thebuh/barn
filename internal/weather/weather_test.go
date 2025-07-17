package weather

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewObservingConditionsDummy(t *testing.T) {
	id := "dummy1"
	name := "Test Dummy"
	desc := "Test Dummy Weather Station"

	dummy := NewObservingConditionsDummy(id, name, desc)

	if dummy.GetId() != id {
		t.Errorf("Expected ID %s, got %s", id, dummy.GetId())
	}

	if dummy.GetName() != name {
		t.Errorf("Expected name %s, got %s", name, dummy.GetName())
	}

	if dummy.GetDescription() != desc {
		t.Errorf("Expected description %s, got %s", desc, dummy.GetDescription())
	}
}

func TestObservingConditionsDummy_Refresh(t *testing.T) {
	dummy := NewObservingConditionsDummy("test", "Test", "Test Station")

	// Refresh should not return an error for dummy implementation
	err := dummy.Refresh()
	if err != nil {
		t.Errorf("Expected no error from Refresh(), got %v", err)
	}
}

func TestObservingConditionsDummy_GetId(t *testing.T) {
	dummy := NewObservingConditionsDummy("test-id", "Test", "Test Station")

	expected := "test-id"
	actual := dummy.GetId()
	if actual != expected {
		t.Errorf("Expected ID %s, got %s", expected, actual)
	}
}

func TestObservingConditionsDummy_GetName(t *testing.T) {
	dummy := NewObservingConditionsDummy("test", "Test Name", "Test Station")

	expected := "Test Name"
	actual := dummy.GetName()
	if actual != expected {
		t.Errorf("Expected name %s, got %s", expected, actual)
	}
}

func TestObservingConditionsDummy_GetDescription(t *testing.T) {
	dummy := NewObservingConditionsDummy("test", "Test", "Test Description")

	expected := "Test Description"
	actual := dummy.GetDescription()
	if actual != expected {
		t.Errorf("Expected description %s, got %s", expected, actual)
	}
}

func TestObservingConditionsDummy_GetState(t *testing.T) {
	dummy := NewObservingConditionsDummy("test", "Test", "Test Station")

	// Set some test values
	dummy.condition = WeatherCondition{
		Temperature: 20.5,
		Humidity:    65.0,
		Pressure:    1013.25,
	}

	state := dummy.GetState()

	// Verify it's valid JSON
	var parsed WeatherCondition
	err := json.Unmarshal([]byte(state), &parsed)
	if err != nil {
		t.Errorf("GetState() should return valid JSON, got error: %v", err)
	}

	if parsed.Temperature != 20.5 {
		t.Errorf("Expected temperature 20.5, got %f", parsed.Temperature)
	}
}

func TestObservingConditionsDummy_GetAveragePeriod(t *testing.T) {
	dummy := NewObservingConditionsDummy("test", "Test", "Test Station")

	// Test default value
	if dummy.GetAveragePeriod() != 0 {
		t.Errorf("Expected default average period 0, got %f", dummy.GetAveragePeriod())
	}

	// Set a value and test
	dummy.condition.AveragePeriod = 15.5
	if dummy.GetAveragePeriod() != 15.5 {
		t.Errorf("Expected average period 15.5, got %f", dummy.GetAveragePeriod())
	}
}

func TestObservingConditionsDummy_SetAveragePeriod(t *testing.T) {
	dummy := NewObservingConditionsDummy("test", "Test", "Test Station")

	// Test valid period
	err := dummy.SetAveragePeriod(10.0)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if dummy.GetAveragePeriod() != 10.0 {
		t.Errorf("Expected average period 10.0, got %f", dummy.GetAveragePeriod())
	}

	// Test invalid period
	err = dummy.SetAveragePeriod(0)
	if err != ErrInvalidPeriod {
		t.Errorf("Expected error %v, got %v", ErrInvalidPeriod, err)
	}

	err = dummy.SetAveragePeriod(-1.0)
	if err != ErrInvalidPeriod {
		t.Errorf("Expected error %v, got %v", ErrInvalidPeriod, err)
	}
}

func TestObservingConditionsDummy_GetCloudCover(t *testing.T) {
	dummy := NewObservingConditionsDummy("test", "Test", "Test Station")

	// Test default value
	if dummy.GetCloudCover() != 0 {
		t.Errorf("Expected default cloud cover 0, got %f", dummy.GetCloudCover())
	}

	// Set a value and test
	dummy.condition.CloudCover = 25.5
	if dummy.GetCloudCover() != 25.5 {
		t.Errorf("Expected cloud cover 25.5, got %f", dummy.GetCloudCover())
	}
}

func TestObservingConditionsDummy_GetDewPoint(t *testing.T) {
	dummy := NewObservingConditionsDummy("test", "Test", "Test Station")

	// Test default value
	if dummy.GetDewPoint() != 0 {
		t.Errorf("Expected default dew point 0, got %f", dummy.GetDewPoint())
	}

	// Set a value and test
	dummy.condition.DewPoint = 12.3
	if dummy.GetDewPoint() != 12.3 {
		t.Errorf("Expected dew point 12.3, got %f", dummy.GetDewPoint())
	}
}

func TestObservingConditionsDummy_GetHumidity(t *testing.T) {
	dummy := NewObservingConditionsDummy("test", "Test", "Test Station")

	// Test default value
	if dummy.GetHumidity() != 0 {
		t.Errorf("Expected default humidity 0, got %f", dummy.GetHumidity())
	}

	// Set a value and test
	dummy.condition.Humidity = 75.2
	if dummy.GetHumidity() != 75.2 {
		t.Errorf("Expected humidity 75.2, got %f", dummy.GetHumidity())
	}
}

func TestObservingConditionsDummy_GetPressure(t *testing.T) {
	dummy := NewObservingConditionsDummy("test", "Test", "Test Station")

	// Test default value
	if dummy.GetPressure() != 0 {
		t.Errorf("Expected default pressure 0, got %f", dummy.GetPressure())
	}

	// Set a value and test
	dummy.condition.Pressure = 1015.8
	if dummy.GetPressure() != 1015.8 {
		t.Errorf("Expected pressure 1015.8, got %f", dummy.GetPressure())
	}
}

func TestObservingConditionsDummy_GetRainRate(t *testing.T) {
	dummy := NewObservingConditionsDummy("test", "Test", "Test Station")

	// Test default value
	if dummy.GetRainRate() != 0 {
		t.Errorf("Expected default rain rate 0, got %f", dummy.GetRainRate())
	}

	// Set a value and test
	dummy.condition.RainRate = 2.5
	if dummy.GetRainRate() != 2.5 {
		t.Errorf("Expected rain rate 2.5, got %f", dummy.GetRainRate())
	}
}

func TestObservingConditionsDummy_GetSkyBrightness(t *testing.T) {
	dummy := NewObservingConditionsDummy("test", "Test", "Test Station")

	// Test default value
	if dummy.GetSkyBrightness() != 0 {
		t.Errorf("Expected default sky brightness 0, got %f", dummy.GetSkyBrightness())
	}

	// Set a value and test
	dummy.condition.SkyBrightness = 18.5
	if dummy.GetSkyBrightness() != 18.5 {
		t.Errorf("Expected sky brightness 18.5, got %f", dummy.GetSkyBrightness())
	}
}

func TestObservingConditionsDummy_GetSkyQuality(t *testing.T) {
	dummy := NewObservingConditionsDummy("test", "Test", "Test Station")

	// Test default value
	if dummy.GetSkyQuality() != 0 {
		t.Errorf("Expected default sky quality 0, got %f", dummy.GetSkyQuality())
	}

	// Set a value and test
	dummy.condition.SkyQuality = 21.2
	if dummy.GetSkyQuality() != 21.2 {
		t.Errorf("Expected sky quality 21.2, got %f", dummy.GetSkyQuality())
	}
}

func TestObservingConditionsDummy_GetSkyTemperature(t *testing.T) {
	dummy := NewObservingConditionsDummy("test", "Test", "Test Station")

	// Test default value
	if dummy.GetSkyTemperature() != 0 {
		t.Errorf("Expected default sky temperature 0, got %f", dummy.GetSkyTemperature())
	}

	// Set a value and test
	dummy.condition.SkyTemperature = -15.7
	if dummy.GetSkyTemperature() != -15.7 {
		t.Errorf("Expected sky temperature -15.7, got %f", dummy.GetSkyTemperature())
	}
}

func TestObservingConditionsDummy_GetStarFWHM(t *testing.T) {
	dummy := NewObservingConditionsDummy("test", "Test", "Test Station")

	// Test default value
	if dummy.GetStarFWHM() != 0 {
		t.Errorf("Expected default star FWHM 0, got %f", dummy.GetStarFWHM())
	}

	// Set a value and test
	dummy.condition.StarFWHM = 3.2
	if dummy.GetStarFWHM() != 3.2 {
		t.Errorf("Expected star FWHM 3.2, got %f", dummy.GetStarFWHM())
	}
}

func TestObservingConditionsDummy_GetTemperature(t *testing.T) {
	dummy := NewObservingConditionsDummy("test", "Test", "Test Station")

	// Test default value
	if dummy.GetTemperature() != 0 {
		t.Errorf("Expected default temperature 0, got %f", dummy.GetTemperature())
	}

	// Set a value and test
	dummy.condition.Temperature = 22.8
	if dummy.GetTemperature() != 22.8 {
		t.Errorf("Expected temperature 22.8, got %f", dummy.GetTemperature())
	}
}

func TestObservingConditionsDummy_GetWindDirection(t *testing.T) {
	dummy := NewObservingConditionsDummy("test", "Test", "Test Station")

	// Test default value
	if dummy.GetWindDirection() != 0 {
		t.Errorf("Expected default wind direction 0, got %f", dummy.GetWindDirection())
	}

	// Set a value and test
	dummy.condition.WindDirection = 180.0
	if dummy.GetWindDirection() != 180.0 {
		t.Errorf("Expected wind direction 180.0, got %f", dummy.GetWindDirection())
	}
}

func TestObservingConditionsDummy_GetWindGust(t *testing.T) {
	dummy := NewObservingConditionsDummy("test", "Test", "Test Station")

	// Test default value
	if dummy.GetWindGust() != 0 {
		t.Errorf("Expected default wind gust 0, got %f", dummy.GetWindGust())
	}

	// Set a value and test
	dummy.condition.WindGust = 15.3
	if dummy.GetWindGust() != 15.3 {
		t.Errorf("Expected wind gust 15.3, got %f", dummy.GetWindGust())
	}
}

func TestObservingConditionsDummy_GetWindSpeed(t *testing.T) {
	dummy := NewObservingConditionsDummy("test", "Test", "Test Station")

	// Test default value
	if dummy.GetWindSpeed() != 0 {
		t.Errorf("Expected default wind speed 0, got %f", dummy.GetWindSpeed())
	}

	// Set a value and test
	dummy.condition.WindSpeed = 8.7
	if dummy.GetWindSpeed() != 8.7 {
		t.Errorf("Expected wind speed 8.7, got %f", dummy.GetWindSpeed())
	}
}

func TestObservingConditionsDummy_GetTimeSinceLastUpdate(t *testing.T) {
	dummy := NewObservingConditionsDummy("test", "Test", "Test Station")

	// Set a known last refresh time
	dummy.lastRefreshTime = time.Now().Add(-5 * time.Second)

	timeSince := dummy.GetTimeSinceLastUpdate()

	// Allow for some small variation in timing
	if timeSince < 4.9 || timeSince > 5.1 {
		t.Errorf("Expected time since last update to be approximately 5.0 seconds, got %f", timeSince)
	}
}

func TestNewObservingConditionsHttp(t *testing.T) {
	id := "http1"
	name := "Test HTTP"
	desc := "Test HTTP Weather Station"
	url := "http://example.com/weather"

	// Test valid URL
	http, err := NewObservingConditionsHttp(id, name, desc, url)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if http.GetId() != id {
		t.Errorf("Expected ID %s, got %s", id, http.GetId())
	}
	if http.GetName() != name {
		t.Errorf("Expected name %s, got %s", name, http.GetName())
	}
	if http.GetDescription() != desc {
		t.Errorf("Expected description %s, got %s", desc, http.GetDescription())
	}

	// Test invalid URL
	http, err = NewObservingConditionsHttp(id, name, desc, "")
	if err != ErrInvalidURL {
		t.Errorf("Expected error %v, got %v", ErrInvalidURL, err)
	}
	if http != nil {
		t.Errorf("Expected nil HTTP client, got %v", http)
	}
}

func TestObservingConditionsHttp_Refresh(t *testing.T) {
	// Create a test server that returns valid weather data
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		weatherData := `{
			"id": 1,
			"indoortemp": 22.5,
			"temp": 18.3,
			"dewpt": 12.1,
			"windchill": 16.8,
			"indoorhumidity": 45,
			"humidity": 65,
			"windspeedms": 5.2,
			"windgustms": 8.1,
			"winddir": 180,
			"absbaromin": 1013.25,
			"baromin": 1013.25,
			"rainin": 0.0,
			"dailyrainin": 2.5,
			"weeklyrainin": 15.3,
			"monthlyrainin": 45.7,
			"solarradiation": 450.2,
			"UV": 3,
			"dateutc": "2023-01-01T12:00:00Z",
			"softwaretype": "test"
		}`
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(weatherData))
	}))
	defer server.Close()

	http, err := NewObservingConditionsHttp("test", "Test", "Test Station", server.URL)
	if err != nil {
		t.Fatalf("Failed to create HTTP client: %v", err)
	}

	// Test refresh
	err = http.Refresh()
	if err != nil {
		t.Errorf("Expected no error from Refresh(), got %v", err)
	}

	// Verify the data was parsed correctly
	if http.GetTemperature() != 18.3 {
		t.Errorf("Expected temperature 18.3, got %f", http.GetTemperature())
	}
	if http.GetDewPoint() != 12.1 {
		t.Errorf("Expected dew point 12.1, got %f", http.GetDewPoint())
	}
	if http.GetHumidity() != 65.0 {
		t.Errorf("Expected humidity 65.0, got %f", http.GetHumidity())
	}
	if http.GetPressure() != 1013.25 {
		t.Errorf("Expected pressure 1013.25, got %f", http.GetPressure())
	}
	if http.GetWindSpeed() != 5.2 {
		t.Errorf("Expected wind speed 5.2, got %f", http.GetWindSpeed())
	}
	if http.GetWindGust() != 8.1 {
		t.Errorf("Expected wind gust 8.1, got %f", http.GetWindGust())
	}
	if http.GetWindDirection() != 180.0 {
		t.Errorf("Expected wind direction 180.0, got %f", http.GetWindDirection())
	}
	if http.GetRainRate() != 0.0 {
		t.Errorf("Expected rain rate 0.0, got %f", http.GetRainRate())
	}
}

func TestObservingConditionsHttp_Refresh_InvalidResponse(t *testing.T) {
	// Create a test server that returns invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	http, err := NewObservingConditionsHttp("test", "Test", "Test Station", server.URL)
	if err != nil {
		t.Fatalf("Failed to create HTTP client: %v", err)
	}

	// Test refresh with invalid JSON
	err = http.Refresh()
	if err == nil {
		t.Error("Expected error from Refresh() with invalid JSON, got nil")
	}
}

func TestObservingConditionsHttp_Refresh_ServerError(t *testing.T) {
	// Create a test server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	http, err := NewObservingConditionsHttp("test", "Test", "Test Station", server.URL)
	if err != nil {
		t.Fatalf("Failed to create HTTP client: %v", err)
	}

	// Test refresh with server error
	err = http.Refresh()
	if err == nil {
		t.Error("Expected error from Refresh() with server error, got nil")
	}
}

func TestObservingConditionsHttp_GetId(t *testing.T) {
	http, _ := NewObservingConditionsHttp("test-id", "Test", "Test Station", "http://example.com")

	expected := "test-id"
	actual := http.GetId()
	if actual != expected {
		t.Errorf("Expected ID %s, got %s", expected, actual)
	}
}

func TestObservingConditionsHttp_GetName(t *testing.T) {
	http, _ := NewObservingConditionsHttp("test", "Test Name", "Test Station", "http://example.com")

	expected := "Test Name"
	actual := http.GetName()
	if actual != expected {
		t.Errorf("Expected name %s, got %s", expected, actual)
	}
}

func TestObservingConditionsHttp_GetDescription(t *testing.T) {
	http, _ := NewObservingConditionsHttp("test", "Test", "Test Description", "http://example.com")

	expected := "Test Description"
	actual := http.GetDescription()
	if actual != expected {
		t.Errorf("Expected description %s, got %s", expected, actual)
	}
}

func TestObservingConditionsHttp_GetState(t *testing.T) {
	http, _ := NewObservingConditionsHttp("test", "Test", "Test Station", "http://example.com")

	// Set some test values
	http.condition = WeatherCondition{
		Temperature: 20.5,
		Humidity:    65.0,
		Pressure:    1013.25,
	}

	state := http.GetState()

	// Verify it's valid JSON
	var parsed WeatherCondition
	err := json.Unmarshal([]byte(state), &parsed)
	if err != nil {
		t.Errorf("GetState() should return valid JSON, got error: %v", err)
	}

	if parsed.Temperature != 20.5 {
		t.Errorf("Expected temperature 20.5, got %f", parsed.Temperature)
	}
}

func TestObservingConditionsHttp_GetAveragePeriod(t *testing.T) {
	http, _ := NewObservingConditionsHttp("test", "Test", "Test Station", "http://example.com")

	// Test default value
	if http.GetAveragePeriod() != 0 {
		t.Errorf("Expected default average period 0, got %f", http.GetAveragePeriod())
	}

	// Set a value and test
	http.condition.AveragePeriod = 15.5
	if http.GetAveragePeriod() != 15.5 {
		t.Errorf("Expected average period 15.5, got %f", http.GetAveragePeriod())
	}
}

func TestObservingConditionsHttp_SetAveragePeriod(t *testing.T) {
	http, _ := NewObservingConditionsHttp("test", "Test", "Test Station", "http://example.com")

	// Test valid period
	err := http.SetAveragePeriod(10.0)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if http.GetAveragePeriod() != 10.0 {
		t.Errorf("Expected average period 10.0, got %f", http.GetAveragePeriod())
	}

	// Test invalid period
	err = http.SetAveragePeriod(0)
	if err != ErrInvalidPeriod {
		t.Errorf("Expected error %v, got %v", ErrInvalidPeriod, err)
	}

	err = http.SetAveragePeriod(-1.0)
	if err != ErrInvalidPeriod {
		t.Errorf("Expected error %v, got %v", ErrInvalidPeriod, err)
	}
}

func TestObservingConditionsHttp_GetCloudCover(t *testing.T) {
	http, _ := NewObservingConditionsHttp("test", "Test", "Test Station", "http://example.com")

	// Test default value
	if http.GetCloudCover() != 0 {
		t.Errorf("Expected default cloud cover 0, got %f", http.GetCloudCover())
	}

	// Set a value and test
	http.condition.CloudCover = 25.5
	if http.GetCloudCover() != 25.5 {
		t.Errorf("Expected cloud cover 25.5, got %f", http.GetCloudCover())
	}
}

func TestObservingConditionsHttp_GetDewPoint(t *testing.T) {
	http, _ := NewObservingConditionsHttp("test", "Test", "Test Station", "http://example.com")

	// Test default value
	if http.GetDewPoint() != 0 {
		t.Errorf("Expected default dew point 0, got %f", http.GetDewPoint())
	}

	// Set a value and test
	http.condition.DewPoint = 12.3
	if http.GetDewPoint() != 12.3 {
		t.Errorf("Expected dew point 12.3, got %f", http.GetDewPoint())
	}
}

func TestObservingConditionsHttp_GetHumidity(t *testing.T) {
	http, _ := NewObservingConditionsHttp("test", "Test", "Test Station", "http://example.com")

	// Test default value
	if http.GetHumidity() != 0 {
		t.Errorf("Expected default humidity 0, got %f", http.GetHumidity())
	}

	// Set a value and test
	http.condition.Humidity = 75.2
	if http.GetHumidity() != 75.2 {
		t.Errorf("Expected humidity 75.2, got %f", http.GetHumidity())
	}
}

func TestObservingConditionsHttp_GetPressure(t *testing.T) {
	http, _ := NewObservingConditionsHttp("test", "Test", "Test Station", "http://example.com")

	// Test default value
	if http.GetPressure() != 0 {
		t.Errorf("Expected default pressure 0, got %f", http.GetPressure())
	}

	// Set a value and test
	http.condition.Pressure = 1015.8
	if http.GetPressure() != 1015.8 {
		t.Errorf("Expected pressure 1015.8, got %f", http.GetPressure())
	}
}

func TestObservingConditionsHttp_GetRainRate(t *testing.T) {
	http, _ := NewObservingConditionsHttp("test", "Test", "Test Station", "http://example.com")

	// Test default value
	if http.GetRainRate() != 0 {
		t.Errorf("Expected default rain rate 0, got %f", http.GetRainRate())
	}

	// Set a value and test
	http.condition.RainRate = 2.5
	if http.GetRainRate() != 2.5 {
		t.Errorf("Expected rain rate 2.5, got %f", http.GetRainRate())
	}
}

func TestObservingConditionsHttp_GetSkyBrightness(t *testing.T) {
	http, _ := NewObservingConditionsHttp("test", "Test", "Test Station", "http://example.com")

	// Test default value
	if http.GetSkyBrightness() != 0 {
		t.Errorf("Expected default sky brightness 0, got %f", http.GetSkyBrightness())
	}

	// Set a value and test
	http.condition.SkyBrightness = 18.5
	if http.GetSkyBrightness() != 18.5 {
		t.Errorf("Expected sky brightness 18.5, got %f", http.GetSkyBrightness())
	}
}

func TestObservingConditionsHttp_GetSkyQuality(t *testing.T) {
	http, _ := NewObservingConditionsHttp("test", "Test", "Test Station", "http://example.com")

	// Test default value
	if http.GetSkyQuality() != 0 {
		t.Errorf("Expected default sky quality 0, got %f", http.GetSkyQuality())
	}

	// Set a value and test
	http.condition.SkyQuality = 21.2
	if http.GetSkyQuality() != 21.2 {
		t.Errorf("Expected sky quality 21.2, got %f", http.GetSkyQuality())
	}
}

func TestObservingConditionsHttp_GetSkyTemperature(t *testing.T) {
	http, _ := NewObservingConditionsHttp("test", "Test", "Test Station", "http://example.com")

	// Test default value
	if http.GetSkyTemperature() != 0 {
		t.Errorf("Expected default sky temperature 0, got %f", http.GetSkyTemperature())
	}

	// Set a value and test
	http.condition.SkyTemperature = -15.7
	if http.GetSkyTemperature() != -15.7 {
		t.Errorf("Expected sky temperature -15.7, got %f", http.GetSkyTemperature())
	}
}

func TestObservingConditionsHttp_GetStarFWHM(t *testing.T) {
	http, _ := NewObservingConditionsHttp("test", "Test", "Test Station", "http://example.com")

	// Test default value
	if http.GetStarFWHM() != 0 {
		t.Errorf("Expected default star FWHM 0, got %f", http.GetStarFWHM())
	}

	// Set a value and test
	http.condition.StarFWHM = 3.2
	if http.GetStarFWHM() != 3.2 {
		t.Errorf("Expected star FWHM 3.2, got %f", http.GetStarFWHM())
	}
}

func TestObservingConditionsHttp_GetTemperature(t *testing.T) {
	http, _ := NewObservingConditionsHttp("test", "Test", "Test Station", "http://example.com")

	// Test default value
	if http.GetTemperature() != 0 {
		t.Errorf("Expected default temperature 0, got %f", http.GetTemperature())
	}

	// Set a value and test
	http.condition.Temperature = 22.8
	if http.GetTemperature() != 22.8 {
		t.Errorf("Expected temperature 22.8, got %f", http.GetTemperature())
	}
}

func TestObservingConditionsHttp_GetWindDirection(t *testing.T) {
	http, _ := NewObservingConditionsHttp("test", "Test", "Test Station", "http://example.com")

	// Test default value
	if http.GetWindDirection() != 0 {
		t.Errorf("Expected default wind direction 0, got %f", http.GetWindDirection())
	}

	// Set a value and test
	http.condition.WindDirection = 180.0
	if http.GetWindDirection() != 180.0 {
		t.Errorf("Expected wind direction 180.0, got %f", http.GetWindDirection())
	}
}

func TestObservingConditionsHttp_GetWindGust(t *testing.T) {
	http, _ := NewObservingConditionsHttp("test", "Test", "Test Station", "http://example.com")

	// Test default value
	if http.GetWindGust() != 0 {
		t.Errorf("Expected default wind gust 0, got %f", http.GetWindGust())
	}

	// Set a value and test
	http.condition.WindGust = 15.3
	if http.GetWindGust() != 15.3 {
		t.Errorf("Expected wind gust 15.3, got %f", http.GetWindGust())
	}
}

func TestObservingConditionsHttp_GetWindSpeed(t *testing.T) {
	http, _ := NewObservingConditionsHttp("test", "Test", "Test Station", "http://example.com")

	// Test default value
	if http.GetWindSpeed() != 0 {
		t.Errorf("Expected default wind speed 0, got %f", http.GetWindSpeed())
	}

	// Set a value and test
	http.condition.WindSpeed = 8.7
	if http.GetWindSpeed() != 8.7 {
		t.Errorf("Expected wind speed 8.7, got %f", http.GetWindSpeed())
	}
}

func TestObservingConditionsHttp_GetTimeSinceLastUpdate(t *testing.T) {
	http, _ := NewObservingConditionsHttp("test", "Test", "Test Station", "http://example.com")
	result := http.GetTimeSinceLastUpdate()
	if result < 0 {
		t.Errorf("Expected non-negative time since last update, got %f", result)
	}
}

func TestIsSensorAvailable(t *testing.T) {
	// Test available sensors
	availableSensors := []string{
		SensorAveragePeriod,
		SensorCloudCover,
		SensorDewPoint,
		SensorHumidity,
		SensorPressure,
		SensorRainRate,
		SensorWindDirection,
		SensorWindGust,
		SensorWindSpeed,
	}

	for _, sensor := range availableSensors {
		if !IsSensorAvailable(sensor) {
			t.Errorf("Expected sensor %s to be available", sensor)
		}
	}

	// Test unavailable sensors
	unavailableSensors := []string{
		SensorSkyBrightness,
		SensorSkyQuality,
		SensorSkyTemperature,
		SensorStarFWHM,
		SensorTemperature,
	}

	for _, sensor := range unavailableSensors {
		if IsSensorAvailable(sensor) {
			t.Errorf("Expected sensor %s to be unavailable", sensor)
		}
	}

	// Test invalid sensor
	if IsSensorAvailable("InvalidSensor") {
		t.Error("Expected invalid sensor to be unavailable")
	}
}
