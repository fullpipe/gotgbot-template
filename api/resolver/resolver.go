package resolver

import "bm/repository"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

func NewResolver(
	userRepo *repository.UserRepo,
) *Resolver {
	return &Resolver{
		userRepo: userRepo,
	}
}

type Resolver struct {
	userRepo *repository.UserRepo
}
