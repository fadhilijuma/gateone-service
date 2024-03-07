package patientdb

import (
	"github.com/fadhilijuma/gateone-service/business/core/crud/patient"
	"time"

	"github.com/google/uuid"
)

type dbPatient struct {
	ID          uuid.UUID `db:"patient_id"`
	UserID      uuid.UUID `db:"user_id"`
	Name        string    `db:"name"`
	Age         int       `db:"age"`
	Condition   string    `db:"condition"`
	Healed      bool      `db:"healed"`
	VideoLinks  []string  `db:"video_links"`
	DateCreated time.Time `db:"date_created"`
	DateUpdated time.Time `db:"date_updated"`
}

func toDBPatient(pn patient.Patient) dbPatient {
	prdDB := dbPatient{
		ID:          pn.ID,
		UserID:      pn.UserID,
		Name:        pn.Name,
		Age:         pn.Age,
		Condition:   pn.Condition,
		Healed:      pn.Healed,
		VideoLinks:  pn.VideoLinks,
		DateCreated: pn.DateCreated.UTC(),
		DateUpdated: pn.DateUpdated.UTC(),
	}

	return prdDB
}

func toCorePatient(dbPn dbPatient) patient.Patient {
	prd := patient.Patient{
		ID:          dbPn.ID,
		UserID:      dbPn.UserID,
		Name:        dbPn.Name,
		Age:         dbPn.Age,
		Condition:   dbPn.Condition,
		Healed:      dbPn.Healed,
		VideoLinks:  dbPn.VideoLinks,
		DateCreated: dbPn.DateCreated.In(time.Local),
		DateUpdated: dbPn.DateUpdated.In(time.Local),
	}

	return prd
}

func toCorePatients(dbPns []dbPatient) []patient.Patient {
	pns := make([]patient.Patient, len(dbPns))

	for i, dbPrd := range dbPns {
		pns[i] = toCorePatient(dbPrd)
	}

	return pns
}
