package internal

import (
	"guitar_processor/internal/config"
	"guitar_processor/internal/db"
	"guitar_processor/internal/effect/service"
	"guitar_processor/internal/repository"
)

func GetProviders() []interface{} {
	var provides []interface{}

	provides = append(provides, config.NewConfig)
	provides = append(provides, db.GetProviders()...)
	provides = append(provides, repository.GetProviders()...)
	provides = append(provides, service.GetProviders()...)

	return provides
}
