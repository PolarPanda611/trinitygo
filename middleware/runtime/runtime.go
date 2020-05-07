package runtime

import (
	"fmt"

	"github.com/PolarPanda611/trinitygo/application"
	"github.com/PolarPanda611/trinitygo/httputil"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
)

// New runtime middleware
func New(app application.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, v := range app.RuntimeKeys() {
			keyValue := c.GetHeader(v.GetKeyName())
			if keyValue == "" {
				if v.GetRequired() {
					c.AbortWithStatusJSON(400, httputil.ResponseData{
						Status: 400,
						Err: map[string]string{
							"code":    codes.Internal.String(),
							"message": fmt.Sprintf("runtime key %v is required ", v.GetKeyName()),
						},
					})
					return
				}
				c.Set(v.GetKeyName(), v.GetDefaultValue())
				if v.IsLog() {
					c.Header(v.GetKeyName(), v.GetDefaultValue())
				}
			}
		}
		c.Next()
	}
}
