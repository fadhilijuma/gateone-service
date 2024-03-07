package conditiondb

import (
	"fmt"
	"github.com/fadhilijuma/gateone-service/business/core/crud/condition"
	"time"

	"github.com/google/uuid"
)

type dbCondition struct {
	ID          uuid.UUID `db:"condition_id"`
	UserID      uuid.UUID `db:"user_id"`
	Name        string    `db:"name"`
	DateCreated time.Time `db:"date_created"`
	DateUpdated time.Time `db:"date_updated"`
}

func toDBCondition(rl condition.Condition) dbCondition {
	rlDB := dbCondition{
		ID:          rl.ID,
		UserID:      rl.UserID,
		Name:        rl.Name,
		DateCreated: rl.DateCreated.UTC(),
		DateUpdated: rl.DateUpdated.UTC(),
	}

	return rlDB
}

func toCoreCondition(dbCn dbCondition) (condition.Condition, error) {
	rl := condition.Condition{
		ID:          dbCn.ID,
		UserID:      dbCn.UserID,
		Name:        dbCn.Name,
		DateCreated: dbCn.DateCreated.In(time.Local),
		DateUpdated: dbCn.DateUpdated.In(time.Local),
	}

	return rl, nil
}

func toCoreConditionsSlice(dbConditions []dbCondition) ([]condition.Condition, error) {
	conditions := make([]condition.Condition, len(dbConditions))

	for i, dbHme := range dbConditions {
		var err error
		conditions[i], err = toCoreCondition(dbHme)
		if err != nil {
			return nil, fmt.Errorf("parse type: %w", err)
		}
	}

	return conditions, nil
}
