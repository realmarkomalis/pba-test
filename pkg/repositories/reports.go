package repositories

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/markomalis/packback-api/pkg/entities"
	"gitlab.com/markomalis/packback-api/pkg/models"
)

type ReportsRepository struct {
	DB *gorm.DB
}

type InventoryStatus map[string]map[int]int

type ReturnRateData struct {
	status models.ReturnStatus
	count  int
}

func (r ReportsRepository) RestaurantReport(rs entities.Restaurant, startInterval, endInterval string) (*entities.RestaurantReport, error) {
	report := entities.RestaurantReport{}

	to, re, nr, err := r.getReturnRate([]uint{rs.ID}, startInterval, endInterval)
	if err != nil {
		return nil, err
	}

	report.TotalCycles = to
	report.Returned = re
	report.NotReturned = nr

	ss, err := r.getSuppliedStats([]uint{rs.ID}, startInterval, endInterval)
	if err != nil {
		return nil, err
	}
	report.Supplied = ss

	cs, err := r.getCurrentStock(rs.ID)
	if err != nil {
		return nil, err
	}
	report.CurrentStock = cs

	ds, err := r.getDropOffPointStatus(rs.DropOffPoint.ID)
	if err != nil {
		return nil, err
	}
	report.DropOffPointStatus = ds

	rts, err := r.getReturnsInPeriod([]uint{rs.ID}, startInterval, endInterval)
	if err != nil {
		return nil, err
	}
	report.Returns = rts

	return &report, nil
}

func (r ReportsRepository) SupplierRestaurantReport(rs entities.Restaurant, userID uint, startInterval, endInterval string) (*entities.RestaurantReport, error) {
	report := entities.RestaurantReport{}

	to, re, nr, err := r.getSupplierReturnRate([]uint{rs.ID}, userID, startInterval, endInterval)
	if err != nil {
		return nil, err
	}

	report.TotalCycles = to
	report.Returned = re
	report.NotReturned = nr

	ss, err := r.getSupplierSuppliedStats([]uint{rs.ID}, userID, startInterval, endInterval)
	if err != nil {
		return nil, err
	}
	report.Supplied = ss

	cs, err := r.getSupplierCurrentStock(rs.ID, userID)
	if err != nil {
		return nil, err
	}
	report.CurrentStock = cs

	ds, err := r.getSupplierDropOffPointStatus(rs.DropOffPoint.ID, userID)
	if err != nil {
		return nil, err
	}
	report.DropOffPointStatus = ds

	rts, err := r.getSupplierRestaurantReturnsInPeriod([]uint{rs.ID}, userID, startInterval, endInterval)
	if err != nil {
		return nil, err
	}
	report.Returns = rts

	return &report, nil
}

func (r ReportsRepository) SupplierTotalsReport(rs []entities.Restaurant, userID uint, startInterval, endInterval string) (*entities.SupplierReport, error) {
	report := entities.SupplierReport{}
	rIDs := []uint{}
	for _, re := range rs {
		rIDs = append(rIDs, re.ID)
	}

	to, re, nr, err := r.getSupplierReturnRate(rIDs, userID, startInterval, endInterval)
	if err != nil {
		return nil, err
	}

	report.TotalCycles = to
	report.Returned = re
	report.NotReturned = nr

	ss, err := r.getSupplierSuppliedStats(rIDs, userID, startInterval, endInterval)
	if err != nil {
		return nil, err
	}
	report.Supplied = ss

	rts, err := r.getSupplierReturnsInPeriod(userID, startInterval, endInterval)
	if err != nil {
		return nil, err
	}
	report.Returns = rts

	return &report, nil
}

func (r ReportsRepository) getReturnRate(rIDs []uint, startInterval, endInterval string) (int, int, int, error) {
	rows, err := r.DB.Raw(RETURN_RATE_QUERY, rIDs, startInterval, endInterval).Rows()
	if err != nil {
		return 0, 0, 0, err
	}
	defer rows.Close()

	total, returned, notReturned := 0, 0, 0
	for rows.Next() {
		rd := ReturnRateData{}
		rows.Scan(&rd.status, &rd.count)
		if rd.status == models.Collected || rd.status == models.Fulfilled {
			returned++
		} else {
			notReturned++
		}
		total++
	}

	return total, returned, notReturned, nil
}

