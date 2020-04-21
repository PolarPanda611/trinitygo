package conf

import (
	"log"

	"github.com/BurntSushi/toml"
)

// Conf Conf interface
type Conf interface {
	GetDebug() bool
	GetTLSEnable() bool
	GetProjectName() string
	GetProjectTags() []string
	GetAppType() string
	GetAppAddress() string
	GetAppPort() int
	GetServiceDiscoveryAddress() string
	GetServiceDiscoveryPort() int
	GetServiceDiscoveryTimeout() int
	GetDeregisterAfterCritical() int
	GetHealthCheckInterval() int
	GetCAPemFile() string
	GetServerPemFile() string
	GetServerKeyFile() string
	GetClientPemFile() string
	GetClientKeyFile() string
	GetProjectVersion() string
	GetLogRootPath() string
	GetLogName() string
	GetServiceDiscoveryType() string
	GetServiceDiscoveryAutoRegister() bool
	GetAtomicRequest() bool
	GetTablePrefix() string
	GetAppMediaURL() string
	GetAppStaticURL() string
	GetAppMediaPath() string
	GetAppStaticPath() string
	GetCacheSize() int
	GetCacheTimeout() int
	GetPageSize() int
	GetAppBaseURL() string
	GetMigrationPath() string
	GetJwtSecretKey() string
	GetJwtVerifyExpireHour() bool
	GetJwtVerifyIssuer() bool
	GetJwtIssuer() string
	GetJwtHeaderPrefix() string
	GetJwtExpireHour() int
	GetAppReadTimeout() int
	GetAppReadHeaderTimeout() int
	GetAppWriteTimeout() int
	GetAppIdleTimeout() int
	GetAppMaxHeaderBytes() int
	GetAllowOrigins() []string
	GetAllowMethods() []string
	GetAllowHeaders() []string
	GetExposeHeaders() []string
	GetAllowCredentials() bool
	GetMaxAgeHour() int
	GetCorsEnable() bool
	GetDBType() string
	GetDBServer() string
	GetDBTablePrefix() string
	GetDbMaxIdleConn() int
	GetDbMaxOpenConn() int
}

// Project struct
type Project struct {
	Name    string
	Version string
	Tags    []string
}

// Database struct
type Database struct {
	DBType      string `toml:"db_type"`
	Server      string `toml:"server"`
	TablePrefix string `toml:"table_prefix"`
	MaxIdleConn int    `toml:"max_idle_conn"`
	MaxOpenConn int    `toml:"max_open_conn"`
}

// Runtime struct
type Runtime struct {
	Debug bool `toml:"debug"`
}

// ServiceDiscovery service delivery
type ServiceDiscovery struct {
	Type                    string // etcd oor consul
	Address                 string
	Port                    int
	Timeout                 int
	DeregisterAfterCritical int  `toml:"deregister_after_critical"` //second
	HealthCheckInterval     int  `toml:"health_check_interval"`     //second
	AutoRegister            bool `toml:"auto_register"`
}

//Jwt jwt
type Jwt struct {
	SecretKey        string `toml:"secret_key"`
	Issuer           string `toml:"issuer"`
	ExpireHour       int    `toml:"expire_hour"`
	HeaderPrefix     string `toml:"header_prefix"`
	VerifyIssuer     bool   `toml:"verify_issuer"`
	VerifyExpireHour bool   `toml:"verify_expire_hour"`
}

// Cors coors
type Cors struct {
	Enable           bool     `toml:"enable"`
	AllowOrigins     []string `toml:"allow_origins"`
	AllowMethods     []string `toml:"allow_methods"`
	AllowHeaders     []string `toml:"allow_headers"`
	ExposeHeaders    []string `toml:"expose_headers"`
	AllowCredentials bool     `toml:"allow_credentials"`
	MaxAgeHour       int      `toml:"max_age_hour"`
}

// TLS tls
type TLS struct {
	Enable        bool   `toml:"enable"`
	CAPemFile     string `toml:"ca_pem_file"`
	ServerPemFile string `toml:"server_pem_file"`
	ServerKeyFile string `toml:"server_key_file"`
	ClientPemFile string `toml:"client_pem_file"`
	ClientKeyFile string `toml:"client_key_file"`
}

