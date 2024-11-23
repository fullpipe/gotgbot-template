package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.56

import (
	"bm/api/auth"
	"bm/api/gen"
	"bm/entity"
	"context"
)

// Me is the resolver for the me field.
func (r *queryResolver) Me(ctx context.Context) (*entity.User, error) {
	return auth.User(ctx), nil
}

// Query returns gen.QueryResolver implementation.
func (r *Resolver) Query() gen.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
