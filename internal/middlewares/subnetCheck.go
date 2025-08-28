package middlewares

import (
	// "bytes"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stepanov-ds/ya-metrics/internal/config/server"
)

func SubnetCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		if *server.TrustedSubnet == "" {
			c.Next()
			return
		}

		xRealIP := c.GetHeader("X-Real-IP") 
		
		ip:= net.ParseIP(xRealIP)

		if !server.TrustedSubnetObj.Contains(ip) && !ip.IsLoopback() {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()

	}
	
}

// func compareIPNet(n1, n2 *net.IPNet) bool {
// 	if n1 == nil && n2 == nil {
// 		return true 
// 	}
// 	if n1 == nil || n2 == nil {
// 		return false 
// 	}
// 	return n1.IP.Equal(n2.IP) && bytes.Equal(n1.Mask, n2.Mask)
// }
