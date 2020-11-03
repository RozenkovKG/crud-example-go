package main

import (
	"crud-example-go/src/config"
	dbImpl "crud-example-go/src/db"
	"crud-example-go/src/intreface/repository"
	repository2 "crud-example-go/src/repository"
	"crud-example-go/src/server"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	conf, err := config.New("properties.yml")
	if err != nil {
		panic(err)
	}

	userRepository, err := getRepository(&conf.DB)
	if err != nil {
		panic(err)
	}

	err = server.New(conf.Port, userRepository).Start()
	if err != nil {
		panic(err)
	}
}

func getRepository(config *config.DbConfig) (repository.UserRepository, error) {
	dsn := "host=" + config.Host +
		" port=" + config.Port +
		" dbname=" + config.Name +
		" user=" + config.User +
		" password=" + config.Pass
	if db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{}); err != nil {
		return nil, err
	} else {
		if err := migrateSchema(db); err != nil {
			return nil, err
		}
		return &repository2.UserRepository{
			DB: &dbImpl.DbWrapper{DB: db}}, nil
	}
}

func migrateSchema(db *gorm.DB) error {
	if err := db.AutoMigrate(&repository2.User{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&repository2.Tag{}); err != nil {
		return err
	}
	return nil
}