// Security security
type Security struct {
	Jwt  Jwt  `toml:"jwt"`
	Cors Cors `toml:"cors"`
	TLS  TLS  `toml:"tls"`
}

// DefaultConf default conf
type DefaultConf struct {
	Project  Project  `toml:"project"`
	Runtime  Runtime  `toml:"runtime"`
	Security Security `toml:"security"`
	App      struct {
		// Type support GRPC HTTP
		Type    string `toml:"type"`
		Address string `toml:"address"`
		Port    int    `toml:"port"`
		// ReadTimeout is the maximum duration for reading the entire
		// request, including the body.
		//
		// Because ReadTimeout does not let Handlers make per-request
		// decisions on each request body's acceptable deadline or
		// upload rate, most users will prefer to use
		// ReadHeaderTimeout. It is valid to use them both.
		ReadTimeout int `toml:"read_timeout"`

		// ReadHeaderTimeout is the amount of time allowed to read
		// request headers. The connection's read deadline is reset
		// after reading the headers and the Handler can decide what
		// is considered too slow for the body. If ReadHeaderTimeout
		// is zero, the value of ReadTimeout is used. If both are
		// zero, there is no timeout.
		ReadHeaderTimeout int `toml:"read_header_timeout"`

		// WriteTimeout is the maximum duration before timing out
		// writes of the response. It is reset whenever a new
		// request's header is read. Like ReadTimeout, it does not
		// let Handlers make decisions on a per-request basis.
		WriteTimeout int `toml:"writer_timeout"`

		// IdleTimeout is the maximum amount of time to wait for the
		// next request when keep-alives are enabled. If IdleTimeout
		// is zero, the value of ReadTimeout is used. If both are
		// zero, there is no timeout.
		IdleTimeout int `toml:"idle_timeout"`

		// MaxHeaderBytes controls the maximum number of bytes the
		// server will read parsing the request header's keys and
		// values, including the request line. It does not limit the
		// size of the request body.
		// If zero, DefaultMaxHeaderBytes is used.
		MaxHeaderBytes int    `toml:"max_header_bytes"`
		TemplatePath   string `toml:"template_path"`
		MediaURL       string `toml:"media_url"`
		MediaPath      string `toml:"media_path"`
		StaticURL      string `toml:"static_url"`
		StaticPath     string `toml:"static_path"`
		MigrationPath  string `toml:"migration_path"`
		PageSize       int    `toml:"page_size"`
		MaxBodySize    int    `toml:"max_body_size"`
		AtomicRequest  bool   `toml:"atomic_request"`
		// if api root is not root , replease with base url
		// e.g : /assetgo
		BaseURL string `toml:"base_url"`
	}
	Log struct {
		LogRootPath string `toml:"log_root_path"` //   /var/log/mold
		LogName     string `toml:"log_name"`      //  app.log
	}
	Cache struct {
		Redis struct {
			Host        string `toml:"host"`
			Port        int    `toml:"port"`
			Password    string `toml:"password"`
			MaxIdle     int    `toml:"max_idle"`
			MaxActive   int    `toml:"max_active"`
			IdleTimeout int    `toml:"idle_timeout"`
		}
		Gcache struct {
			CacheAlgorithm string `toml:"cache_algorithm"`
			CacheSize      int    `toml:"cachesize"`
			Timeout        int    `toml:"timeout"` // hour
		}
	}
	Database         Database         `toml:"database"`
	ServiceDiscovery ServiceDiscovery `toml:"service_discovery"`
}

// GetDbMaxIdleConn get db max idle connection
func (s *DefaultConf) GetDbMaxIdleConn() int { return s.Database.MaxIdleConn }

// GetDbMaxOpenConn get db max open connection
func (s *DefaultConf) GetDbMaxOpenConn() int { return s.Database.MaxOpenConn }

// GetDBServer get db host
func (s *DefaultConf) GetDBServer() string { return s.Database.Server }

// GetCorsEnable get if enable cors
func (s *DefaultConf) GetCorsEnable() bool { return s.Security.Cors.Enable }

