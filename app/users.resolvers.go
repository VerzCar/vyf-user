package app

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/vektah/gqlparser/v2/gqlerror"
	"gitlab.vecomentman.com/vote-your-face/service/user/api/model"
)

func (r *queryResolver) User(ctx context.Context) (*model.User, error) {
	user, err := r.userService.User(ctx)

	if err != nil {
		return nil, gqlerror.Errorf("cannot find user")
	}

	return user, nil
}
