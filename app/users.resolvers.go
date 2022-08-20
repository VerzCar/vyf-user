package app

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/vektah/gqlparser/v2/gqlerror"
	"gitlab.vecomentman.com/vote-your-face/service/user/api/model"
)

// UpdateUser is the resolver for the updateUser field.
func (r *mutationResolver) UpdateUser(ctx context.Context, user model.UserUpdateInput) (*model.User, error) {
	gqlError := gqlerror.Errorf("user cannot be updated")

	if err := r.validate.Struct(user); err != nil {
		r.log.Error(err)
		return nil, gqlError
	}

	updatedUser, err := r.userService.UpdateUser(ctx, &user)

	if err != nil {
		return nil, gqlError
	}

	return updatedUser, nil
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, identityID *string) (*model.User, error) {
	user, err := r.userService.User(ctx, identityID)

	if err != nil {
		return nil, gqlerror.Errorf("cannot find user")
	}

	return user, nil
}
