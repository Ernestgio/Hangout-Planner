package mappings_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/Ernestgio/Hangout-Planner/services/Hangout/internal/dto"
	"github.com/Ernestgio/Hangout-Planner/services/Hangout/internal/mappings"
	"github.com/Ernestgio/Hangout-Planner/services/Hangout/internal/models"
)

func TestUserCreateRequestToModel(t *testing.T) {
	req := dto.UserCreateRequest{
		Name:  "Alice",
		Email: "alice@example.com",
	}

	user := mappings.UserCreateRequestToModel(req)

	assert.Equal(t, req.Name, user.Name)
	assert.Equal(t, req.Email, user.Email)
}

func TestUserToResponseDTO(t *testing.T) {
	user := &models.User{
		ID:    uuid.New(),
		Name:  "Bob",
		Email: "bob@example.com",
	}

	resp := mappings.UserToResponseDTO(user)

	assert.Equal(t, user.ID, resp.ID)
	assert.Equal(t, user.Name, resp.Name)
	assert.Equal(t, user.Email, resp.Email)
}
