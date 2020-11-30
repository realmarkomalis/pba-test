```
package repositories

import (
	"fmt"

	"github.com/jinzhu/gorm"

	"gitlab.com/markomalis/packback-api/pkg/entities"
)

type TestRepository struct {
	DB *gorm.DB
}

func (r TestRepository) Test() (*entities.PickupSlot, error) {
	fmt.Println("TEST REPO TEST")
	rows, err := r.DB.Raw(`
		SELECT
			r.id, r.created_at, r.status,
			p.id, p.name, p.package_code,
			rr.id, rr.created_at, rr.user_id
		FROM returns AS r
		JOIN return_requests AS rr ON r.id = rr.return_id
		JOIN pickup_slots AS ps ON rr.pickup_slot_id = ps.id
		JOIN packages AS p ON r.package_id = p.id
		`,
	).Rows()
	if err != nil {
		fmt.Println("TEST REPO TEST ERRROR")
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	fmt.Println("TEST REPO TEST SUCCES")
	fmt.Println(rows.Columns())
	ret := entities.Return{}
	for rows.Next() {
		rows.Scan(
			&ret.ID, &ret.CreatedAt, &ret.Status,
			&ret.Package.ID, &ret.Package.Name, &ret.Package.PackageCode,
			&ret.PickupRequest.ID, &ret.PickupRequest.CreatedAt, &ret.PickupRequest.User.ID,
		)
		fmt.Println("Scan into entities return")
		fmt.Printf("%+v\n", ret)
	}

	return &entities.PickupSlot{}, nil
}
```
