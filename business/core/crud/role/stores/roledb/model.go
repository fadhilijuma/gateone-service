package roledb

import (
	"fmt"
	"github.com/fadhilijuma/gateone-service/business/core/crud/role"
	"time"

	"github.com/google/uuid"
)

type dbRole struct {
	ID          uuid.UUID `db:"role_id"`
	UserID      uuid.UUID `db:"user_id"`
	Name        string    `db:"name"`
	DateCreated time.Time `db:"date_created"`
	DateUpdated time.Time `db:"date_updated"`
}

func toDBRole(rl role.Role) dbRole {
	rlDB := dbRole{
		ID:          rl.ID,
		UserID:      rl.UserID,
		Name:        rl.Name,
		DateCreated: rl.DateCreated.UTC(),
		DateUpdated: rl.DateUpdated.UTC(),
	}

	return rlDB
}

func toCoreRole(dbRl dbRole) (role.Role, error) {
	rl := role.Role{
		ID:          dbRl.ID,
		UserID:      dbRl.UserID,
		Name:        dbRl.Name,
		DateCreated: dbRl.DateCreated.In(time.Local),
		DateUpdated: dbRl.DateUpdated.In(time.Local),
	}

	return rl, nil
}

func toCoreRolesSlice(dbRoles []dbRole) ([]role.Role, error) {
	roles := make([]role.Role, len(dbRoles))

	for i, dbHme := range dbRoles {
		var err error
		roles[i], err = toCoreRole(dbHme)
		if err != nil {
			return nil, fmt.Errorf("parse type: %w", err)
		}
	}

	return roles, nil
}
