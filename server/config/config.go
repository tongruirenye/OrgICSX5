package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	WebDevRoot  string
	WebDevUser  string
	WebDevPass  string
	Sign        string
	Port        string
	Mode        string
	CalName     string
	Project     string
	viperConfig *viper.Viper
}

var AppConfig *Config
var AppPath string

func InitConfig(path string) error {
	if AppConfig != nil {
		return nil
	}
	AppConfig = &Config{
		viperConfig: viper.New(),
	}
	if path != "" {
		AppConfig.viperConfig.SetConfigFile(path)
	} else {
		AppConfig.viperConfig.AddConfigPath("./conf")
		AppConfig.viperConfig.SetConfigName("conf.toml")
	}

	AppConfig.viperConfig.SetConfigType("toml")
	if err := AppConfig.viperConfig.ReadInConfig(); err != nil {
		return err
	}

	AppConfig.WebDevRoot = AppConfig.viperConfig.GetString("webdev.root")
	AppConfig.WebDevUser = AppConfig.viperConfig.GetString("webdev.user")
	AppConfig.WebDevPass = AppConfig.viperConfig.GetString("webdev.pass")
	AppConfig.Sign = AppConfig.viperConfig.GetString("auth.sign")
	AppConfig.Port = AppConfig.viperConfig.GetString("base.port")
	AppConfig.Mode = AppConfig.viperConfig.GetString("base.mode")
	AppConfig.CalName = AppConfig.viperConfig.GetString("base.cal")
	AppConfig.Project = AppConfig.viperConfig.GetString("base.project")

	return nil
}

func init() {
	var err error
	if AppPath, err = filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
		panic(err)
	}
	AppPath = AppPath + "/"
}
