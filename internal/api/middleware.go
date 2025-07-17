package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ValidationContext holds validated request data
type ValidationContext struct {
	DeviceID            int
	ClientID            int
	ClientTransactionID int
	FullClientID        ClientId
	IsValid             bool
	ErrorMessage        string
}

// alpacaValidationMiddleware validates Alpaca API requests
func alpacaValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract and validate device_id
		deviceIDStr := c.Param("device_id")
		deviceID, err := strconv.Atoi(deviceIDStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"Error": "Invalid device_id parameter",
			})
			return
		}

		// Extract and validate ClientID
		clientID := getClientId(c)
		if clientID < 1 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"Error": "Invalid or missing ClientID parameter",
			})
			return
		}

		// Extract and validate ClientTransactionID
		clientTransactionID := getClientTransactionId(c)
		if clientTransactionID < 1 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"Error": "Invalid or missing ClientTransactionID parameter",
			})
			return
		}

		// Create validation context
		validationCtx := &ValidationContext{
			DeviceID:            deviceID,
			ClientID:            clientID,
			ClientTransactionID: clientTransactionID,
			FullClientID:        getFullClientId(c),
			IsValid:             true,
		}

		// Store validation context in gin context for handlers to access
		c.Set("validation", validationCtx)
		c.Next()
	}
}

// GetValidationContext retrieves the validation context from gin context
func GetValidationContext(c *gin.Context) *ValidationContext {
	if v, exists := c.Get("validation"); exists {
		if validationCtx, ok := v.(*ValidationContext); ok {
			return validationCtx
		}
	}
	return nil
}

// deviceValidationMiddleware validates that the device exists and is accessible
func deviceValidationMiddleware(deviceType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		validationCtx := GetValidationContext(c)
		if validationCtx == nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"Error": "Validation context not found",
			})
			return
		}

		// Get the API server from context (we'll need to set this up)
		apiServer, exists := c.Get("apiServer")
		if !exists {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"Error": "API server not found in context",
			})
			return
		}

		server := apiServer.(*ApiServer)

		// Check if device exists
		devices, exists := server.Devices[deviceType]
		if !exists {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"Error": "Device type not supported",
			})
			return
		}

		device, exists := devices[validationCtx.DeviceID]
		if !exists {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"Error": "Device not found",
			})
			return
		}

		// Store device in context for handlers
		c.Set("device", device)
		c.Next()
	}
}

// alpacaResponseMiddleware prepares the standard Alpaca response
func alpacaResponseMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		validationCtx := GetValidationContext(c)
		if validationCtx == nil {
			c.Next()
			return
		}

		// Store response preparation function in context
		c.Set("prepareResponse", func(resp *alpacaResponse) {
			ctid := validationCtx.ClientTransactionID
			if ctid < 0 {
				ctid = 0
			}

			// Get API server from context
			if apiServer, exists := c.Get("apiServer"); exists {
				if server, ok := apiServer.(*ApiServer); ok {
					server.ServerTransactionID += 1
					resp.ClientTransactionID = uint32(ctid)
					resp.ServerTransactionID = server.ServerTransactionID
				}
			}
		})

		c.Next()
	}
}
