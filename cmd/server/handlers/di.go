package handlers

import "guitar_processor/cmd/server/middlewares"

func GetProviders() []interface{} {
	var provides []interface{}

	provides = append(provides, NewAuthHandler)
	provides = append(provides, NewDevicesHandler)
	provides = append(provides, NewWebSocketHandler)
	provides = append(provides, NewEffectsHandler)
	provides = append(provides, middlewares.NewAuthMiddleware)

	return provides
}
