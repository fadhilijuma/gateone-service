package condition

import "github.com/fadhilijuma/gateone-service/business/web/v1/order"

// DefaultOrderBy represents the default way we sort.
var DefaultOrderBy = order.NewBy(OrderByID, order.ASC)

// Set of fields that the results can be ordered by.
const (
	OrderByID     = "condition_id"
	OrderByName   = "name"
	OrderByUserID = "user_id"
)
