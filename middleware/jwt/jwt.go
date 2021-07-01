package jwt

import (
	"errors"
	"fmt"
	"strings"

	"github.com/PolarPanda611/trinitygo/application"
	"github.com/PolarPanda611/trinitygo/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
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
	MethodWhiteList    []string
	PathWhiteList      []string
	SkipValidation     bool
}

// DefaultConfig default jwt config
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
		MethodWhiteList:    []string{"OPTION"},
		PathWhiteList:      []string{},
	}

}

// Claims data to sign
type Claims struct {
	Userkey string `json:"userkey"`
	jwt.StandardClaims
}

// JWTVerifier jwt verifier
type JWTVerifier interface {
	parseUnverifiedToken(token string) (jwt.Claims, error)
	Middleware() gin.HandlerFunc
}

// JWTVerifierImpl impl
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

func (m *JWTVerifierImpl) parseToken(token string) (jwt.Claims, error) {
	p := new(jwt.Parser)
	p.SkipClaimsValidation = true
	tkn, err := p.ParseWithClaims(token, m.config.Claims, func(*jwt.Token) (interface{}, error) {
		return []byte(m.config.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !tkn.Valid {
		return nil, errors.New("unauthorized")
	}
	return m.config.Claims, nil
}

// checkUnverifiedTokenValid check authorization header token is valid
func (m *JWTVerifierImpl) checkTokenValid(c *gin.Context) (jwt.Claims, error) {
	if c.Request.Header.Get(m.config.HeaderName) == "" || len(strings.Fields(c.Request.Header.Get(m.config.HeaderName))) != 2 {
		return nil, ErrTokenWrongAuthorization
	}
	prefix := strings.Fields(c.Request.Header.Get(m.config.HeaderName))[0]
	token := strings.Fields(c.Request.Header.Get(m.config.HeaderName))[1]
	if prefix != m.config.JwtPrefix {
		return nil, ErrTokenWrongHeaderPrefix
	}
	if m.config.SkipValidation {
		tokenClaims, err := m.parseUnverifiedToken(token)
		if err != nil {
			return nil, err
		}
		return tokenClaims, nil
	} else {
		tokenClaims, err := m.parseToken(token)
		if err != nil {
			return nil, err
		}
		return tokenClaims, nil
	}

}

// Middleware jwt middleware
func (m *JWTVerifierImpl) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if util.StringInSlice(c.Request.Method, m.config.MethodWhiteList) {
			c.Next()
			return
		}
		if util.StringContainsInSlice(c.Request.URL.Path, m.config.PathWhiteList) {
			c.Next()
			return
		}

		tokenClaims, err := m.checkTokenValid(c)
		if err != nil {

			if m.app.ResponseFactory() != nil {
				c.JSON(401, m.app.ResponseFactory()(401, map[string]string{
					"code":    codes.Internal.String(),
					"message": fmt.Sprintf("Unauthenticated header"),
				},
					nil,
				))
				c.Abort()
			} else {
				c.AbortWithStatusJSON(401, map[string]string{
					"code":    codes.Internal.String(),
					"message": fmt.Sprintf("Unauthenticated header"),
				})
			}
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

// GenerateToken generate tokens used for auth
func GenerateToken(app application.Application, claims jwt.Claims) (string, error) {
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString([]byte(app.Conf().GetJwtSecretKey()))
	if err != nil {
		return "", err
	}
	return token, nil
}