func (r ReportsRepository) getSupplierReturnRate(rIDs []uint, userID uint, startInterval, endInterval string) (int, int, int, error) {
	rows, err := r.DB.Raw(SUPPLIER_RETURN_RATE_QUERY, rIDs, userID, startInterval, endInterval).Rows()
	if err != nil {
		return 0, 0, 0, err
	}
	defer rows.Close()

	total, returned, notReturned := 0, 0, 0
	for rows.Next() {
		rd := ReturnRateData{}
		rows.Scan(&rd.status, &rd.count)
		if rd.status == models.Collected || rd.status == models.Fulfilled {
			returned++
		} else {
			notReturned++
		}
		total++
	}

	return total, returned, notReturned, nil
}

func (r ReportsRepository) getSuppliedStats(rIDs []uint, startInterval, endInterval string) (int, error) {
	rows, err := r.DB.Raw(SUPPLIES_IN_PERIOD, rIDs, startInterval, endInterval).Rows()
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	total := 0
	for rows.Next() {
		rows.Scan(&total)
	}

	return total, nil
}

func (r ReportsRepository) getSupplierSuppliedStats(rIDs []uint, userID uint, startInterval, endInterval string) (int, error) {
	rows, err := r.DB.Raw(SUPPLIER_SUPPLIES_IN_PERIOD, rIDs, userID, startInterval, endInterval).Rows()
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	total := 0
	for rows.Next() {
		rows.Scan(&total)
	}

	return total, nil
}

func (r ReportsRepository) getCurrentStock(rID uint) (int, error) {
	rows, err := r.DB.Raw(CURRENT_STOCK_QUERY, rID).Rows()
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	total := 0
	for rows.Next() {
		rows.Scan(&total)
	}

	return total, nil
}

func (r ReportsRepository) getSupplierCurrentStock(rID, userID uint) (int, error) {
	rows, err := r.DB.Raw(SUPPLIER_CURRENT_STOCK_QUERY, rID, userID).Rows()
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	total := 0
	for rows.Next() {
		rows.Scan(&total)
	}

	return total, nil
}

func (r ReportsRepository) getDropOffPointStatus(dID uint) (int, error) {
	rows, err := r.DB.Raw(CURRENT_DROP_OFF_POINT_CONTENT_QUERY, dID).Rows()
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	total := 0
	for rows.Next() {
		rows.Scan(&total)
	}

	return total, nil
}

func (r ReportsRepository) getSupplierDropOffPointStatus(dID, userID uint) (int, error) {
	rows, err := r.DB.Raw(SUPPLIER_CURRENT_DROP_OFF_POINT_CONTENT_QUERY, dID, userID).Rows()
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	total := 0
	for rows.Next() {
		rows.Scan(&total)
	}

	return total, nil
}

func (r ReportsRepository) getReturnsInPeriod(rIDs []uint, startInterval, endInterval string) ([]entities.Return, error) {
	ps := []models.PackageSupply{}
	err := r.DB.
		Where("restaurant_id in (?)", rIDs).
		Where("created_at BETWEEN ? AND ?", startInterval, endInterval).
		Find(&ps).
		Error
	if err != nil {
		return nil, err
	}

	ri := []uint{}
	for _, p := range ps {
		ri = append(ri, p.ReturnID)
	}

	rs := []models.Return{}
	err = r.DB.
		Preload("Package").
		Preload("Package.PackageType").
		Where("id in (?)", ri).
		Find(&rs).
		Error
	if err != nil {
		return nil, err
	}

	rts := []entities.Return{}
	for _, r := range rs {
		rts = append(rts, r.ModelToEntity())
	}

	return rts, nil
}

func (r ReportsRepository) getSupplierReturnsInPeriod(userID uint, startInterval, endInterval string) ([]entities.Return, error) {
	ps := []models.PackageSupply{}
	err := r.DB.
		Where("user_id = ?", userID).
		Where("created_at BETWEEN ? AND ?", startInterval, endInterval).
		Find(&ps).
		Error
	if err != nil {
		return nil, err
	}

	ri := []uint{}
	for _, p := range ps {
		ri = append(ri, p.ReturnID)
	}

	rs := []models.Return{}
	err = r.DB.
		Preload("Package").
		Preload("Package.PackageType").
		Where("id in (?)", ri).
		Find(&rs).
		Error
	if err != nil {
		return nil, err
	}

	rts := []entities.Return{}
	for _, r := range rs {
		rts = append(rts, r.ModelToEntity())
	}

	return rts, nil
}

