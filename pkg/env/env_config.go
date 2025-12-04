package env

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

func LoadEnv(config interface{}) {
	val := reflect.ValueOf(config).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		envJsonTag := fieldType.Tag.Get("env_json")
		if envJsonTag != "" {
			envValue := os.Getenv(envJsonTag)
			if envValue == "" {
				defValue := fieldType.Tag.Get("def")
				if defValue != "" {
					envValue = defValue
				}
			}
			if envValue != "" {
				err := json.Unmarshal([]byte(envValue), field.Addr().Interface())
				if err != nil {
					fmt.Println(fmt.Errorf("invalid value for %s: %v", envJsonTag, err))
				}
			}
		}

		// Check if the field is a struct and need to recursively process it.
		if field.Kind() == reflect.Struct {
			// Recursively call loadEnv for nested structs
			LoadEnv(field.Addr().Interface())
			continue
		}
		envTag := fieldType.Tag.Get("env")
		if envTag == "" {
			continue
		}

		envValue := os.Getenv(envTag)
		if envValue == "" {
			defValue := fieldType.Tag.Get("def")
			if defValue != "" {
				envValue = defValue
			} else {
				continue
			}
		}

		switch field.Kind() {
		case reflect.String:
			field.SetString(envValue)
		case reflect.Slice:
			// Assuming we want to handle it as a slice of strings
			if field.Type().Elem().Kind() == reflect.String {
				field.Index(0).SetString(envValue)
			}
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
		case reflect.Struct:
			err := json.Unmarshal([]byte(envValue), field.Addr().Interface())
			if err != nil {
				fmt.Println(fmt.Errorf("invalid value for %s: %v", envTag, err))
				continue
			}
		default:
			fmt.Println(fmt.Errorf("unsupported type for %s", envTag))
			continue
		}
	}
}
