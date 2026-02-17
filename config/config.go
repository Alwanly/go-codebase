package config

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

func LoadConfig(configName string) (GlobalConfig, error) {
	// Load default config
	loadDefaults()

	// Load config from file
	viper.AddConfigPath(".")
	viper.SetConfigType("env")
	viper.SetConfigName(configName)

	// load from environment variables
	viper.AutomaticEnv()

	// -- binding environment variables as config source
	emptyConfig := &GlobalConfig{}

	// count the number of fields (config) on this struct
	fieldCount := reflect.TypeOf(emptyConfig).Elem().NumField()

	for i := 0; i < fieldCount; i++ {
		// get the field tag name, for example `mapstructure:"ENV"`
		field := string(reflect.TypeOf(emptyConfig).Elem().Field(i).Tag)

		// trim the tag and just return the tag name, for example `ENV`
		// then bind it as env to load to viper config
		if err := viper.BindEnv(field[14 : len(field)-1]); err != nil {
			fmt.Println("Error bind env.", err)
		}
	}

	// read config from sources (config file and environment variables)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Println("Error reading config file:", err)
			return GlobalConfig{}, err
		}
	}

	// unmarshal config
	c := GlobalConfig{}
	if err := viper.Unmarshal(&c); err != nil {
		return GlobalConfig{}, err
	}

	// validate config
	if err := c.Validate(); err != nil {
		return GlobalConfig{}, fmt.Errorf("config validation failed: %w", err)
	}

	return c, nil
}

// Validate validates the configuration
func (c *GlobalConfig) Validate() error {
	var errs []string

	// Validate required fields
	if c.ServiceName == "" {
		errs = append(errs, "SERVICE_NAME is required")
	}

	if c.Port <= 0 {
		errs = append(errs, "PORT must be greater than 0")
	}

	if c.PostgresURI == "" {
		errs = append(errs, "POSTGRES_URI is required")
	}

	// Validate authentication config if JWT is used
	if c.JwtIssuer == "" {
		errs = append(errs, "JWT_ISSUER is required")
	}

	if c.JwtAudience == "" {
		errs = append(errs, "JWT_AUDIENCE is required")
	}

	if c.JwtExpirationTime <= 0 {
		errs = append(errs, "JWT_EXPIRATION must be greater than 0")
	}

	if c.JwtRefreshTime <= 0 {
		errs = append(errs, "JWT_REFRESH_EXPIRATION must be greater than 0")
	}

	// Warn about missing keys in production (but don't fail)
	if c.Environment == "production" {
		if c.PrivateKey == "" {
			errs = append(errs, "PRIVATE_KEY should be set in production")
		}
		if c.PublicKey == "" {
			errs = append(errs, "PUBLIC_KEY should be set in production")
		}
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}

	return nil
}
