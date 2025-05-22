package umi

import (
	"net/http"

	"github.com/victorguidi/umi/umi/middleware"
	"github.com/victorguidi/umi/umi/types"
)

type (
	Context = types.Context

	// Umi struct
	Umi struct {
		*http.ServeMux
		Context
		UmiOptions
		routes      []route
		middlewares []middleware.Middleware
	}

	// Options for Umi
	UmiOptions struct {
		PrintRoutes bool
		LogRequests bool
		Cors        bool
	}

	// Route holds the info of the registered routes
	route struct {
		method  string
		path    string
		handler types.HandlerFunc
	}

	// Configs for cors
	Cors struct {
		ORIGIN            string
		METHODS           string
		ALLOW_HEADERS     string
		ALLOW_CREDENTIALS string
	}
)
