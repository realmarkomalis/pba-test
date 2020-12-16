package entities

type RestaurantReport struct {
	Returned           int      `json:"returned"`
	NotReturned        int      `json:"notReturned"`
	TotalCycles        int      `json:"totalCycles"`
	Supplied           int      `json:"supplied"`
	CurrentStock       int      `json:"currentStock"`
	DropOffPointStatus int      `json:"dropOffPointStatus"`
	Returns            []Return `json:"returns"`
}

type SupplierReport struct {
	Returned     int      `json:"returned"`
	NotReturned  int      `json:"notReturned"`
	TotalCycles  int      `json:"totalCycles"`
	Supplied     int      `json:"supplied"`
	AverageLoops int      `json:"averageLoops"`
	Returns      []Return `json:"returns"`
}
