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
	meRouter.HandleFunc("/", MeHandler.GetMe).
		Methods(http.MethodGet, http.MethodOptions)
	meRouter.HandleFunc("/", MeHandler.UpdateMe).
		Methods(http.MethodPut, http.MethodOptions)
	meRouter.HandleFunc("/qr-codes/", MeHandler.GetMyQR).
		Methods(http.MethodGet, http.MethodOptions)

	myRouter := restAPIRouter.PathPrefix("/my").Subrouter()
	myRouter.Use(middleware.AuthMiddleware)

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
	packagesRouter := restAPIRouter.PathPrefix("/packages").Subrouter()
	packagesRouter.
		HandleFunc("/{package_id}/", packagesHandler.GetPackage).
		Methods(http.MethodGet, http.MethodOptions)

	packageTypesHandler := handlers.PackageTypesHandler{DB: db}
	packageTypesRouter := restAPIRouter.PathPrefix("/package-types").Subrouter()
	packageTypesRouter.Use(middleware.AuthMiddleware)
	packageTypesRouter.
		HandleFunc("/", packageTypesHandler.GetPackageTypes).
		Methods(http.MethodGet, http.MethodOptions)

	packageReturnsRouter := restAPIRouter.PathPrefix("/packages").Subrouter()
	packageReturnsRouter.Use(middleware.AuthMiddleware)

	packageReturnsRouter.
		HandleFunc("/", packagesHandler.CreateOrUpdatePackage).
		Methods(http.MethodPut, http.MethodOptions)

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
		HandleFunc("/restaurants/", rph.ListReportableRestaurants).
		Methods(http.MethodGet, http.MethodOptions)
	reportsRouter.
		HandleFunc("/restaurants/", rph.GetRestaurantReport).
		Methods(http.MethodPost, http.MethodOptions)

	dih := handlers.DropOffIntentsHandler{DB: db}
	dropOffIntentsRouter := restAPIRouter.PathPrefix("/drop-off-intents/").Subrouter()
	dropOffIntentsRouter.Use(middleware.AuthMiddleware)
	dropOffIntentsRouter.
		HandleFunc("/", dih.CreateDropOffIntent).
		Methods(http.MethodPost, http.MethodOptions)
	dropOffIntentsRouter.
		HandleFunc("/", dih.ListDropOffIntents).
		Methods(http.MethodGet, http.MethodOptions)

	doh := handlers.DropOffPointsHandler{DB: db}
	myRouter.HandleFunc("/drop-off-points/", doh.ListMyDropOffPoints).
		Methods(http.MethodGet, http.MethodOptions)
	dropOffPointsRouter := restAPIRouter.PathPrefix("/drop-off-points/").Subrouter()
	dropOffPointsRouter.
		HandleFunc("/", doh.ListDropOffPoints).
		Methods(http.MethodGet, http.MethodOptions)
	dropOffPointsRouter.
		HandleFunc("/{dop_id:[0-9]+}/", doh.GetDropOffPoint).
		Methods(http.MethodGet, http.MethodOptions)

	dch := handlers.DropOffCollectsHandler{DB: db}
	dropOffCollectsRouter := restAPIRouter.PathPrefix("/drop-off-collects/").Subrouter()
	dropOffCollectsRouter.Use(middleware.AuthMiddleware)
	dropOffCollectsRouter.
		HandleFunc("/", dch.CreateDropOffCollect).
		Methods(http.MethodPost, http.MethodOptions)

	pch := handlers.PackageCyclesHandler{DB: db}
	packageCyclesRouter := restAPIRouter.PathPrefix("/packages/").Subrouter()
	packageCyclesRouter.Use(middleware.AuthMiddleware)
	packageCyclesRouter.
		HandleFunc("/{package_id:[0-9]+}/cycles/", pch.GetPackageCycles).
		Methods(http.MethodGet, http.MethodOptions)

	dph := handlers.DepositsHandler{DB: db}
	depositsRouter := restAPIRouter.PathPrefix("/deposit-items").Subrouter()
	depositsRouter.Use(middleware.AuthMiddleware)
	depositsRouter.
		HandleFunc("/", dph.GetUserDepositItems).
		Methods(http.MethodGet, http.MethodOptions)

	customerDepositPayoutsRouter := restAPIRouter.PathPrefix("/customer-deposit-payouts").Subrouter()
	customerDepositPayoutsRouter.Use(middleware.AuthMiddleware)
	customerDepositPayoutsRouter.
		HandleFunc("/", dph.CreateCustomerDepositPayout).
		Methods(http.MethodPost, http.MethodOptions)
	customerDepositPayoutsRouter.
		HandleFunc("/", dph.ListCustomerDepositPayouts).
		Methods(http.MethodGet, http.MethodOptions)
	customerDepositPayoutsRouter.
		HandleFunc("/", dph.UpdateCustomerDepositPayout).
		Methods(http.MethodPut, http.MethodOptions)
}
