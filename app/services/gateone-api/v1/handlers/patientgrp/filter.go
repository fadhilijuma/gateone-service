package patientgrp

import (
	"github.com/fadhilijuma/gateone-service/business/core/crud/patient"
	"github.com/fadhilijuma/gateone-service/foundation/validate"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

func parseFilter(r *http.Request) (patient.QueryFilter, error) {
	const (
		filterByPatientID = "patient_id"
		filterByUserID    = "user_id"
		filterByAge       = "age"
		filterByName      = "name"
		filterByCondition = "condition"
		filterByHealed    = "healed"
	)

	values := r.URL.Query()

	var filter patient.QueryFilter

	if productID := values.Get(filterByPatientID); productID != "" {
		id, err := uuid.Parse(productID)
		if err != nil {
			return patient.QueryFilter{}, validate.NewFieldsError(filterByPatientID, err)
		}
		filter.WithPatientID(id)
	}

	if age := values.Get(filterByAge); age != "" {
		ag, err := strconv.ParseInt(age, 10, 64)
		if err != nil {
			return patient.QueryFilter{}, validate.NewFieldsError(filterByAge, err)
		}
		filter.WithAge(int(ag))
	}

	if userID := values.Get(filterByUserID); userID != "" {
		id, err := uuid.Parse(userID)
		if err != nil {
			return patient.QueryFilter{}, validate.NewFieldsError(filterByUserID, err)
		}
		filter.WithUserID(id)
	}

	if name := values.Get(filterByName); name != "" {
		filter.WithName(name)
	}
	if condition := values.Get(filterByCondition); condition != "" {
		filter.WithName(condition)
	}
	if healed := values.Get(filterByHealed); healed != "" {
		hl, err := strconv.ParseBool(healed)
		if err != nil {
			return patient.QueryFilter{}, validate.NewFieldsError(filterByHealed, err)
		}
		filter.WithHealed(hl)
	}

	return filter, nil
}
