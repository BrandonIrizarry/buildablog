package main

import (
	"net/http"

	"github.com/BrandonIrizarry/buildablog/internal/environment"
)

type config struct {
	environment.Env
	handler http.HandlerFunc
}
