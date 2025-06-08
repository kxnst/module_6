package repository

func GetProviders() []interface{} {
	var provides []interface{}

	provides = append(provides, NewEffectRepository)
	provides = append(provides, NewUserRepository)

	return provides
}
