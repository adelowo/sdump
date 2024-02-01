package config

// ENUM(psql)
type DatabaseType string

type SSHConfig struct {
	// Port defines where the ssh server runs at
	Port int `mapstructure:"port" json:"port,omitempty" yaml:"port"`

	Host string `mapstructure:"host" yaml:"host" json:"host,omitempty"`
	// Identities is an array containing private keys for the ssh server
	// By default it uses .ssh/id_rsa only
	Identities []string `mapstructure:"identities" json:"identities,omitempty" yaml:"identities"`
}

type HTTPConfig struct {
	// Port to run http server on
	// The server
	Port int `mapstructure:"port" json:"port,omitempty" yaml:"port"`

	// AdminSecret is used to protect routes that are meant to be internal or
	// only ran by an admin
	// Endpoints to create a new url as an example should only be ran by an admin
	// or the ssh server ( after it has verified we have a verified connection)
	// If empty, server would crash
	AdminSecret string `mapstructure:"admin_secret" json:"admin_secret,omitempty" yaml:"admin_secret"`

	Database struct {
		DSN        string `mapstructure:"dsn" json:"dsn,omitempty" yaml:"dsn"`
		LogQueries bool   `mapstructure:"log_queries" json:"log_queries,omitempty" yaml:"log_queries"`
	} `mapstructure:"database" json:"database,omitempty" yaml:"database"`

	Domain             string `json:"domain,omitempty" yaml:"domain" mapstructure:"domain"`
	MaxRequestBodySize int64  `json:"max_request_body_size,omitempty" yaml:"max_request_body_size" mapstructure:"max_request_body_size"`

	OTEL struct {
		UseTLS      bool   `json:"use_tls,omitempty" mapstructure:"use_tls" yaml:"use_tls"`
		ServiceName string `json:"service_name,omitempty" mapstructure:"service_name" yaml:"service_name"`
		Endpoint    string `json:"endpoint,omitempty" mapstructure:"endpoint" yaml:"endpoint"`
		IsEnabled   bool   `json:"is_enabled,omitempty" mapstructure:"is_enabled" yaml:"is_enabled"`
	} `json:"otel,omitempty" mapstructure:"otel" yaml:"otel"`

	// Prometheus config to protect the /metrics endpoint
	// This will be used as basic auth information
	Prometheus struct {
		Username  string `json:"username,omitempty" mapstructure:"username" yaml:"username"`
		Password  string `json:"password,omitempty" mapstructure:"password" yaml:"password"`
		IsEnabled bool   `json:"is_enabled,omitempty" mapstructure:"is_enabled" yaml:"is_enabled"`
	} `json:"prometheus,omitempty" mapstructure:"prometheus" yaml:"prometheus"`
}

type TUIConfig struct {
	ColorScheme string `mapstructure:"color_scheme" yaml:"color_scheme" json:"color_scheme,omitempty"`
}

type Config struct {
	SSH      SSHConfig  `mapstructure:"ssh" json:"ssh,omitempty" yaml:"ssh"`
	HTTP     HTTPConfig `json:"http,omitempty" mapstructure:"http" yaml:"http"`
	LogLevel string     `mapstructure:"log_level" json:"log_level,omitempty" yaml:"log_level"`
	TUI      TUIConfig  `mapstructure:"tui" json:"tui,omitempty" yaml:"tui"`
}
