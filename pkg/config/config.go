package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
)

var (
	_, b, _, _ = runtime.Caller(0)
	// Root folder of this project
	Root = filepath.Join(filepath.Dir(b), "../../..")
)
var Config config

type config struct {
	Env struct {
		HttpClientTimeOut int64  `yaml:"httpClientTimeOut" env:"DEC_TIME_OUT"`
		FileAddr          string `yaml:"fileAddr" env:"FILE_ADDR"`
		DecApiAddr        string `yaml:"decApiAddr" env:"DEC_API_ADDR"`
		AK                string `yaml:"appAk" env:"APP_AK"`
		SK                string `yaml:"appSk" env:"APP_SK"`
		PORT              string `yaml:"port" env:"PORT"`
		DriveId           string `yaml:"driveId" env:"DRIVE_ID"`
		DBAddress         string `yaml:"dbAddress" env:"REDIS_ADDR"`
		DBUserName        string `yaml:"dbUserName" env:"REDIS_USER"`
		DBPassWord        string `yaml:"dbPassWord" env:"REDIS_PW"`
	} `yaml:"env"`
	GoVERSION string `yaml:"dbPassWord" env:"GOVERSION"`
}

func UnmarshalConfig(config interface{}, configName string) {
	var env string
	if configName == "config.yaml" {
		env = "CONFIG_NAME"
	} else {
		panic("configName must be config.yaml")
	}
	cfgName := os.Getenv(env)
	if len(cfgName) != 0 {
		bytes, err := os.ReadFile(filepath.Join(cfgName, "config", configName))
		if err != nil {
			bytes, err = os.ReadFile(filepath.Join(Root, "config", configName))
			if err != nil {
				panic(err.Error() + " config: " + filepath.Join(cfgName, "config", configName))
			}
		} else {
			Root = cfgName
		}
		if err = yaml.Unmarshal(bytes, config); err != nil {
			panic(err.Error())
		}
	} else {
		bytes, err := os.ReadFile(fmt.Sprintf("%s", configName))
		if err != nil {
			bytes, err = os.ReadFile(fmt.Sprintf("./config/%s", configName))
			if err != nil {
				panic(err.Error() + configName)
			}
		}
		if err = yaml.Unmarshal(bytes, config); err != nil {
			panic(err.Error())
		}
	}
	loadEnv(config)
}
func init() {
	UnmarshalConfig(&Config, "config.yaml")
}

//	func loadEnv(config interface{}) {
//		val := reflect.ValueOf(config).Elem()
//		typ := val.Type()
//
//		for i := 0; i < val.NumField(); i++ {
//			field := val.Field(i)
//			fieldType := typ.Field(i)
//			envTag := fieldType.Tag.Get("env")
//			if envTag == "" {
//				continue
//			}
//			envValue := os.Getenv(envTag)
//			if envValue == "" {
//				continue
//			}
//			switch field.Kind() {
//			case reflect.String:
//				field.SetString(envValue)
//			case reflect.Int:
//				intValue, err := strconv.Atoi(envValue)
//				if err != nil {
//					fmt.Println(fmt.Errorf("invalid value for %s: %v", envTag, err))
//					continue
//				}
//				field.SetInt(int64(intValue))
//			case reflect.Bool:
//				boolValue, err := strconv.ParseBool(envValue)
//				if err != nil {
//					fmt.Println(fmt.Errorf("invalid value for %s: %v", envTag, err))
//					continue
//				}
//				field.SetBool(boolValue)
//			default:
//				fmt.Println(fmt.Errorf("unsupported type for %s", envTag))
//				continue
//			}
//		}
//	}
func loadEnv(config interface{}) {
	val := reflect.ValueOf(config).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		envTag := fieldType.Tag.Get("env")
		if envTag == "" {
			continue
		}

		envValue := os.Getenv(envTag)
		if envValue == "" {
			continue
		}

		if field.Kind() == reflect.Struct {
			loadEnv(field.Addr().Interface()) // 递归处理嵌套结构体
			continue
		}

		switch field.Kind() {
		case reflect.String:
			field.SetString(envValue)
		case reflect.Int:
			intValue, err := strconv.Atoi(envValue)
			if err != nil {
				fmt.Println(fmt.Errorf("invalid value for %s: %v", envTag, err))
				continue
			}
			field.SetInt(int64(intValue))
		case reflect.Bool:
			boolValue, err := strconv.ParseBool(envValue)
			if err != nil {
				fmt.Println(fmt.Errorf("invalid value for %s: %v", envTag, err))
				continue
			}
			field.SetBool(boolValue)
		default:
			fmt.Println(fmt.Errorf("unsupported type for %s", envTag))
			continue
		}
	}
}
