package mapper_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/domain"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/dto"
	mapper "github.com/Ernestgio/Hangout-Planner/services/hangout/internal/mapper"
)

func TestCreateUserRequestToModel(t *testing.T) {
	req := dto.CreateUserRequest{
		Name:  "Alice",
		Email: "alice@example.com",
	}

	user := mapper.CreateUserRequestToModel(req)

	assert.Equal(t, req.Name, user.Name)
	assert.Equal(t, req.Email, user.Email)
}

func TestUserToResponseDTO(t *testing.T) {
	user := &domain.User{
		ID:    uuid.New(),
		Name:  "Bob",
		Email: "bob@example.com",
	}

	resp := mapper.UserToResponseDTO(user)

	assert.Equal(t, user.ID, resp.ID)
	assert.Equal(t, user.Name, resp.Name)
	assert.Equal(t, user.Email, resp.Email)
}
