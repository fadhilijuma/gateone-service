package patient

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/google/uuid"
)

// TestGenerateNewPatients is a helper method for testing.
func TestGenerateNewPatients(n int, userID uuid.UUID) []NewPatient {
	newPrds := make([]NewPatient, n)

	idx := rand.Intn(10000)
	for i := 0; i < n; i++ {
		idx++

		np := NewPatient{
			UserID:     userID,
			Name:       "test product",
			Age:        10,
			Condition:  "deaf",
			VideoLinks: []string{"https://www.youtube.com/watch?v=1234"},
			Healed:     false,
		}

		newPrds[i] = np
	}

	return newPrds
}

// TestGenerateSeedPatients is a helper method for testing.
func TestGenerateSeedPatients(n int, api *Core, userID uuid.UUID) ([]Patient, error) {
	newPrds := TestGenerateNewPatients(n, userID)

	prds := make([]Patient, len(newPrds))
	for i, np := range newPrds {
		prd, err := api.Create(context.Background(), np)
		if err != nil {
			return nil, fmt.Errorf("seeding patient: idx: %d : %w", i, err)
		}

		prds[i] = prd
	}

	return prds, nil
}
