package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/nzwice/surl/pkg/validateutils"
)

func validateReqMw() endpoint.Middleware {
	return func(e endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if err := validateutils.Struct(ctx, request); err != nil {
				return nil, err
			}
			return e(ctx, request)
		}
	}
}
