package util

import (
	"fmt"

	"github.com/spf13/viper"
)

func ReadConfig[T any](name string) (T, error) {
	var result T
	viper.AddConfigPath("conf")
	viper.SetConfigType("json")
	viper.SetConfigName(name)
	if err := viper.ReadInConfig(); err != nil {
		return result, fmt.Errorf("Error reading config file, %s", err.Error())
	}
	values := viper.AllSettings()
	fmt.Println(values)
	if err := viper.Unmarshal(&result); err != nil {
		return result, fmt.Errorf("Unable to decode into struct, %v", err.Error())
	}
	return result, nil
}
func ReadConfigAsString(name string) (string, error) {
	viper.AddConfigPath(".")
	viper.SetConfigName(name)
	if err := viper.ReadInConfig(); err != nil {
		return "", fmt.Errorf("Error reading config file, %s", err.Error())
	}
	content := viper.GetString("client_id")
	return content, nil
}

func WriteConfig(values map[string]interface{}, name string) error {
	for k, v := range values {
		viper.Set(k, v)
	}
	viper.SetConfigName(name)
	viper.AddConfigPath("conf")

	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("Unable to cache oauth token: %v", err)
	}
	return nil
}
