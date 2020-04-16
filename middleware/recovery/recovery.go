package recovery

import (
	"fmt"
	"runtime/debug"

	"github.com/PolarPanda611/trinitygo/application"
	"github.com/PolarPanda611/trinitygo/httputil"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
)

// New runtime middleware
func New(app application.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// when stack finishes
				logMessage := fmt.Sprintf("Recovered from HTTP Request %v %v \n", c.Request.Method, c.Request.URL)
				logMessage += fmt.Sprintf("Trace: %s\n", err)
				logMessage += fmt.Sprintf("\n%s", debug.Stack())
				app.Logger().Warn(logMessage)
				c.AbortWithStatusJSON(400, httputil.ResponseData{
					Status: 400,
					Error: map[string]string{
						"code":    codes.Internal.String(),
						"message": fmt.Sprintf("Internal err : %v", err),
					},
				})
				return
			}

		}()
		c.Next()
	}
}
