package conf

var (
	// DefaultLogLevel default log level
	DefaultLogLevel = "disable"
	// DefaultProjectName default project name
	DefaultProjectName = "TrinityGO"
	// DefaultProjectVersion default project version
	DefaultProjectVersion = "v0.0.1"
	// DefaultDebugMode default debug mode
	DefaultDebugMode = false
	// RandomPort if random port
	RandomPort = false
)

// Conf Conf interface
type Conf struct {
	ProjectName    string
	ProjectVersion string
	Debug          bool
	Loglevel       string
	RandomPort     bool
}

// DefaultConf default conf
func DefaultConf() *Conf {
	return &Conf{
		ProjectName:    DefaultProjectName,
		ProjectVersion: DefaultProjectVersion,
		Debug:          DefaultDebugMode,
		Loglevel:       DefaultLogLevel,
		RandomPort:     RandomPort,
	}
}
