package jwt

import (
	"errors"
	"fmt"
	"strings"

	"github.com/PolarPanda611/trinitygo/application"
	"github.com/PolarPanda611/trinitygo/httputils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var (
	// ErrTokenExpired return token expired err
	ErrTokenExpired = errors.New("app.err.TokenExpired")
	// ErrTokenWrongIssuer return wrong issuer err
	ErrTokenWrongIssuer = errors.New("app.err.TokenWrongIssuer")
	// ErrTokenWrongHeaderPrefix return wrong header prefix
	ErrTokenWrongHeaderPrefix = errors.New("app.err.TokenWrongHeaderPrefix")
	// ErrTokenWrongAuthorization return wrong authorization
	ErrTokenWrongAuthorization = errors.New("app.err.TokenWrongAuthorization")
)

// Config jwt config
type Config struct {
	HeaderName         string
	SecretKey          string
	JwtPrefix          string
	Issuer             string
	ExpireHour         int
	IsVerifyIssuer     bool
	IsVerifyExpireHour bool
	IsVerifySecretKey  bool
	Claims             jwt.Claims
	SuccessCallback    func(c *gin.Context, claim jwt.Claims)
}

func DefaultConfig() *Config {
	return &Config{
		HeaderName:         "Authorization",
		SecretKey:          "test",
		JwtPrefix:          "Bearer",
		Issuer:             "xxx",
		ExpireHour:         2,
		IsVerifyIssuer:     false,
		IsVerifyExpireHour: false,
		IsVerifySecretKey:  false,
	}

}

// Claims data to sign
type Claims struct {
	Userkey string `json:"userkey"`
	jwt.StandardClaims
}

type JWTVerifier interface {
	parseUnverifiedToken(token string) (jwt.Claims, error)
	Middleware() gin.HandlerFunc
}

type JWTVerifierImpl struct {
	app    application.Application
	config *Config
}

func (m *JWTVerifierImpl) parseUnverifiedToken(token string) (jwt.Claims, error) {
	p := new(jwt.Parser)
	p.SkipClaimsValidation = true
	_, _, err := p.ParseUnverified(token, m.config.Claims)
	if err != nil {
		return nil, err
	}

	return m.config.Claims, nil
}

// checkUnverifiedTokenValid check authorization header token is valid
func (m *JWTVerifierImpl) checkUnverifiedTokenValid(c *gin.Context) (jwt.Claims, error) {
	if c.Request.Header.Get(m.config.HeaderName) == "" || len(strings.Fields(c.Request.Header.Get(m.config.HeaderName))) != 2 {
		return nil, ErrTokenWrongAuthorization
	}
	prefix := strings.Fields(c.Request.Header.Get(m.config.HeaderName))[0]
	token := strings.Fields(c.Request.Header.Get(m.config.HeaderName))[1]
	if prefix != m.config.JwtPrefix {
		return nil, ErrTokenWrongHeaderPrefix
	}
	tokenClaims, err := m.parseUnverifiedToken(token)
	if err != nil {
		return nil, err
	}
	return tokenClaims, nil
}

func (m *JWTVerifierImpl) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenClaims, err := m.checkUnverifiedTokenValid(c)
		if err != nil {
			c.AbortWithStatusJSON(401, httputils.ResponseData{
				Status: 401,
				Result: fmt.Sprintf("runtime key %v is required ", err),
			})
			return
		}

		m.config.SuccessCallback(c, tokenClaims)
		c.Next()
	}
}

// New jwt middleware
func New(app application.Application, config ...*Config) gin.HandlerFunc {
	c := DefaultConfig()
	if len(config) > 0 {
		c = config[0]
	}
	jwt := JWTVerifierImpl{
		app:    app,
		config: c,
	}
	return jwt.Middleware()
}
