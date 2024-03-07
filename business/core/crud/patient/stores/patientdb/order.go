package patientdb

import (
	"fmt"
	"github.com/fadhilijuma/gateone-service/business/core/crud/patient"
	"github.com/fadhilijuma/gateone-service/business/web/v1/order"
)

var orderByFields = map[string]string{
	patient.OrderByPatientID: "patient_id",
	patient.OrderByUserID:    "user_id",
	patient.OrderByName:      "name",
	patient.OrderByAge:       "age",
	patient.OrderByCondition: "condition",
	patient.OrderByHealed:    "healed",
}

func orderByClause(orderBy order.By) (string, error) {
	by, exists := orderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}
