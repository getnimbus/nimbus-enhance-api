package conf

import (
	"reflect"
	"strings"

	"github.com/creasty/defaults"
	"github.com/spf13/viper"
	"github.com/tikivn/ultrago/u_logger"
)

var Config config

func init() {
	if err := defaults.Set(&Config); err != nil {
		panic(err)
	}
	bindEnvs(Config)
}

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type config struct {
	Env       string `mapstructure:"ENV" default:"dev"`
	Debug     string `mapstructure:"DEBUG" default:"no"`
	NumBatch  int    `mapstructure:"NUM_BATCH" default:"5"`
	NumWorker int    `mapstructure:"NUM_WORKER" default:"3"`

	// db
	Migration string `mapstructure:"MIGRATION" default:"no"`
	GormDsn   string `mapstructure:"GORM_DSN" default:"-"`

	// redis
	RedisAddress string `mapstructure:"REDIS_ADDRESS" default:"localhost:6379"`
	RedisDB      int    `mapstructure:"REDIS_DB" default:"0"`
	RedisUser    string `mapstructure:"REDIS_USER" default:"-"`
	RedisPass    string `mapstructure:"REDIS_PASS" default:"-"`
}

func (c *config) IsLocal() bool {
	return !c.IsDev() && !c.IsProd()
}

func (c *config) IsDev() bool {
	return strings.ToLower(c.Env) == "dev"
}

func (c *config) IsProd() bool {
	return strings.ToLower(c.Env) == "prod"
}

func (c *config) IsDebug() bool {
	return strings.ToLower(c.Debug) == "yes"
}

func (c *config) IsMigration() bool {
	return strings.ToLower(c.Migration) == "yes"
}

// LoadConfig read configuration for both file and system environment
func LoadConfig(path string) error {
	logger := u_logger.NewLogger()

	// Read from .env file
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")

	// priority load os ENV before .env file
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logger.Warnf("config file not found: %v", err)
		} else {
			logger.Warnf("config file is invalid: %v", err)
		}
	}

	if err := viper.Unmarshal(&Config); err != nil {
		return err
	}

	return nil
}

func bindEnvs(iface interface{}, parts ...string) {
	ifv := reflect.ValueOf(iface)
	ift := reflect.TypeOf(iface)
	for i := 0; i < ift.NumField(); i++ {
		v := ifv.Field(i)
		t := ift.Field(i)
		tv, ok := t.Tag.Lookup("mapstructure")
		if !ok {
			continue
		}
		switch v.Kind() {
		case reflect.Struct:
			bindEnvs(v.Interface(), append(parts, tv)...)
		default:
			viper.BindEnv(strings.Join(append(parts, tv), "."))
		}
	}
}
