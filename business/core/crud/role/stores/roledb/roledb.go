// Package roledb contains role related CRUD functionality.
package roledb

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/fadhilijuma/gateone-service/business/core/crud/role"
	"github.com/fadhilijuma/gateone-service/business/data/sqldb"
	"github.com/fadhilijuma/gateone-service/business/data/transaction"
	"github.com/fadhilijuma/gateone-service/business/web/v1/order"
	"github.com/fadhilijuma/gateone-service/foundation/logger"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// Store manages the set of APIs for role database access.
type Store struct {
	log *logger.Logger
	db  sqlx.ExtContext
}

// NewStore constructs the api for data access.
func NewStore(log *logger.Logger, db *sqlx.DB) *Store {
	return &Store{
		log: log,
		db:  db,
	}
}

// ExecuteUnderTransaction constructs a new Store value replacing the sqlx DB
// value with a sqlx DB value that is currently inside a transaction.
func (s *Store) ExecuteUnderTransaction(tx transaction.Transaction) (role.Storer, error) {
	ec, err := sqldb.GetExtContext(tx)
	if err != nil {
		return nil, err
	}

	store := Store{
		log: s.log,
		db:  ec,
	}

	return &store, nil
}

// Create inserts a new role into the database.
func (s *Store) Create(ctx context.Context, rl role.Role) error {
	const q = `
    INSERT INTO roles
        (role_id, user_id, name, date_created, date_updated)
    VALUES
        (:role_id, :user_id, :name, :date_created, :date_updated)`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBRole(rl)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Delete removes a role from the database.
func (s *Store) Delete(ctx context.Context, role role.Role) error {
	data := struct {
		ID string `db:"role_id"`
	}{
		ID: role.ID.String(),
	}

	const q = `
    DELETE FROM
	    roles
	WHERE
	  	role_id = :role_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Update replaces a role document in the database.
func (s *Store) Update(ctx context.Context, rl role.Role) error {
	const q = `
    UPDATE
        roles
    SET
        "name"          = :name,
        "date_updated"  = :date_updated
    WHERE
        role_id = :role_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBRole(rl)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Query retrieves a list of existing roles from the database.
func (s *Store) Query(ctx context.Context, filter role.QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]role.Role, error) {
	data := map[string]interface{}{
		"offset":        (pageNumber - 1) * rowsPerPage,
		"rows_per_page": rowsPerPage,
	}

	const q = `
    SELECT
	    role_id, user_id, name, date_created, date_updated
	FROM
	  	roles`

	buf := bytes.NewBufferString(q)
	s.applyFilter(filter, data, buf)

	orderByClause, err := orderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbRoles []dbRole
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbRoles); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	roles, err := toCoreRolesSlice(dbRoles)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

// Count returns the total number of roles in the DB.
func (s *Store) Count(ctx context.Context, filter role.QueryFilter) (int, error) {
	data := map[string]interface{}{}

	const q = `
    SELECT
        count(1)
    FROM
        roles`

	buf := bytes.NewBufferString(q)
	s.applyFilter(filter, data, buf)

	var count struct {
		Count int `db:"count"`
	}
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &count); err != nil {
		return 0, fmt.Errorf("namedquerystruct: %w", err)
	}

	return count.Count, nil
}

// QueryByID gets the specified role from the database.
func (s *Store) QueryByID(ctx context.Context, roleID uuid.UUID) (role.Role, error) {
	data := struct {
		ID string `db:"role_id"`
	}{
		ID: roleID.String(),
	}

	const q = `
    SELECT
	  	role_id, user_id, name, date_created, date_updated
    FROM
        roles
    WHERE
        role_id = :role_id`

	var dbRl dbRole
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, q, data, &dbRl); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return role.Role{}, fmt.Errorf("namedquerystruct: %w", role.ErrNotFound)
		}
		return role.Role{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	return toCoreRole(dbRl)
}

// QueryByUserID gets the specified role from the database by user id.
func (s *Store) QueryByUserID(ctx context.Context, userID uuid.UUID) ([]role.Role, error) {
	data := struct {
		ID string `db:"user_id"`
	}{
		ID: userID.String(),
	}

	const q = `
	SELECT
	    role_id, user_id, name, date_created, date_updated
	FROM
		roles
	WHERE
		user_id = :user_id`

	var dbRoles []dbRole
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, q, data, &dbRoles); err != nil {
		return nil, fmt.Errorf("namedquerystruct: %w", err)
	}

	return toCoreRolesSlice(dbRoles)
}
