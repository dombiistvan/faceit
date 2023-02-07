package helpers

import (
	"faceit/common"
	"faceit/db"
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const configFlag = "config"

type AppConfig struct {
	AppPort uint      `yaml:"app_port" json:"app_port"`
	DB      db.Config `yaml:"db" json:"db"`
}

var appConfig AppConfig

// validate AppConfig struct
func (ac AppConfig) Validate() error {
	return validation.ValidateStruct(&ac,
		validation.Field(&ac.AppPort, validation.Required, validation.Min(uint(1))),
	)
}

// set config file and tries to parse
func init() {
	var configFile string
	flag.StringVar(&configFile, configFlag, ``, "you can add the yaml config file location with this parameter")
	flag.Parse()

	if !common.IsEmptyString(configFile) {
		content, err := common.GetFileContent(configFile)
		if err != nil {
			panic(fmt.Errorf(`could not retrieve file's content: %s`, configFile))
		}

		if err = ParseYaml(content, &appConfig); err != nil {
			panic(fmt.Errorf(`could not parse config %w`, err))
		}
		return
	}

	//panic(fmt.Errorf(`flag does not present: %s`, configFlag))

	appPort, err := strconv.ParseInt(os.Getenv("APP_PORT"), 10, 64)
	if err != nil {
		panic(fmt.Errorf(`you can specify application port in a config .yaml file by running go executable with --config=config.yaml command`, err))
	}
	appConfig.AppPort = uint(appPort)

	if common.IsEmptyString(appConfig.DB.DBName) {
		appConfig.DB.Charset = os.Getenv("CHARSET")
		appConfig.DB.DBName = os.Getenv("MYSQL_DATABASE")
		appConfig.DB.Loc = os.Getenv("LOC")
		appConfig.DB.DBHost = os.Getenv("MYSQL_HOST")
		appConfig.DB.DBUser = os.Getenv("MYSQL_USER")
		appConfig.DB.DBPassword = os.Getenv("MYSQL_PASSWORD")
		appConfig.DB.ParseTime, _ = strconv.ParseBool(os.Getenv("PARSE_TIME"))
	}
}

// GetConfig returns app config
func GetConfig() (*AppConfig, error) {
	if err := appConfig.Validate(); err != nil {
		return nil, err
	}
	return &appConfig, nil
}

// ParseYaml unmarshalls []byte data of yaml content to struct
func ParseYaml(content []byte, into interface{}) error {
	return yaml.Unmarshal(content, into)
}
