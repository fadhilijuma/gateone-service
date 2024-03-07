package role

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"math/rand"
)

// TestGenerateNewRoles is a helper method for testing.
func TestGenerateNewRoles(n int, userID uuid.UUID) []NewRole {
	newRoles := make([]NewRole, n)

	idx := rand.Intn(10000)
	for i := 0; i < n; i++ {
		idx++
		nr := NewRole{
			Name:   fmt.Sprintf("Name%d", idx),
			UserID: userID,
		}

		newRoles[i] = nr
	}

	return newRoles
}

// TestGenerateSeedRoles is a helper method for testing.
func TestGenerateSeedRoles(n int, api *Core, userID uuid.UUID) ([]Role, error) {
	newRoles := TestGenerateNewRoles(n, userID)

	roles := make([]Role, len(newRoles))
	for i, nh := range newRoles {
		hme, err := api.Create(context.Background(), nh)
		if err != nil {
			return nil, fmt.Errorf("seeding role: idx: %d : %w", i, err)
		}

		roles[i] = hme
	}

	return roles, nil
}