// GetMaxAgeHour get max age hour
func (s *DefaultConf) GetMaxAgeHour() int { return s.Security.Cors.MaxAgeHour }

// GetAllowOrigins get allow origins
func (s *DefaultConf) GetAllowOrigins() []string { return s.Security.Cors.AllowOrigins }

// GetAllowMethods get allow method
func (s *DefaultConf) GetAllowMethods() []string { return s.Security.Cors.AllowMethods }

// GetAllowHeaders get allow headers
func (s *DefaultConf) GetAllowHeaders() []string { return s.Security.Cors.AllowHeaders }

// GetExposeHeaders get expoose headers
func (s *DefaultConf) GetExposeHeaders() []string { return s.Security.Cors.ExposeHeaders }

// GetAllowCredentials get allow credentials
func (s *DefaultConf) GetAllowCredentials() bool {
	return s.Security.Cors.AllowCredentials
}

// GetAppReadTimeout get readtimeoout
func (s *DefaultConf) GetAppReadTimeout() int { return s.App.ReadTimeout }

// GetAppReadHeaderTimeout get GetReadHeaderTimeoutSecond
func (s *DefaultConf) GetAppReadHeaderTimeout() int { return s.App.ReadHeaderTimeout }

// GetAppWriteTimeout get GetWriteTimeoutSecond
func (s *DefaultConf) GetAppWriteTimeout() int { return s.App.WriteTimeout }

// GetAppIdleTimeout get GetIdleTimeoutSecond
func (s *DefaultConf) GetAppIdleTimeout() int { return s.App.IdleTimeout }

// GetAppMaxHeaderBytes get GetMaxHeaderBytes
func (s *DefaultConf) GetAppMaxHeaderBytes() int { return s.App.MaxHeaderBytes }

// GetJwtSecretKey get GetSecretKey
func (s *DefaultConf) GetJwtSecretKey() string {
	return s.Security.Jwt.SecretKey
}

// GetJwtExpireHour get GetJwtExpireHour
func (s *DefaultConf) GetJwtExpireHour() int {
	return s.Security.Jwt.ExpireHour
}

// GetJwtHeaderPrefix get GetJwtHeaderPrefix
func (s *DefaultConf) GetJwtHeaderPrefix() string {
	return s.Security.Jwt.HeaderPrefix
}

// GetJwtIssuer get GetJwtIssuer
func (s *DefaultConf) GetJwtIssuer() string {
	return s.Security.Jwt.Issuer
}

// GetJwtVerifyIssuer get GetJwtVerifyIssuer
func (s *DefaultConf) GetJwtVerifyIssuer() bool {
	return s.Security.Jwt.VerifyIssuer
}

// GetJwtVerifyExpireHour get GetJwtVerifyExpireHour
func (s *DefaultConf) GetJwtVerifyExpireHour() bool {
	return s.Security.Jwt.VerifyExpireHour
}

// GetMigrationPath get GetMigrationPath
func (s *DefaultConf) GetMigrationPath() string { return s.App.MigrationPath }

// GetAppBaseURL get GetAppBaseURL
func (s *DefaultConf) GetAppBaseURL() string { return s.App.BaseURL }

// GetPageSize get GetPageSize
func (s *DefaultConf) GetPageSize() int { return s.App.PageSize }

// GetCacheSize get GetCacheSize
func (s *DefaultConf) GetCacheSize() int { return s.Cache.Gcache.CacheSize }

// GetCacheTimeout get GetCacheTimeout
func (s *DefaultConf) GetCacheTimeout() int { return s.Cache.Gcache.Timeout }

// GetAppMediaURL get web app media url
func (s *DefaultConf) GetAppMediaURL() string { return s.App.MediaURL }

// GetAppMediaPath get web app media path
func (s *DefaultConf) GetAppMediaPath() string { return s.App.MediaPath }

// GetAppStaticPath get web app static path
func (s *DefaultConf) GetAppStaticPath() string { return s.App.StaticPath }

// GetAppStaticURL get web app static url
func (s *DefaultConf) GetAppStaticURL() string { return s.App.StaticURL }

// GetLogRootPath get log root path
func (s *DefaultConf) GetLogRootPath() string {
	return s.Log.LogRootPath
}

