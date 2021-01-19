package main

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"
	"gitlab.com/markomalis/packback-api/pkg/models"
)

func main() {
	loadConfig("config", ".")
	dbString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		viper.GetString("db.host"),
		viper.GetString("db.port"),
		viper.GetString("db.user"),
		viper.GetString("db.db_name"),
		viper.GetString("db.password"),
	)
	db, err := gorm.Open("postgres", dbString)
	if err != nil {
		log.Fatalf("Could not connect to the database %s", err)
	}
	defer db.Close()

	for i := 1000; i < 5000; i++ {
		createPackage(db, i)
	}
}

func createPackage(db *gorm.DB, i int) {
	p := models.Package{
		Name: fmt.Sprintf("Packie %d", i),
		PackageTypeID: 4,
	}
	db.Create(&p)
}

func loadConfig(filename, path string) {
	viper.SetConfigName(filename)
	viper.AddConfigPath(path)
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}
}
