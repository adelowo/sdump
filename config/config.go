package config

// ENUM(sqlite3,psql,mysql)
type DatabaseType string

type SSHConfig struct {
	// Port defines where the ssh server runs at
	Port int `mapstructure:"port" json:"port,omitempty" yaml:"port"`
	// Allowlist is an array pointing to a bunch of public keys that
	// are allowed to connect to the ssh server
	Allowlist []string `mapstructure:"allowlist" json:"allowlist,omitempty" yaml:"allowlist"`
	// Identities is an array containing private keys for the ssh server
	// By default it uses ~/.ssh/id_rsa only
	Identities []string `mapstructure:"identities" json:"identities,omitempty" yaml:"identities"`
}

type HTTPConfig struct {
	// Port to run http server on
	// The server
	Port int `mapstructure:"port" json:"port,omitempty"`

	// AdminSecret is used to protect routes that are meant to be internal or
	// only ran by an admin
	// Endpoints to create a new url as an example should only be ran by an admin
	// or the ssh server ( after it has verified we have a verified connection)
	// If empty, server would crash
	AdminSecret string `mapstructure:"admin_secret" json:"admin_secret,omitempty"`

	Database struct {
		DSN  string       `mapstructure:"dsn" json:"dsn,omitempty" yaml:"dsn"`
		Type DatabaseType `mapstructure:"type" json:"type,omitempty" yaml:"type"`
	} `mapstructure:"database" json:"database,omitempty" yaml:"database"`

	Domain string `json:"domain,omitempty"`
}

type Config struct {
	SSH  SSHConfig  `mapstructure:"ssh" json:"ssh,omitempty" yaml:"ssh"`
	HTTP HTTPConfig `json:"http,omitempty" mapstructure:"http"`
	Log  string     `mapstructure:"log" json:"log,omitempty"`
}
