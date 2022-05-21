package app

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/vektah/gqlparser/v2/gqlerror"
	"gitlab.vecomentman.com/service/user/api/model"
)

func (r *mutationResolver) CreateCompany(ctx context.Context, company model.CompanyCreateInput) (*model.Company, error) {
	//currentUser, err := r.authService.UserFromContext(ctx)
	//
	//if err != nil {
	//	return nil, gqlerror.Errorf("authentication failed")
	//}

	gqlError := gqlerror.Errorf("company cannot be created")

	if err := r.validate.Struct(company); err != nil {
		r.log.Error(err)
		return nil, gqlError
	}

	// TODO change user here
	newCompany, err := r.companyService.CreateCompany(ctx, &company, &model.User{})

	if err != nil {
		return nil, gqlError
	}

	return newCompany, nil
}

func (r *mutationResolver) VerifyCompany(ctx context.Context, payload model.CompanyVerifyInput) (*model.Company, error) {
	gqlError := gqlerror.Errorf("company cannot be verified")

	company, err := r.companyService.VerifyCompany(ctx, &payload)

	if err != nil {
		return nil, gqlError
	}

	return company, nil
}
