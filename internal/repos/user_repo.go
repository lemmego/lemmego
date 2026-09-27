package repos

import (
	"context"
	"fmt"
	"strconv"

	"github.com/lemmego/auth"
	"github.com/lemmego/gpa"
	"github.com/lemmego/lemmego/internal/models"
)

type UserRepository struct {
	gpa.MigratableRepository[models.User]
}

func User(instanceName ...string) *UserRepository {
	repo := SQLRepo[models.User](instanceName...)
	return &UserRepository{repo}
}

// FindByID resolves the identifier a credential carries back to a row.
//
// It deliberately shadows the embedded repository's FindByID, which takes an
// untyped key. auth works in string ids — a JWT subject is a string — and this
// is the single place that knows the column behind it is a uint64. The
// embedded method is still reachable as r.MigratableRepository.FindByID for a caller that
// already holds a typed key.
//
// An id that does not parse is a user that does not exist, not a failure. A
// credential issued before the identifier format changed lands here; reporting
// it as an error would make auth answer 503, so on the day of an upgrade a
// wave of stale tokens would read as an outage instead of asking one person to
// sign in again.
func (r *UserRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	pk, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("%w: %q is not a user id", auth.ErrUserNotFound, id)
	}
	user, err := r.MigratableRepository.FindByID(ctx, pk)
	if gpa.IsNotFound(err) {
		return nil, fmt.Errorf("%w: id %d", auth.ErrUserNotFound, pk)
	}
	return user, err
}
