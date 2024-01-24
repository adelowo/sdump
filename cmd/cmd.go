package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/adelowo/sdump/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Version describes the version of the current build.
	Version = "dev"

	// Commit describes the commit of the current build.
	Commit = "none"

	// Date describes the date of the current build.
	Date = time.Now().UTC()
)

const (
	defaultConfigFilePath = "config"
	envPrefix             = "SDUMP"
)

func main() {
	if err := Execute(); err != nil {
		log.Fatal(err)
	}
}

func initializeConfig(cfg *config.Config) error {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	viper.AddConfigPath(filepath.Join(homePath, ".config", defaultConfigFilePath))
	viper.AddConfigPath(".")

	viper.SetConfigName(defaultConfigFilePath)
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	viper.SetEnvPrefix(envPrefix)

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetDefault("tui.color_scheme", "monokai")
	viper.SetDefault("log_level", "debug")
	viper.SetDefault("http.database.log_queries", false)
	viper.SetDefault("http.port", 4200)
	viper.SetDefault("http.domain", "sdump.app")

	return viper.Unmarshal(cfg)
}

func Execute() error {
	cfg := &config.Config{}

	rootCmd := &cobra.Command{
		Use:   "sdump",
		Short: "sdump runs a SSH server that helps you view and inspect incoming http requests",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initializeConfig(cfg)
		},
	}

	rootCmd.SetVersionTemplate(
		fmt.Sprintf("Version: %v\nCommit: %v\nDate: %v\n", Version, Commit, Date))

	rootCmd.Flags().StringP("config", "c", defaultConfigFilePath, "Config file. This is in YAML")

	createHTTPCommand(rootCmd, cfg)
	createSSHCommand(rootCmd, cfg)

	return rootCmd.Execute()
}
