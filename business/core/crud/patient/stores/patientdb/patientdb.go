// Package patientdb contains patient related CRUD functionality.
package patientdb

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/fadhilijuma/gateone-service/business/core/crud/patient"
	"github.com/fadhilijuma/gateone-service/business/data/sqldb"
	"github.com/fadhilijuma/gateone-service/business/data/transaction"
	"github.com/fadhilijuma/gateone-service/business/web/v1/order"
	"github.com/fadhilijuma/gateone-service/foundation/logger"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// Store manages the set of APIs for patient database access.
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
func (s *Store) ExecuteUnderTransaction(tx transaction.Transaction) (patient.Storer, error) {
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

// Create adds a Patient to the sqldb. It returns the created Patient with
// fields like ID and DateCreated populated.
func (s *Store) Create(ctx context.Context, prd patient.Patient) error {
	const q = `
	INSERT INTO patients
		(patient_id, user_id, name, age, condition, healed,video_links, date_created, date_updated)
	VALUES
		(:patient_id, :user_id, :name, :age, :condition, :healed,:video_links, :date_created, :date_updated)`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBPatient(prd)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Update modifies data about a Patient. It will error if the specified ID is
// invalid or does not reference an existing Patient.
func (s *Store) Update(ctx context.Context, prd patient.Patient) error {
	const q = `
	UPDATE
		patients
	SET
		"name" = :name,
		"age" = :age,
		"condition" = :condition,
		"healed" = :healed,
		"video_links" = :video_links,
		"date_updated" = :date_updated
	WHERE
		patient_id = :patient_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBPatient(prd)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Delete removes the patient identified by a given ID.
func (s *Store) Delete(ctx context.Context, prd patient.Patient) error {
	data := struct {
		ID string `db:"patient_id"`
	}{
		ID: prd.ID.String(),
	}

	const q = `
	DELETE FROM
		patients
	WHERE
		patient_id = :patient_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Query gets all Patients from the database.
func (s *Store) Query(ctx context.Context, filter patient.QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]patient.Patient, error) {
	data := map[string]interface{}{
		"offset":        (pageNumber - 1) * rowsPerPage,
		"rows_per_page": rowsPerPage,
	}

	const q = `
	SELECT
	    patient_id, user_id, name, age, condition, healed, video_links, date_created, date_updated
	FROM
		patients`

	buf := bytes.NewBufferString(q)
	s.applyFilter(filter, data, buf)

	orderByClause, err := orderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbPrds []dbPatient
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbPrds); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toCorePatients(dbPrds), nil
}

// Count returns the total number of Patients in the DB.
func (s *Store) Count(ctx context.Context, filter patient.QueryFilter) (int, error) {
	data := map[string]interface{}{}

	const q = `
	SELECT
		count(1)
	FROM
		patients`

	buf := bytes.NewBufferString(q)
	s.applyFilter(filter, data, buf)

	var count struct {
		Count  int `db:"count"`
		Healed int `db:"healed"`
	}
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &count); err != nil {
		return 0, fmt.Errorf("namedquerystruct: %w", err)
	}

	return count.Count, nil
}

// QueryByID finds the patient identified by a given ID.
func (s *Store) QueryByID(ctx context.Context, patientID uuid.UUID) (patient.Patient, error) {
	data := struct {
		ID string `db:"patient_id"`
	}{
		ID: patientID.String(),
	}

	const q = `
	SELECT
	    patient_id, user_id, name, age, condition, healed, video_links, date_created, date_updated
	FROM
		patients
	WHERE
		patient_id = :patient_id`

	var dbPn dbPatient
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, q, data, &dbPn); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return patient.Patient{}, fmt.Errorf("namedquerystruct: %w", patient.ErrNotFound)
		}
		return patient.Patient{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	return toCorePatient(dbPn), nil
}

// QueryByUserID finds the patient identified by a given User ID.
func (s *Store) QueryByUserID(ctx context.Context, userID uuid.UUID) ([]patient.Patient, error) {
	data := struct {
		ID string `db:"user_id"`
	}{
		ID: userID.String(),
	}

	const q = `
	SELECT
	    patient_id, user_id, name, age, condition, healed, video_links, date_created, date_updated
	FROM
		patients
	WHERE
		user_id = :user_id`

	var dbPrds []dbPatient
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, q, data, &dbPrds); err != nil {
		return nil, fmt.Errorf("namedquerystruct: %w", err)
	}

	return toCorePatients(dbPrds), nil
}
