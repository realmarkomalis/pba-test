package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	"gitlab.com/markomalis/packback-api/pkg/handlers"
	"gitlab.com/markomalis/packback-api/pkg/middleware"
)

func InitializeRoutes(r *mux.Router, db *gorm.DB) {
	restAPIRouter := r.PathPrefix("/api/rest").Subrouter()

	MeHandler := handlers.MeHandler{DB: db}
	meRouter := restAPIRouter.PathPrefix("/me").Subrouter()
	meRouter.Use(middleware.AuthMiddleware)
	meRouter.HandleFunc("/", MeHandler.GetMe).Methods(http.MethodGet, http.MethodOptions)

	LoginRequestsHandler := handlers.LoginRequestsHandler{DB: db}
	restAPIRouter.HandleFunc("/login-requests/", LoginRequestsHandler.Create).Methods(http.MethodPost, http.MethodOptions)

	LoginAttemptsHandler := handlers.LoginAttemptsHandler{DB: db}
	restAPIRouter.HandleFunc("/login-attempts/", LoginAttemptsHandler.Create).Methods(http.MethodPost, http.MethodOptions)

	ReturnsRequestsHandler := handlers.ReturnsRequestsHandler{DB: db}
	returnsRouter := restAPIRouter.PathPrefix("/returns").Subrouter()
	returnsRouter.Use(middleware.AuthMiddleware)
	returnsRouter.HandleFunc("/", ReturnsRequestsHandler.GetReturns).Methods(http.MethodGet, http.MethodOptions)
	returnsRouter.HandleFunc("/", ReturnsRequestsHandler.Create).Methods(http.MethodPost, http.MethodOptions)

	returnsCountRouter := restAPIRouter.PathPrefix("/returns/count").Subrouter()
	returnsCountRouter.HandleFunc("/", ReturnsRequestsHandler.ReturnsCount).Methods(http.MethodGet, http.MethodOptions)

	packagesHandler := handlers.PackagesHandler{DB: db}
	packageReturnsRouter := restAPIRouter.PathPrefix("/packages").Subrouter()
	packageReturnsRouter.Use(middleware.AuthMiddleware)
	// packageReturnsRouter.HandleFunc("/{id:[0-9]+}/package-dispatches/", ReturnsRequestsHandler.GetReturns).Methods(http.MethodPost, http.MethodOptions)
	// packageReturnsRouter.HandleFunc("/{id:[0-9]+}/pickup-requests/", ReturnsRequestsHandler.CreateReturnRequest).Methods(http.MethodPost, http.MethodOptions)
	packageReturnsRouter.
		HandleFunc("/{package_id:[0-9]+}/return-requests/", ReturnsRequestsHandler.CreateReturnRequest).
		Methods(http.MethodPost, http.MethodOptions)
	packageReturnsRouter.
		HandleFunc("/{package_id:[0-9]+}/package-pickups/", ReturnsRequestsHandler.CreatePackagePickup).
		Methods(http.MethodPost, http.MethodOptions)
	packageReturnsRouter.
		HandleFunc("/{package_id:[0-9]+}/codes/", packagesHandler.GetPackageCode).
		Methods(http.MethodGet, http.MethodOptions)

	psh := handlers.PickupSlotsHandler{DB: db}
	pickupSlotsRouter := restAPIRouter.PathPrefix("/pickup-slots/").Subrouter()
	pickupSlotsRouter.
		HandleFunc("/", psh.GetPickupSlots).
		Methods(http.MethodGet, http.MethodOptions)

	pickupShedulesHandler := handlers.PickupShedulesHandler{DB: db}
	pickupShedulesRouter := restAPIRouter.PathPrefix("/pickup-shedules").Subrouter()
	pickupShedulesRouter.Use(middleware.AuthMiddleware)
	pickupShedulesRouter.
		HandleFunc("/", pickupShedulesHandler.GetPickupShedules).
		Methods(http.MethodGet, http.MethodOptions)

	userAddressRouter := restAPIRouter.PathPrefix("/user-addresses").Subrouter()
	userAddressRouter.Use(middleware.AuthMiddleware)
	uah := handlers.UserAddressesHandler{DB: db}
	userAddressRouter.
		HandleFunc("/", uah.CreateUserAddress).
		Methods(http.MethodPost, http.MethodOptions)
	userAddressRouter.
		HandleFunc("/", uah.UpdateUserAddress).
		Methods(http.MethodPut, http.MethodOptions)

	userChipsRouter := restAPIRouter.PathPrefix("/chips").Subrouter()
	userChipsRouter.Use(middleware.AuthMiddleware)
	uch := handlers.UserChipsHandler{DB: db}
	userChipsRouter.
		HandleFunc("/", uch.GetUserChips).
		Methods(http.MethodGet, http.MethodOptions)

	restaurantsRouter := restAPIRouter.PathPrefix("/restaurants").Subrouter()
	rh := handlers.RestaurantsHandler{DB: db}
	restaurantsRouter.
		HandleFunc("/", rh.GetRestaurants).
		Methods(http.MethodGet, http.MethodOptions)

	reportsRouter := restAPIRouter.PathPrefix("/reports").Subrouter()
	reportsRouter.Use(middleware.AuthMiddleware)
	rph := handlers.ReportsHandler{DB: db}
	reportsRouter.
		HandleFunc("/", rph.StatusesPerRestaurant).
		Methods(http.MethodGet, http.MethodOptions)
}
