package repos

import (
	"github.com/lemmego/gpa"
	"github.com/lemmego/lemmego/internal/models"
)

// UserRepository provides user-specific database operations
type UserRepository struct {
	gpa.MigratableRepository[models.User]
}

// User returns a UserRepository instance
func User(instanceName ...string) *UserRepository {
	repo := SQLRepo[models.User](instanceName...)
	return &UserRepository{repo}
}

// Custom methods specific to User operations
// Add your custom repository methods here
