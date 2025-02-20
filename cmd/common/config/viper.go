// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Khulnasoft

package config

import (
	"strings"

	"github.com/khulnasoft/triangle/pkg/defaults"
	"github.com/spf13/viper"
)

// NewViper creates a new viper instance configured for Triangle.
func NewViper() *viper.Viper {
	vp := viper.New()

	// read config from a file
	vp.SetConfigName("config") // name of config file (without extension)
	vp.SetConfigType("yaml")   // useful if the given config file does not have the extension in the name
	vp.AddConfigPath(".")      // look for a config in the working directory first
	if defaults.ConfigDir != "" {
		vp.AddConfigPath(defaults.ConfigDir)
	}
	if defaults.ConfigDirFallback != "" {
		vp.AddConfigPath(defaults.ConfigDirFallback)
	}

	// read config from environment variables
	vp.SetEnvPrefix("triangle") // env var must start with TRIANGLE_
	// replace - by _ for environment variable names
	// (eg: the env var for tls-server-name is TLS_SERVER_NAME)
	vp.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	vp.AutomaticEnv() // read in environment variables that match
	return vp
}
