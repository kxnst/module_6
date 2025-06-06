package internal

import "guitar_processor/internal/config"

func GetProviders() []interface{} {
	var provides []interface{}

	provides = append(provides, config.NewConfig)

	return provides
}
