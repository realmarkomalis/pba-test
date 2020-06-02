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

	slots := []int{11, 12, 17, 18}
	now := time.Now()
	for i := 0; i < 3; i++ {
		for _, s := range slots {
			d := time.Date(now.Year(), now.Month(), now.Day(), s, 0, 0, 0, time.UTC)
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
