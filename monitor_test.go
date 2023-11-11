package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestSafetyMonitorDummy(t *testing.T) {
	var dummy = NewSafetyMonitorDummy("dummy", "name", "description", true)
	assert.Equal(t, true, dummy.IsSafe(), "they should be equal")
	dummy.safe = false
	assert.Equal(t, false, dummy.IsSafe(), "they should be equal")
	assert.Equal(t, "name", dummy.GetName(), "they should be equal")
	assert.Equal(t, "description", dummy.GetDescription(), "they should be equal")

}

func TestSafetyMonitorFile_IsSafeRefresh(t *testing.T) {
	f, err := os.CreateTemp("", "SafetyMonitorFileTest")
	if err != nil {
		panic(err)
	}
	defer os.Remove(f.Name())

	fmt.Println("Temp file name:", f.Name())
	_, err = f.Write([]byte("TrUe"))
	if err != nil {
		panic(err)
	}
	var file = NewSafetyMonitorFile("file", "name", "description", f.Name())
	assert.Equal(t, true, file.IsSafe(), "they should be equal")
	f.Truncate(0)

	_, err = f.Write([]byte("false"))
	file.Refresh()
	assert.Equal(t, false, file.IsSafe(), "they should be equal")
}

func TestSafetyMonitorFile_InvalidPath(t *testing.T) {
	f, err := os.CreateTemp("", "SafetyMonitorFileTest")
	if err != nil {
		panic(err)
	}
	defer os.Remove(f.Name())

	fmt.Println("Temp file name:", f.Name())
	_, err = f.Write([]byte("TrUe"))
	if err != nil {
		panic(err)
	}
	var file = NewSafetyMonitorFile("invalid", "name", "description", f.Name()+"invalid")
	assert.Equal(t, false, file.IsSafe(), "they should be equal")
}

func startHttpServer(content string) *gin.Engine {
	router := gin.Default()
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, content)
	})
	return router
}

func TestSafetyMonitorHttp_IsSafe(t *testing.T) {
	r := startHttpServer("TrUe")
	go r.Run("127.0.0.1:12345")
	time.Sleep(1 * time.Second)
	var httpsm = NewSafetyMonitorHttp("id", "name", "description", "http://127.0.0.1:12345/test")
	assert.Equal(t, true, httpsm.IsSafe(), "they should be equal")
}

func TestSafetyMonitorHttp_InvalidUrl(t *testing.T) {
	var file = NewSafetyMonitorHttp("id", "name", "description", "http://127.0.0.1/test")
	assert.Equal(t, false, file.IsSafe(), "they should be equal")
}
