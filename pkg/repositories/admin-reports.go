package repositories

import (
	"gitlab.com/markomalis/packback-api/pkg/entities"
	"gitlab.com/markomalis/packback-api/pkg/models"
)

func (r ReportsRepository) AdminTotalsReport(rs []entities.Restaurant, startInterval, endInterval string) (*entities.SupplierReport, error) {
	report := entities.SupplierReport{}
	rIDs := []uint{}
	for _, re := range rs {
		rIDs = append(rIDs, re.ID)
	}

	to, re, nr, err := r.getAdminReturnRate(rIDs, startInterval, endInterval)
	if err != nil {
		return nil, err
	}

	report.TotalCycles = to
	report.Returned = re
	report.NotReturned = nr

	ss, err := r.getAdminSuppliedStats(rIDs, startInterval, endInterval)
	if err != nil {
		return nil, err
	}
	report.Supplied = ss

	rts, err := r.getAdminReturnsInPeriod(startInterval, endInterval)
	if err != nil {
		return nil, err
	}
	report.Returns = rts

	al, err := r.getAdminAverageLoops()
	if err != nil {
		return nil, err
	}
	report.AverageLoops = al

	return &report, nil
}

func (r ReportsRepository) getAdminReturnRate(rIDs []uint, startInterval, endInterval string) (int, int, int, error) {
	rows, err := r.DB.Raw(ADMIN_RETURN_RATE_QUERY, rIDs, startInterval, endInterval).Rows()
	if err != nil {
		return 0, 0, 0, err
	}
	defer rows.Close()

	total, returned, notReturned := 0, 0, 0
	for rows.Next() {
		rd := ReturnRateData{}
		rows.Scan(&rd.status, &rd.count)
		if rd.status == models.Collected || rd.status == models.Fulfilled {
			returned = returned + rd.count
		} else {
			notReturned = notReturned + rd.count
		}
		total = total + rd.count
	}

	return total, returned, notReturned, nil
}

const ADMIN_RETURN_RATE_QUERY = `
    SELECT status, COUNT(status)
    FROM package_supplies
    JOIN returns on returns.id = package_supplies.return_id
    WHERE package_supplies.restaurant_id in (?)
    AND package_supplies.created_at BETWEEN ?::timestamp AND ?::timestamp
    GROUP BY status;
`

func (r ReportsRepository) getAdminSuppliedStats(rIDs []uint, startInterval, endInterval string) (int, error) {
	rows, err := r.DB.Raw(ADMIN_SUPPLIES_IN_PERIOD, rIDs, startInterval, endInterval).Rows()
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

const ADMIN_SUPPLIES_IN_PERIOD = `
    SELECT count(*)
    FROM package_supplies
    WHERE package_supplies.restaurant_id in (?)
    AND package_supplies.created_at BETWEEN ?::timestamp AND ?::timestamp;
`

func (r ReportsRepository) getAdminAverageLoops() (float64, error) {
	rows, err := r.DB.Raw(ADMIN_AVERAGE_LOOPS_PER_PACKAGE_QUERY).Rows()
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	total, returned, notReturned := 0, 0, 0
	for rows.Next() {
		rd := averageLoopsData{}
		rows.Scan(&rd.count, &rd.status, &rd.packageID)
		if rd.status == models.Collected || rd.status == models.Fulfilled {
			returned = returned + rd.count
		} else {
			notReturned = notReturned + rd.count
		}
		total = total + rd.count
	}

	rows, err = r.DB.Raw(ADMIN_UNIQUE_PACKAGES_FOR_SUPPLIER_QUERY).Rows()
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	up := 0
	for rows.Next() {
		rows.Scan(&up)
	}

	average := float64(total) / float64(up)
	return average, nil
}

const ADMIN_AVERAGE_LOOPS_PER_PACKAGE_QUERY = `
    SELECT COUNT(returns.status), returns.status, returns.package_id
    FROM "returns"
    GROUP BY returns.package_id, returns.status;
`

const ADMIN_UNIQUE_PACKAGES_FOR_SUPPLIER_QUERY = `
    SELECT COUNT(DISTINCT returns.package_id)
    FROM "returns";
`

func (r ReportsRepository) getAdminReturnsInPeriod(startInterval, endInterval string) ([]entities.Return, error) {
	ps := []models.PackageSupply{}
	err := r.DB.
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
