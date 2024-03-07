package patient

import "github.com/fadhilijuma/gateone-service/business/web/v1/order"

// DefaultOrderBy represents the default way we sort.
var DefaultOrderBy = order.NewBy(OrderByPatientID, order.ASC)

// Set of fields that the results can be ordered by.
const (
	OrderByPatientID = "patient_id"
	OrderByUserID    = "user_id"
	OrderByName      = "name"
	OrderByAge       = "age"
	OrderByCondition = "condition"
	OrderByHealed    = "healed"
)
