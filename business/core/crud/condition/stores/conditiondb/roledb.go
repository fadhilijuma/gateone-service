// Package conditiondb contains condition related CRUD functionality.
package conditiondb

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/fadhilijuma/gateone-service/business/core/crud/condition"
	"github.com/fadhilijuma/gateone-service/business/data/sqldb"
	"github.com/fadhilijuma/gateone-service/business/data/transaction"
	"github.com/fadhilijuma/gateone-service/business/web/v1/order"
	"github.com/fadhilijuma/gateone-service/foundation/logger"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// Store manages the set of APIs for condition database access.
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
func (s *Store) ExecuteUnderTransaction(tx transaction.Transaction) (condition.Storer, error) {
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

// Create inserts a new condition into the database.
func (s *Store) Create(ctx context.Context, cn condition.Condition) error {
	const q = `
    INSERT INTO conditions
        (condition_id, user_id, name, date_created, date_updated)
    VALUES
        (:condition_id, :user_id, :name, :date_created, :date_updated)`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBCondition(cn)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Delete removes a condition from the database.
func (s *Store) Delete(ctx context.Context, condition condition.Condition) error {
	data := struct {
		ID string `db:"condition_id"`
	}{
		ID: condition.ID.String(),
	}

	const q = `
    DELETE FROM
	    conditions
	WHERE
	  	condition_id = :condition_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Update replaces a condition document in the database.
func (s *Store) Update(ctx context.Context, rl condition.Condition) error {
	const q = `
    UPDATE
        conditions
    SET
        "name"          = :name,
        "date_updated"  = :date_updated
    WHERE
        condition_id = :condition_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBCondition(rl)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Query retrieves a list of existing conditions from the database.
func (s *Store) Query(ctx context.Context, filter condition.QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]condition.Condition, error) {
	data := map[string]interface{}{
		"offset":        (pageNumber - 1) * rowsPerPage,
		"rows_per_page": rowsPerPage,
	}

	const q = `
    SELECT
	    condition_id, user_id, name, date_created, date_updated
	FROM
	  	conditions`

	buf := bytes.NewBufferString(q)
	s.applyFilter(filter, data, buf)

	orderByClause, err := orderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbConditions []dbCondition
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbConditions); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	conditions, err := toCoreConditionsSlice(dbConditions)
	if err != nil {
		return nil, err
	}

	return conditions, nil
}

// Count returns the total number of conditions in the DB.
func (s *Store) Count(ctx context.Context, filter condition.QueryFilter) (int, error) {
	data := map[string]interface{}{}

	const q = `
    SELECT
        count(1)
    FROM
        conditions`

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

// QueryByID gets the specified condition from the database.
func (s *Store) QueryByID(ctx context.Context, conditionID uuid.UUID) (condition.Condition, error) {
	data := struct {
		ID string `db:"condition_id"`
	}{
		ID: conditionID.String(),
	}

	const q = `
    SELECT
	  	condition_id, user_id, name, date_created, date_updated
    FROM
        conditions
    WHERE
        condition_id = :condition_id`

	var dbRl dbCondition
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, q, data, &dbRl); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return condition.Condition{}, fmt.Errorf("namedquerystruct: %w", condition.ErrNotFound)
		}
		return condition.Condition{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	return toCoreCondition(dbRl)
}

// QueryByUserID gets the specified condition from the database by user id.
func (s *Store) QueryByUserID(ctx context.Context, userID uuid.UUID) ([]condition.Condition, error) {
	data := struct {
		ID string `db:"user_id"`
	}{
		ID: userID.String(),
	}

	const q = `
	SELECT
	    condition_id, user_id, name, date_created, date_updated
	FROM
		conditions
	WHERE
		user_id = :user_id`

	var dbConditions []dbCondition
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, q, data, &dbConditions); err != nil {
		return nil, fmt.Errorf("namedquerystruct: %w", err)
	}

	return toCoreConditionsSlice(dbConditions)
}
