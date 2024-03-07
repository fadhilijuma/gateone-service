package region

import "github.com/fadhilijuma/gateone-service/business/web/v1/order"

// DefaultOrderBy represents the default way we sort.
var DefaultOrderBy = order.NewBy(OrderByID, order.ASC)

// Set of fields that the results can be ordered by.
const (
	OrderByID     = "region_id"
	OrderByName   = "name"
	OrderByUserID = "user_id"
)
