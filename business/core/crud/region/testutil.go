package region

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"math/rand"
)

// TestGenerateNewRegions is a helper method for testing.
func TestGenerateNewRegions(n int, userID uuid.UUID) []NewRegion {
	newRegions := make([]NewRegion, n)

	idx := rand.Intn(10000)
	for i := 0; i < n; i++ {
		idx++
		nr := NewRegion{
			Name:   fmt.Sprintf("Name%d", idx),
			UserID: userID,
		}

		newRegions[i] = nr
	}

	return newRegions
}

// TestGenerateSeedRegions is a helper method for testing.
func TestGenerateSeedRegions(n int, api *Core, userID uuid.UUID) ([]Region, error) {
	newRegions := TestGenerateNewRegions(n, userID)

	Regions := make([]Region, len(newRegions))
	for i, nh := range newRegions {
		hme, err := api.Create(context.Background(), nh)
		if err != nil {
			return nil, fmt.Errorf("seeding Region: idx: %d : %w", i, err)
		}

		Regions[i] = hme
	}

	return Regions, nil
}
