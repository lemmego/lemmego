package repos

import (
	"github.com/lemmego/gpa"
	"github.com/lemmego/gpagorm"
	"github.com/lemmego/lemmego/internal/models"
)

// UserRepository provides user-specific database operations
type UserRepository struct {
	gpa.Repository[models.User]
}

// User returns a UserRepository instance
func User(instanceName ...string) *UserRepository {
	repo := gpagorm.GetRepository[models.User](instanceName...)
	return &UserRepository{Repository: repo}
}

// Custom methods specific to User operations
// Add your custom repository methods here