// GetTablePrefix get table prefix
func (s *DefaultConf) GetTablePrefix() string {
	return s.Database.TablePrefix
}

// GetServiceDiscoveryAutoRegister get auto register
func (s *DefaultConf) GetServiceDiscoveryAutoRegister() bool {
	return s.ServiceDiscovery.AutoRegister
}

// GetAtomicRequest get automic request is open
func (s *DefaultConf) GetAtomicRequest() bool {
	return s.App.AtomicRequest
}

//GetServiceDiscoveryType get s m type
func (s *DefaultConf) GetServiceDiscoveryType() string {
	return s.ServiceDiscovery.Type
}

// GetTLSEnable get tls enabled
func (s *DefaultConf) GetTLSEnable() bool {
	return s.Security.TLS.Enable
}

// GetLogName get log name
func (s *DefaultConf) GetLogName() string {
	return s.Log.LogName
}

// GetDebug get debug
func (s *DefaultConf) GetDebug() bool {
	return s.Runtime.Debug
}

// GetDefaultConf get DefaultConf
func (s *DefaultConf) GetDefaultConf() *DefaultConf {
	return s
}

//GetCAPemFile get ca pem file
func (s *DefaultConf) GetCAPemFile() string {
	return s.Security.TLS.CAPemFile
}

//GetServerPemFile get server pem file
func (s *DefaultConf) GetServerPemFile() string {
	return s.Security.TLS.ServerPemFile
}

//GetServerKeyFile get server key file
func (s *DefaultConf) GetServerKeyFile() string {
	return s.Security.TLS.ServerKeyFile
}

//GetClientPemFile get client pem file
func (s *DefaultConf) GetClientPemFile() string {
	return s.Security.TLS.ClientPemFile
}

//GetClientKeyFile get client key file
func (s *DefaultConf) GetClientKeyFile() string {
	return s.Security.TLS.ClientKeyFile
}

// GetDeregisterAfterCritical deregister service after critical second
func (s *DefaultConf) GetDeregisterAfterCritical() int {
	return s.ServiceDiscovery.DeregisterAfterCritical
}

// GetHealthCheckInterval health check interval
func (s *DefaultConf) GetHealthCheckInterval() int {
	return s.ServiceDiscovery.HealthCheckInterval

}

//GetProjectTags get project tags
func (s *DefaultConf) GetProjectTags() []string {
	return s.Project.Tags
}

// GetProjectName get project name
func (s *DefaultConf) GetProjectName() string {
	return s.Project.Name
}

// GetProjectVersion get project name
func (s *DefaultConf) GetProjectVersion() string {
	return s.Project.Version
}

// GetAppType get web app type
func (s *DefaultConf) GetAppType() string {
	return s.App.Type
}

// GetAppAddress get web service  ip address
func (s *DefaultConf) GetAppAddress() string {
	return s.App.Address
}

// GetAppPort get web service port
func (s *DefaultConf) GetAppPort() int {
	return s.App.Port
}

// GetServiceDiscoveryAddress get service mesh address
func (s *DefaultConf) GetServiceDiscoveryAddress() string {
	return s.ServiceDiscovery.Address
}

// GetServiceDiscoveryPort get service mesh port
func (s *DefaultConf) GetServiceDiscoveryPort() int {
	return s.ServiceDiscovery.Port
}

// GetServiceDiscoveryTimeout get service mesh port
func (s *DefaultConf) GetServiceDiscoveryTimeout() int {
	return s.ServiceDiscovery.Timeout

}

// GetDBType get db type
func (s *DefaultConf) GetDBType() string {
	return s.Database.DBType
}

// GetDBTablePrefix get db table prefix
func (s *DefaultConf) GetDBTablePrefix() string {
	return s.Database.TablePrefix
}

// NewSetting new setting
func NewSetting(configFilePath string) Conf {
	return LoadConfigFile(configFilePath)
}

// LoadConfigFile load config file
func LoadConfigFile(configFilePath string) Conf {
	var d DefaultConf
	_, err := toml.DecodeFile(configFilePath, &d)
	if err != nil {
		log.Fatal("load log error ", err)
	}
	return &d
}
