package rolegrp

import (
	"errors"
	"github.com/fadhilijuma/gateone-service/business/core/crud/role"
	"github.com/fadhilijuma/gateone-service/business/web/v1/order"
	"github.com/fadhilijuma/gateone-service/foundation/validate"
	"net/http"
)

func parseOrder(r *http.Request) (order.By, error) {
	const (
		orderByRoleID = "role_id"
		orderByUserID = "user_id"
		orderByName   = "name"
	)

	var orderByFields = map[string]string{
		orderByRoleID: role.OrderByID,
		orderByName:   role.OrderByName,
		orderByUserID: role.OrderByUserID,
	}

	orderBy, err := order.Parse(r, order.NewBy(orderByRoleID, order.ASC))
	if err != nil {
		return order.By{}, err
	}

	if _, exists := orderByFields[orderBy.Field]; !exists {
		return order.By{}, validate.NewFieldsError(orderBy.Field, errors.New("order field does not exist"))
	}

	orderBy.Field = orderByFields[orderBy.Field]

	return orderBy, nil
}
