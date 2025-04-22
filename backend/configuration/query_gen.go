package main

import (
	model "github.com/cat-spmog/iothubgo/models"
	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
)

// Dynamic SQL
type Querier interface {
	// SELECT * FROM @@table WHERE name = @name{{if role !=""}} AND role = @role{{end}}
	FilterWithNameAndRole(name, role string) ([]gen.T, error)
	// SELECT * FROM @@table WHERE id=@id
	GetByID(id int) (gen.T, error) // GetByID query data by id and return it as *struct*
	// SELECT * FROM @@table WHERE role IN @rolesName
	GetByRoles(rolesName ...string) ([]*gen.T, error)
}

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "../query",
		Mode:    gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	dsn := "host=localhost user=postgres dbname=iotgo port=5432 sslmode=disable"
	gormdb, err := gorm.Open(postgres.New(postgres.Config{DSN: dsn}), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to the database: " + err.Error())
	}
	g.UseDB(gormdb) // reuse your gorm db

	// Generate basic type-safe DAO API for struct `model.User` following conventions
	g.ApplyBasic(&model.DataSink{}, &model.DataSource{}, &model.Device{}, &model.SensorData{}, &model.User{}, &model.UserInteraction{}) // Get all models from DBHandler

	// Generate Type Safe API with Dynamic SQL defined on Querier interface for `model.User` and `model.Company`
	g.ApplyInterface(func(Querier) {}, &model.User{}, &model.Device{}, &model.DataSink{}, &model.DataSource{}, &model.SensorData{}, &model.UserInteraction{})

	// Generate the code
	g.Execute()
}
