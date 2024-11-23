package directive

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/go-playground/validator/v10"
)

var (
	validate = validator.New()
)

func Validate(ctx context.Context, obj interface{}, next graphql.Resolver, constraint string) (interface{}, error) {
	val, err := next(ctx)
	if err != nil {
		return nil, err
	}

	err = validate.VarCtx(ctx, val, constraint)
	if err != nil {
		return nil, err
	}

	return val, nil
}
