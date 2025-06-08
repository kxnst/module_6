package service

func GetProviders() []interface{} {
	var provides []interface{}

	provides = append(provides, NewAudioService)
	provides = append(provides, NewEffectsRegistry)

	return provides
}
