package model

type ModelHandler interface {
	GetAllModels() []interface{}
}

// GetAllModels returns a slice of all models in the system
// This is used for database migrations and seeding
// It is important to keep this list updated with all models in the system

// Do this dynamically to avoid having to update this list manually
// when adding new models

func GetAllModels() []interface{} {
	// This function creates new instances of each model and returns them
	// The reflection approach would be more dynamic but requires more complex code

	// For now, we'll use a semi-dynamic approach where we initialize models
	// based on a registry of model constructors
	modelConstructors := []func() interface{}{
		func() interface{} { return &User{} },
		func() interface{} { return &UserInteraction{} },
		func() interface{} { return &Role{} },
		func() interface{} { return &Device{} },
		func() interface{} { return &SensorData{} },
		func() interface{} { return &DataSink{} },
		func() interface{} { return &DataSource{} },
		// Add new model constructors here
	}

	models := make([]interface{}, 0, len(modelConstructors))
	for _, constructor := range modelConstructors {
		models = append(models, constructor())
	}

	return models
}
