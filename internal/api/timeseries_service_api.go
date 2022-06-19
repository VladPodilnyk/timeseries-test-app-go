package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/golang/gddo/httputil/header"
	"google.golang.org/protobuf/proto"

	"github.com/VladPodilnyk/timeseries-test-app-go/internal/config"
	"github.com/VladPodilnyk/timeseries-test-app-go/internal/model"
	"github.com/VladPodilnyk/timeseries-test-app-go/internal/repo"
)

type TimeSeriesServiceApi struct {
	service GrpcService
}

func New(repo repo.TimeSeries, config config.LimitsConfig) *TimeSeriesServiceApi {
	service := GrpcServiceIml{repo, config}
	return &TimeSeriesServiceApi{&service}
}

func (api *TimeSeriesServiceApi) Routes() *http.ServeMux {
	mux := http.NewServeMux()

	staticPage := http.FileServer(http.Dir("./static"))
	fetchData := dataHandler(api)

	mux.Handle("/timeseries", staticPage)
	mux.HandleFunc("/fetch", fetchData)

	return mux
}

func dataHandler(api *TimeSeriesServiceApi) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			http.Error(w, "Method not allowed", 405)
			return
		}

		decodedData, err := decodeJSONBody[model.UserRequest](w, request)
		if err != nil {
			var malformedRequest *model.MalformedRequest
			if errors.As(err, &malformedRequest) {
				http.Error(w, malformedRequest.Error(), malformedRequest.Status)
			} else {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}

		// TODO fetch data
		result, err := api.service.FetchData(*decodedData)
		if err != nil {
			// TODO: better error messages
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		encoded, err := proto.Marshal(result)
		w.Write(encoded)
	}
}

func decodeJSONBody[T any](w http.ResponseWriter, request *http.Request) (*T, error) {
	if request.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(request.Header, "Content-Type")
		if value != "application/json" {
			msg := "Content-Type header is not application/json"
			return nil, &model.MalformedRequest{http.StatusUnsupportedMediaType, model.DomainError{msg}}
		}
	}

	// TODO: remove magic numbers
	request.Body = http.MaxBytesReader(w, request.Body, 1048576)
	dec := json.NewDecoder(request.Body)
	dec.DisallowUnknownFields()

	var decodedValue *T
	err := dec.Decode(decodedValue)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshallTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-fromed JSON (at position %d)", syntaxError.Offset)
			return nil, &model.MalformedRequest{Status: http.StatusBadRequest, Diagnostics: model.DomainError{Message: msg}}

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Request body contains badly-fromed JSON")
			return nil, &model.MalformedRequest{Status: http.StatusBadRequest, Diagnostics: model.DomainError{Message: msg}}

		case errors.As(err, &unmarshallTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshallTypeError.Field, unmarshallTypeError.Offset)
			return nil, &model.MalformedRequest{Status: http.StatusBadRequest, Diagnostics: model.DomainError{Message: msg}}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return nil, &model.MalformedRequest{Status: http.StatusBadRequest, Diagnostics: model.DomainError{Message: msg}}

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return nil, &model.MalformedRequest{Status: http.StatusBadRequest, Diagnostics: model.DomainError{Message: msg}}

		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			return nil, &model.MalformedRequest{Status: http.StatusBadRequest, Diagnostics: model.DomainError{Message: msg}}

		default:
			return nil, err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		msg := "Request body must contain a single JSON object"
		return nil, &model.MalformedRequest{Status: http.StatusBadRequest, Diagnostics: model.DomainError{Message: msg}}
	}

	return decodedValue, nil
}
