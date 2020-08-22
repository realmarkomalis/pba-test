package repositories

import (
	"github.com/jinzhu/gorm"
)

type ReportsRepository struct {
	DB *gorm.DB
}

type InventoryStatus map[string]map[int]int

func (r ReportsRepository) InventoryPerStatusPerRestaurant() (map[string]map[int]int, error) {

	rows, err := r.DB.Raw(`
		SELECT restaurants.name, status, count(status) 
		FROM returns 
		JOIN package_supplies ON returns.id = package_supplies.return_id
		JOIN restaurants ON restaurants.id = package_supplies.restaurant_id
		GROUP BY restaurants.name, returns.status;
		`,
	).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var statuses = make(InventoryStatus)
	result := struct {
		name   string
		status int
		count  int
	}{}
	for rows.Next() {
		rows.Scan(
			&result.name, &result.status, &result.count,
		)
		if _, ok := statuses[result.name]; !ok {
			statuses[result.name] = map[int]int{}
		}
		statuses[result.name][result.status] = result.count
	}

	return statuses, nil
}
