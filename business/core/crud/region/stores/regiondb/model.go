package regiondb

import (
	"fmt"
	"github.com/fadhilijuma/gateone-service/business/core/crud/region"
	"time"

	"github.com/google/uuid"
)

type dbRegion struct {
	ID          uuid.UUID `db:"Region_id"`
	UserID      uuid.UUID `db:"user_id"`
	Name        string    `db:"name"`
	DateCreated time.Time `db:"date_created"`
	DateUpdated time.Time `db:"date_updated"`
}

func toDBRegion(rl region.Region) dbRegion {
	rlDB := dbRegion{
		ID:          rl.ID,
		UserID:      rl.UserID,
		Name:        rl.Name,
		DateCreated: rl.DateCreated.UTC(),
		DateUpdated: rl.DateUpdated.UTC(),
	}

	return rlDB
}

func toCoreRegion(dbCn dbRegion) (region.Region, error) {
	rl := region.Region{
		ID:          dbCn.ID,
		UserID:      dbCn.UserID,
		Name:        dbCn.Name,
		DateCreated: dbCn.DateCreated.In(time.Local),
		DateUpdated: dbCn.DateUpdated.In(time.Local),
	}

	return rl, nil
}

func toCoreRegionsSlice(dbRegions []dbRegion) ([]region.Region, error) {
	Regions := make([]region.Region, len(dbRegions))

	for i, dbHme := range dbRegions {
		var err error
		Regions[i], err = toCoreRegion(dbHme)
		if err != nil {
			return nil, fmt.Errorf("parse type: %w", err)
		}
	}

	return Regions, nil
}
