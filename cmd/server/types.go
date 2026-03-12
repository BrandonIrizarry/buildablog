package main

import (
	"net/http"

	"github.com/BrandonIrizarry/buildablog/v2/internal/environment"
)

type config struct {
	environment.Env
	handler http.HandlerFunc
}
