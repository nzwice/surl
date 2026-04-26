package endpoints

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/nzwice/surl/pkg/shortensvc"
)

type Set struct {
	ShortenUrl     endpoint.Endpoint
	GetOriginalUrl endpoint.Endpoint
}

func MakeEndpoints(shortenService shortensvc.Service) Set {
	return Set{
		ShortenUrl: endpoint.Chain(
			loggingMw(),
		)(makeShortenUrlEndpoint(shortenService)),
		GetOriginalUrl: endpoint.Chain(
			loggingMw(),
		)(makeGetOrignalUrlEndpoint(shortenService)),
	}
}

type ShortenUrlRequest struct {
	OriginalUrl string     `json:"original_url"`
	Alias       *string    `json:"alias"`
	ExpiredAt   *time.Time `json:"expired_at"`
}

type ShortenUrlResponse struct {
	ShortCode string `json:"short_code"`
}

type GetOriginalUrlRequest struct {
	ShortCode string `json:"short_code"`
}

type GetOriginalUrlResponse struct {
	OriginalUrl string `json:"original_url"`
}

func makeShortenUrlEndpoint(svc shortensvc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(ShortenUrlRequest)
		resp, err := svc.ShortenUrl(ctx, req.OriginalUrl, req.Alias, req.ExpiredAt)
		if err != nil {
			return nil, err
		}
		return ShortenUrlResponse{
			ShortCode: resp,
		}, nil
	}
}

func makeGetOrignalUrlEndpoint(svc shortensvc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetOriginalUrlRequest)
		resp, err := svc.GetOriginalUrl(ctx, req.ShortCode)
		if err != nil {
			return nil, err
		}
		return GetOriginalUrlResponse{
			OriginalUrl: resp,
		}, nil
	}
}
