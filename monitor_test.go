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
	assert.NoError(t, err, "should work")
	defer os.Remove(f.Name())

	fmt.Println("Temp file name:", f.Name())
	_, err = f.Write([]byte("TrUe"))
	assert.NoError(t, err, "should work")
	var file = NewSafetyMonitorFile("file", "name", "description", f.Name(), NewSafetyMatchingRule(false, ""))
	assert.Equal(t, "name", file.GetName(), "they should be equal")
	assert.Equal(t, "description", file.GetDescription(), "they should be equal")
	assert.Equal(t, true, file.IsSafe(), "they should be equal")
	assert.Equal(t, "TrUe", file.GetRawValue(), "they should be equal")
	f.Truncate(0)
	f.Seek(0, 0)
	_, err = f.Write([]byte("1"))
	assert.NoError(t, err, "should work")
	file.Refresh()
	assert.Equal(t, true, file.IsSafe(), "they should be equal")
	f.Truncate(0)
	f.Seek(0, 0)
	_, err = f.Write([]byte("false"))
	assert.NoError(t, err, "should work")
	file.Refresh()
	assert.Equal(t, false, file.IsSafe(), "they should be equal")
	f.Truncate(0)
	f.Seek(0, 0)
	_, err = f.Write([]byte("0"))
	assert.NoError(t, err, "should work")
	file.Refresh()
	assert.Equal(t, false, file.IsSafe(), "they should be equal")
	file = NewSafetyMonitorFile("file", "name", "description", f.Name(), NewSafetyMatchingRule(false, "open"))
	f.Truncate(0)
	f.Seek(0, 0)
	_, err = f.Write([]byte("open"))
	assert.NoError(t, err, "should work")
	file.Refresh()
	assert.Equal(t, true, file.IsSafe(), "they should be equal")
	f.Truncate(0)
	f.Seek(0, 0)
	_, err = f.Write([]byte("closing"))
	assert.NoError(t, err, "should work")
	file.Refresh()
	assert.Equal(t, false, file.IsSafe(), "they should be equal")
}

func TestSafetyMonitorFile_InvalidPath(t *testing.T) {
	f, err := os.CreateTemp("", "SafetyMonitorFileTest")
	assert.NoError(t, err, "should work")
	defer os.Remove(f.Name())

	fmt.Println("Temp file name:", f.Name())
	_, err = f.Write([]byte("TrUe"))
	assert.NoError(t, err, "should work")
	var file = NewSafetyMonitorFile("invalid", "name", "description", f.Name()+"invalid", NewSafetyMatchingRule(false, ""))
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
	var httpsm = NewSafetyMonitorHttp("id", "name", "description", "http://127.0.0.1:12345/test", NewSafetyMatchingRule(false, ""))
	assert.Equal(t, true, httpsm.IsSafe(), "they should be equal")
	assert.Equal(t, "TrUe", httpsm.GetRawValue(), "they should be equal")
}

func TestSafetyMonitorHttp_IsUnsafe(t *testing.T) {
	r := startHttpServer("0")
	go r.Run("127.0.0.1:12346")
	time.Sleep(1 * time.Second)
	var httpsm = NewSafetyMonitorHttp("id", "name", "description", "http://127.0.0.1:12346/test", NewSafetyMatchingRule(false, ""))
	assert.Equal(t, false, httpsm.IsSafe(), "they should be equal")
	assert.Equal(t, "0", httpsm.GetRawValue(), "they should be equal")
}

func TestSafetyMonitorHttp_InvalidUrl(t *testing.T) {
	var file = NewSafetyMonitorHttp("id", "name", "description", "http://127.0.0.1/test", NewSafetyMatchingRule(false, ""))
	assert.Equal(t, false, file.IsSafe(), "they should be equal")
}

func TestSafetyMatchingRule_Regex(t *testing.T) {
	rule := NewSafetyMatchingRule(false, "[a-z]+")
	assert.Equal(t, false, rule.isSafe("1"), "they should be equal")
	assert.Equal(t, false, rule.isSafe("0"), "they should be equal")
	assert.Equal(t, true, rule.isSafe("abc"), "they should be equal")
}

func TestSafetyMatchingRule_InvertRegex(t *testing.T) {
	rule := NewSafetyMatchingRule(true, "[a-z]+")
	assert.Equal(t, false, rule.isSafe("1"), "they should be equal")
	assert.Equal(t, false, rule.isSafe("0"), "they should be equal")
	assert.Equal(t, true, rule.isSafe("abc"), "they should be equal")
}

func TestSafetyMatchingRule_InvalidRegex(t *testing.T) {
	rule := NewSafetyMatchingRule(false, "")
	assert.Equal(t, true, rule.isSafe("1"), "they should be equal")
	assert.Equal(t, true, rule.isSafe("TrUe"), "they should be equal")
	assert.Equal(t, false, rule.isSafe("0"), "they should be equal")
	assert.Equal(t, false, rule.isSafe("abc"), "they should be equal")
}
