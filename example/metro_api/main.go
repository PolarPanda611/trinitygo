package main

import (
	"fmt"
	"os"
	"strings"

	_ "metro_api/domain/controller/http"
	"metro_api/domain/model"
	"metro_api/infra/db"
	"metro_api/infra/migrate"

	"github.com/PolarPanda611/trinitygo/keyword"

	_ "metro_api/docs"

	"github.com/PolarPanda611/trinitygo"
	tjwt "github.com/PolarPanda611/trinitygo/middleware/jwt"
	truntime "github.com/PolarPanda611/trinitygo/runtime"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @title metro_api
// @version 1.0
// @description  metro_api
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 127.0.0.1:8088
// @BasePath /metro_api/
func main() {
	currentPath, _ := os.Getwd()
	configPath := fmt.Sprintf(currentPath + "/conf/conf.toml")
	trinitygo.SetConfigPath(configPath)
	trinitygo.SetResponseFactory(CustomizeResponseFactory)
	trinitygo.SetKeyword(keyword.Keyword{
		SearchBy:      "SearchBy",
		PageNum:       "current",
		PageSize:      "pageSize",
		OrderBy:       "OrderBy",
		PaginationOff: "PaginationOff",
	})
	trinitygo.EnableHealthCheckURL()
	t := trinitygo.DefaultHTTP()
	jwtConfig := &tjwt.Config{
		HeaderName:         "Authorization",
		SecretKey:          "test",
		JwtPrefix:          "Bearer",
		Issuer:             "xxx",
		ExpireHour:         2,
		IsVerifyIssuer:     false,
		IsVerifyExpireHour: false,
		IsVerifySecretKey:  false,
		Claims:             &Claims{},
		SuccessCallback: func(c *gin.Context, claim jwt.Claims) {
			dktClaim, _ := claim.(*Claims)
			userNameLowerCase := strings.ToLower(dktClaim.UID)
			c.Set("user_name", userNameLowerCase)
			var user model.User
			if err := db.DB.Where("user_name = ? ", userNameLowerCase).First(&user).Error; err != nil {
				c.AbortWithError(400, err)
				return
			}
			c.Set(fmt.Sprintf("%v", "user_name"), userNameLowerCase)
			c.Set(fmt.Sprintf("%v", "user_id"), fmt.Sprintf("%v", user.ID))
		},
		MethodWhiteList: []string{"OPTIONS"},
		PathWhiteList:   []string{"/metro_api/ping", "/metro_api/v1/login"},
	}
	t.UseMiddleware(tjwt.New(t, jwtConfig))
	t.RegRuntimeKey(truntime.NewRuntimeKey("trace_id", false, func() string { return uuid.New().String() }, true))
	t.RegRuntimeKey(truntime.NewRuntimeKey("user_id", false, func() string { return "" }, false))
	t.RegRuntimeKey(truntime.NewRuntimeKey("user_name", false, func() string { return "" }, false))
	t.InitHTTP()
	db.DB = t.DB()
	migrate.Migrate()
	t.ServeHTTP()
}

// Claims data to sign
type Claims struct {
	ClientID string `json:"client_id,omitempty"`
	UID      string `json:"uid,omitempty"`
	Origin   string `json:"origin,omitempty"`
	jwt.StandardClaims
}

type Response struct {
	Status  int         `json:"status"`
	Data    interface{} `json:"data"`
	Runtime interface{} `json:"runtime"`
}

func CustomizeResponseFactory(status int, res interface{}, runtime map[string]string) interface{} {
	resMap, ok := res.(map[string]interface{})
	if !ok {
		resMap = make(map[string]interface{})
		resMap["data"] = res
	}
	resMap["status"] = status
	for k, v := range runtime {
		resMap[k] = v
	}
	return resMap
}
