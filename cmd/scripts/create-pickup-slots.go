package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"
	"gitlab.com/markomalis/packback-api/pkg/models"
)

func main() {
	loadConfigS("config", ".")
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

	slots := []int{11, 11, 11, 12, 12, 12, 13, 13, 13, 14, 14, 14, 15, 15, 15, 16, 16, 16, 17, 17, 17, 17, 17, 17, 18, 18, 18, 18, 18, 18, 19, 19, 19, 19, 19, 19, 20, 20, 20, 20, 20, 20, 20, 20, 20, 21, 21, 21, 21, 21, 21, 21, 21, 21}
	now := time.Now()
	for i := 0; i < 20; i++ {
		for _, s := range slots {
			d := time.Date(now.Year(), now.Month(), now.Day(), s, 0, 0, 0, now.Local().Location())
			d = d.AddDate(0, 0, i)

			fmt.Printf("start = %v\n", d)
			fmt.Printf("end = %v\n", d.Add(time.Hour*1))
			createSlot(db, d)
		}
	}
}

func createSlot(db *gorm.DB, d time.Time) {
	slot := models.PickupSlot{
		StartDateTime: d,
		EndDateTime:   d.Add(time.Hour * 1),
	}
	db.Create(&slot)
}

func loadConfigS(filename, path string) {
	viper.SetConfigName(filename)
	viper.AddConfigPath(path)
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}
}
