package repos

import (
	"github.com/lemmego/gpa"
	"github.com/lemmego/gpaorm"
)

// SQLRepo returns a SQL repository for the given entity type.
// The repository implements the full MigratableRepository interface with SQL-specific operations.
func SQLRepo[T any](instanceName ...string) gpa.MigratableRepository[T] {
	return gpaorm.GetRepositoryByName[T](instanceName...)
}
