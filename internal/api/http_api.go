package api

import "net/http"

type HttpApi interface {
	Routes() http.ServeMux
}
