package runtime

import (
	"fmt"

	"github.com/PolarPanda611/trinitygo/application"
	"github.com/PolarPanda611/trinitygo/httputils"
	"github.com/gin-gonic/gin"
)

// DefaultHeaderPrefix will be used in set header as prefix
var DefaultHeaderPrefix = "trinity_"

// New runtime middleware
func New(app application.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, v := range app.RuntimeKeys() {
			keyValue := c.GetHeader(v.GetKeyName())
			if keyValue == "" {
				if v.GetRequired() {
					c.AbortWithStatusJSON(400, httputils.ResponseData{
						Status: 400,
						Result: fmt.Sprintf("runtime key %v is required ", v.GetKeyName()),
					})
					return
				}
				c.Set(v.GetKeyName(), v.GetDefaultValue())

				c.Header(fmt.Sprintf("%v%v", DefaultHeaderPrefix, v.GetKeyName()), v.GetDefaultValue())
			}
		}
		c.Next()
	}
}
