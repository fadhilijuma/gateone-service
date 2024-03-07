package condition

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"math/rand"
)

// TestGenerateNewConditions is a helper method for testing.
func TestGenerateNewConditions(n int, userID uuid.UUID) []NewCondition {
	newConditions := make([]NewCondition, n)

	idx := rand.Intn(10000)
	for i := 0; i < n; i++ {
		idx++
		nr := NewCondition{
			Name:   fmt.Sprintf("Name%d", idx),
			UserID: userID,
		}

		newConditions[i] = nr
	}

	return newConditions
}

// TestGenerateSeedConditions is a helper method for testing.
func TestGenerateSeedConditions(n int, api *Core, userID uuid.UUID) ([]Condition, error) {
	newConditions := TestGenerateNewConditions(n, userID)

	conditions := make([]Condition, len(newConditions))
	for i, nh := range newConditions {
		hme, err := api.Create(context.Background(), nh)
		if err != nil {
			return nil, fmt.Errorf("seeding condition: idx: %d : %w", i, err)
		}

		conditions[i] = hme
	}

	return conditions, nil
}
