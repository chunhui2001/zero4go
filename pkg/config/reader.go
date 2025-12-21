package config

import (
	"strings"

	"github.com/spf13/viper"
)

func readProperties(v *viper.Viper, responseMap map[string]any) {
	if config := responseMap["configurations"]; config != nil {
		for key, val := range config.(map[string]any) {
			v.SetDefault(strings.TrimSpace(key), strings.TrimSpace(val.(string)))
		}
	}
}
