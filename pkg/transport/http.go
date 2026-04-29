package transport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-playground/validator/v10"

	"github.com/nzwice/surl/pkg/endpoints"
)

func HttpHandler(e endpoints.Set) http.Handler {

	mux := http.NewServeMux()

	var handlerOpts = []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(errorEncoder),
	}

	shortenUrlHandler := httptransport.NewServer(
		e.ShortenUrl,
		jsonDecodeRequest[endpoints.ShortenUrlRequest](),
		jsonEncodeResponse[endpoints.ShortenUrlResponse](),
		handlerOpts...,
	)

	redirectUrlHandler := httptransport.NewServer(
		e.GetOriginalUrl,
		decodeRedirectUrlRequest(),
		encodeRedirectUrlResponse(),
		handlerOpts...,
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
		w.Header().Set("Content-Type", "application/json")
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

func errorEncoder(ctx context.Context, err error, w http.ResponseWriter) {
	status, resp := extractError(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}

func extractValidationError(vErrors validator.ValidationErrors) (int, errorResp) {
	var data []string
	for _, f := range vErrors {
		// TODO: need to custom error message. Preferraly having a tag in struct fields for validation message.
		data = append(data, f.Error())
	}
	return http.StatusBadRequest, errorResp{
		Error: "bad request",
		Data:  data,
	}
}

func extractError(err error) (int, errorResp) {
	var vErrors validator.ValidationErrors
	if errors.As(err, &vErrors) {
		return extractValidationError(vErrors)
	}
	return http.StatusInternalServerError, errorResp{
		Error: err.Error(),
	}
}

type errorResp struct {
	Error string `json:"error"`
	Data  any    `json:"data"`
}