func (r ReportsRepository) getSupplierRestaurantReturnsInPeriod(rIDs []uint, userID uint, startInterval, endInterval string) ([]entities.Return, error) {
	ps := []models.PackageSupply{}
	err := r.DB.
		Where("user_id = ?", userID).
		Where("restaurant_id in (?)", rIDs).
		Where("created_at BETWEEN ? AND ?", startInterval, endInterval).
		Find(&ps).
		Error
	if err != nil {
		return nil, err
	}

	ri := []uint{}
	for _, p := range ps {
		ri = append(ri, p.ReturnID)
	}

	rs := []models.Return{}
	err = r.DB.
		Preload("Package").
		Preload("Package.PackageType").
		Where("id in (?)", ri).
		Find(&rs).
		Error
	if err != nil {
		return nil, err
	}

	rts := []entities.Return{}
	for _, r := range rs {
		rts = append(rts, r.ModelToEntity())
	}

	return rts, nil
}

const RETURNS_IN_PERIOD_QUERY = `
    SELECT id, status
    FROM package_supplies
    JOIN returns on returns.id = package_supplies.return_id
    JOIN packages on packages.id = returns.package_id
    JOIN package_types on package_types.id = packages.package_type_id
    WHERE package_supplies.restaurant_id in (?)
    AND package_supplies.created_at BETWEEN ?::timestamp AND ?::timestamp;
`

const RETURN_RATE_QUERY = `
    SELECT status, COUNT(status)
    FROM package_supplies
    JOIN returns on returns.id = package_supplies.return_id
    WHERE package_supplies.restaurant_id in (?)
    AND returns.status != 0
    AND package_supplies.created_at BETWEEN ?::timestamp AND ?::timestamp
    GROUP BY status;
`

const SUPPLIER_RETURN_RATE_QUERY = `
    SELECT status, return_id, restaurant_id
    FROM package_supplies
    JOIN returns on returns.id = package_supplies.return_id
    WHERE package_supplies.restaurant_id in (?)
    AND package_supplies.user_id = ?
    AND returns.status != 0
    AND package_supplies.created_at BETWEEN ?::timestamp AND ?::timestamp;
`

const AVERAGE_LOOPS_PER_PACKAGE_QUERY = `
    SELECT COUNT(returns.status), returns.status, returns.package_id
    FROM returns
    GROUP BY returns.package_id, returns.status;
`
const SUPPLIES_IN_PERIOD = `
    SELECT count(*)
    FROM package_supplies
    WHERE package_supplies.restaurant_id in (?)
    AND package_supplies.created_at BETWEEN ?::timestamp AND ?::timestamp;
`

const SUPPLIER_SUPPLIES_IN_PERIOD = `
    SELECT count(*)
    FROM package_supplies
    WHERE package_supplies.restaurant_id in (?)
    AND package_supplies.user_id = ?
    AND package_supplies.created_at BETWEEN ?::timestamp AND ?::timestamp;
`

const CURRENT_STOCK_QUERY = `
    SELECT count(*)
    FROM package_supplies
    JOIN returns on returns.id = package_supplies.return_id
    WHERE package_supplies.restaurant_id = ?
    AND returns.status = 0;
`

const SUPPLIER_CURRENT_STOCK_QUERY = `
    SELECT count(*)
    FROM package_supplies
    JOIN returns on returns.id = package_supplies.return_id
    WHERE package_supplies.restaurant_id = ?
    AND package_supplies.user_id = ?
    AND returns.status = 0;
`

const CURRENT_DROP_OFF_POINT_CONTENT_QUERY = `
    SELECT COUNT(drop_off_intents.id)
    FROM drop_off_intents
    JOIN "returns" ON "returns".id = drop_off_intents.return_id
    WHERE drop_off_intents.drop_off_point_id = ?
    AND "returns".status = 5;
`

const SUPPLIER_CURRENT_DROP_OFF_POINT_CONTENT_QUERY = `
    SELECT COUNT(drop_off_intents.id)
    FROM drop_off_intents
    JOIN "returns" ON "returns".id = drop_off_intents.return_id
    JOIN package_supplies ON package_supplies.return_id = drop_off_intents.return_id
    WHERE drop_off_intents.drop_off_point_id = ?
    AND package_supplies.user_id = ?
    AND "returns".status = 5;
`
