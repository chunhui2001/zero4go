package config

import (
	"strings"

	"github.com/spf13/viper"
)

func readProperties(v *viper.Viper, responseMap map[string]any) {
	config := responseMap["configurations"].(map[string]any)

	for key, val := range config {
		v.SetDefault(strings.TrimSpace(key), strings.TrimSpace(val.(string)))
	}
}
