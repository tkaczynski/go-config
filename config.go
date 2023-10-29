package go_config

import (
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	settings       map[string]interface{}
	loadedSettings map[string]interface{}
}

func ConfigDefault() (*Config, error) {
	const (
		configLocationEnvVar = "APP_CONFIG"
		configFileName       = ".app-config"
	)

	if envConfigFileName, ok := os.LookupEnv(configLocationEnvVar); ok {
		return ConfigFromFile(envConfigFileName)
	} else {
		homeConfigFileName := os.Getenv("HOME") + "/" + configFileName
		if _, err := os.Stat(configFileName); err == nil {
			return ConfigFromFile(configFileName)
		} else if _, err := os.Stat(homeConfigFileName); err == nil {
			return ConfigFromFile(homeConfigFileName)
		}
	}

	return &Config{
		settings:       map[string]interface{}{},
		loadedSettings: map[string]interface{}{},
	}, nil
}

func ConfigFromFile(filename string) (*Config, error) {
	config := &Config{
		settings:       map[string]interface{}{},
		loadedSettings: map[string]interface{}{},
	}

	var err error
	if _, err := os.Stat(filename); err == nil {
		err = config.loadSettings(filename)
	}

	if err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) GetBool(setting string) (bool, bool) {
	value, ok := c.getSetting(setting)
	if !ok {
		return false, false
	}

	switch value.(type) {
	default:
		return false, false
	case bool:
		return value.(bool), true
	case int:
		return value.(int) != 0, true
	case string:
		v := strings.ToLower(value.(string))
		return v == "true" || v == "yes", true
	}
}

func (c *Config) MustBool(setting string) bool {
	if value, ok := c.GetBool(setting); ok {
		return value
	}
	
	return false
}

func (c *Config) GetInt(setting string) (int, bool) {
	value, ok := c.getSetting(setting)
	if !ok {
		return 0, false
	}
	
	switch value.(type) {
	default:
		return 0, false
	case int:
		return value.(int), true
	case string:
		v, err := strconv.Atoi(value.(string))
		if err != nil {
			return 0, false
		}
		
		return v, true
	}
}

func (c *Config) MustInt(setting string) int {
	if value, ok := c.GetInt(setting); ok {
		return value
	}

	return 0
}

func (c *Config) GetString(setting string) (string, bool) {
	value, ok := c.getSetting(setting)
	if !ok {
		return "", false
	}
	
	switch value.(type) {
	default:
		return "", false
	case int:
		return strconv.Itoa(value.(int)), true
	case bool:
		return strconv.FormatBool(value.(bool)), true
	case string:
		return value.(string), true
	}
}

func (c *Config) MustString(setting string) string {
	value, _ := c.GetString(setting)
	return value
}

func (c *Config) loadSettings(fname string) error {
	file, err := os.ReadFile(fname)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(file, &c.loadedSettings)
}

func (c *Config) getSetting(setting string) (interface{}, bool) {
	if value, ok := c.settings[setting]; ok {
		return value, true
	}

	envVarName := settingToEnvVarName(setting)
	if value, ok := os.LookupEnv(envVarName); ok {
		c.settings[setting] = value
		return value, true
	}

	settingParts := strings.Split(setting, ".")
	value, ok := getSettingFromMap(setting, settingParts, c.loadedSettings)
	if ok {
		c.settings[setting] = value
	}

	return value, ok
}

func settingToEnvVarName(setting string) string {
	return strings.ToUpper(strings.ReplaceAll(setting, ".", "_"))
}

func getSettingFromMap(setting string, settingParts []string, settings map[string]interface{}) (interface{}, bool) {
	if len(settingParts) == 1 {
		if value := settings[settingParts[0]]; value != nil {
			return value, true
		}
	} else if subSettings, ok := settings[settingParts[0]]; ok {
		return getSettingFromMap(setting, settingParts[1:], subSettings.(map[string]interface{}))
	}

	return nil, false
}
