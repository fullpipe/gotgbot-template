package directive

import (
	"bm/api/auth"
	"bm/entity"
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/samber/lo"
)

func HasRole(ctx context.Context, _ interface{}, next graphql.Resolver, role entity.Role) (interface{}, error) {
	if !lo.Contains(auth.Roles(ctx), string(role)) {
		return nil, fmt.Errorf("Access denied")
	}

	return next(ctx)
}
