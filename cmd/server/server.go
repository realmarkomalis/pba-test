package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/rs/cors"
	"github.com/spf13/viper"

	"gitlab.com/markomalis/packback-api/pkg/models"
	"gitlab.com/markomalis/packback-api/pkg/routes"
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

	r := mux.NewRouter()
	routes.InitializeRoutes(r, db)
	db.AutoMigrate(
		&models.User{},
		&models.UserAddress{},
		&models.UserRole{},
		&models.Restaurant{},
		&models.Package{},
		&models.PackageType{},
		&models.Return{},
		&models.UserReturnEntry{},
		&models.UserReturnEntryReturns{},
		&models.PackageSupply{},
		&models.PackageDispatch{},
		&models.PackagePickup{},
		&models.ReturnRequest{},
		&models.PickupSlot{},
		&models.DropOffPoint{},
		&models.DropOffIntent{},
		&models.DropOffCollect{},
		&models.CustomerDepositPayout{},
	)

	db.Model(&models.User{}).AddForeignKey("user_role_id", "user_roles(id)", "RESTRICT", "RESTRICT")

	db.Model(&models.UserAddress{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")

	db.Model(&models.Restaurant{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")

	db.Model(&models.Return{}).AddForeignKey("package_id", "packages(id)", "RESTRICT", "RESTRICT")

	db.Model(&models.Package{}).AddForeignKey("package_type_id", "package_types(id)", "RESTRICT", "RESTRICT")

	db.Model(&models.PackageSupply{}).AddForeignKey("return_id", "returns(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.PackageSupply{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.PackageSupply{}).AddForeignKey("restaurant_id", "restaurants(id)", "RESTRICT", "RESTRICT")

	db.Model(&models.PackageDispatch{}).AddForeignKey("return_id", "returns(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.PackageDispatch{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")

	db.Model(&models.ReturnRequest{}).AddForeignKey("return_id", "returns(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.ReturnRequest{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.ReturnRequest{}).AddForeignKey("pickup_slot_id", "pickup_slots(id)", "RESTRICT", "RESTRICT")

	db.Model(&models.PackagePickup{}).AddForeignKey("return_id", "returns(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.PackagePickup{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")

	db.Model(&models.UserReturnEntry{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.UserReturnEntryReturns{}).AddForeignKey("user_return_entry_id", "user_return_entries(id)", "CASCADE", "CASCADE")
	db.Model(&models.UserReturnEntryReturns{}).AddForeignKey("return_id", "returns(id)", "CASCADE", "CASCADE")

	db.Model(&models.DropOffIntent{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.DropOffIntent{}).AddForeignKey("return_id", "returns(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.DropOffIntent{}).AddForeignKey("drop_off_point_id", "drop_off_points(id)", "CASCADE", "CASCADE")

	db.Model(&models.DropOffCollect{}).AddForeignKey("return_id", "returns(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.DropOffCollect{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.DropOffCollect{}).AddForeignKey("drop_off_point_id", "drop_off_points(id)", "CASCADE", "CASCADE")

	db.Model(&models.CustomerDepositPayout{}).AddForeignKey("return_id", "returns(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.CustomerDepositPayout{}).AddForeignKey("receiving_user_id", "users(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.CustomerDepositPayout{}).AddForeignKey("paying_user_id", "users(id)", "RESTRICT", "RESTRICT")

	c := cors.New(cors.Options{
		AllowedOrigins:   viper.GetStringSlice("allowed_origins"),
		AllowCredentials: true,
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodOptions},
		AllowedHeaders:   []string{"Content-Type", "X-Packback-Auth"},
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})

	srv := &http.Server{
		Handler: c.Handler(r),
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func loadConfig(filename, path string) {
	viper.SetConfigName(filename)
	viper.AddConfigPath(path)
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}
}
