package db

func GetProviders() []interface{} {
	var provides []interface{}

	provides = append(provides, NewMongoClient)
	provides = append(provides, NewMongoDatabase)

	return provides
}
