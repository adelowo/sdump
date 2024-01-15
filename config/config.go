package config

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

type Config struct {
	SSH SSHConfig `mapstructure:"ssh" json:"ssh,omitempty" yaml:"ssh"`
}
