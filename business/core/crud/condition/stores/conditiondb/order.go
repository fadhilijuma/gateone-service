package conditiondb

import (
	"fmt"
	"github.com/fadhilijuma/gateone-service/business/core/crud/condition"
	"github.com/fadhilijuma/gateone-service/business/web/v1/order"
)

var orderByFields = map[string]string{
	condition.OrderByID:     "condition_id",
	condition.OrderByName:   "name",
	condition.OrderByUserID: "user_id",
}

func orderByClause(orderBy order.By) (string, error) {
	by, exists := orderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}
