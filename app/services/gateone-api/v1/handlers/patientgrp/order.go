package patientgrp

import (
	"errors"
	"github.com/fadhilijuma/gateone-service/business/core/crud/patient"
	"github.com/fadhilijuma/gateone-service/business/web/v1/order"
	"github.com/fadhilijuma/gateone-service/foundation/validate"
	"net/http"
)

func parseOrder(r *http.Request) (order.By, error) {
	const (
		orderByPatientID = "patient_id"
		orderByUserID    = "user_id"
		orderByName      = "name"
		orderByAge       = "age"
		orderByCondition = "condition"
		orderByHealed    = "healed"
	)

	var orderByFields = map[string]string{
		orderByPatientID: patient.OrderByPatientID,
		orderByName:      patient.OrderByName,
		orderByUserID:    patient.OrderByUserID,
		orderByAge:       patient.OrderByAge,
		orderByCondition: patient.OrderByCondition,
		orderByHealed:    patient.OrderByHealed,
	}

	orderBy, err := order.Parse(r, order.NewBy(orderByPatientID, order.ASC))
	if err != nil {
		return order.By{}, err
	}

	if _, exists := orderByFields[orderBy.Field]; !exists {
		return order.By{}, validate.NewFieldsError(orderBy.Field, errors.New("order field does not exist"))
	}

	orderBy.Field = orderByFields[orderBy.Field]

	return orderBy, nil
}
