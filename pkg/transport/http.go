package transport

import (
	"context"
	"encoding/json"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/nzwice/surl/pkg/endpoints"
)

func HttpHandler(e endpoints.Set) http.Handler {

	mux := http.NewServeMux()

	shortenUrlHandler := httptransport.NewServer(
		e.ShortenUrl,
		jsonDecodeRequest[endpoints.ShortenUrlRequest](),
		jsonEncodeResponse[endpoints.ShortenUrlResponse](),
	)

	redirectUrlHandler := httptransport.NewServer(
		e.GetOriginalUrl,
		decodeRedirectUrlRequest(),
		encodeRedirectUrlResponse(),
	)

	mux.Handle("POST /api/v1/surl", shortenUrlHandler)
	mux.Handle("GET /r/{shortCode}", redirectUrlHandler)

	return withRecovery(withRequestId(withLogging(mux)))
}

func jsonDecodeRequest[T any]() httptransport.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (request interface{}, err error) {
		v := new(T)
		err = json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			return nil, err
		}
		return *v, nil
	}
}

func jsonEncodeResponse[T any]() httptransport.EncodeResponseFunc {
	return func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
		response = response.(T)
		err := json.NewEncoder(w).Encode(response)
		return err
	}
}

func decodeRedirectUrlRequest() httptransport.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (request interface{}, err error) {
		shortCode := r.PathValue("shortCode")
		return endpoints.GetOriginalUrlRequest{
			ShortCode: shortCode,
		}, nil
	}
}

func encodeRedirectUrlResponse() httptransport.EncodeResponseFunc {
	return func(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
		parsed := resp.(endpoints.GetOriginalUrlResponse)
		w.Header().Set("Location", parsed.OriginalUrl)
		w.WriteHeader(http.StatusFound)
		return nil
	}
}
