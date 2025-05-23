package umi

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/victorguidi/umi/middleware"
	"github.com/victorguidi/umi/types"
)

func (u *Umi) registerRoute(method, path string, handler types.HandlerFunc) {
	u.routes = append(u.routes, route{
		method:  method,
		path:    path,
		handler: handler,
	})
}

// This is the handler that allows all others handlers to simply return an error
func (u *Umi) defaultHandler(handler types.HandlerFunc) http.HandlerFunc {
	// Copy of middlewares. This will make routers declared before the USE function to not have the middleware
	middlewares := make([]middleware.Middleware, len(u.middlewares))
	copy(middlewares, u.middlewares)
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		baseHandler := func(w http.ResponseWriter, r *http.Request) {
			u.Context = Context{ResponseWriter: w, Request: r}
			err := handler(&u.Context)
			if err != nil {
				if error, ok := err.(types.Error); ok {
					w.WriteHeader(error.Status)
					fmt.Fprint(w, err.Error())
				} else {
					w.WriteHeader(500)
					fmt.Fprint(w, err)
				}
			}
		}

		if len(middlewares) > 0 {
			chainedMiddlewares := middleware.Chain(middlewares...)(baseHandler)
			chainedMiddlewares(w, r)
		} else {
			baseHandler(w, r)
		}
	}
}

func (u *Umi) rebuildRoutes() {
	u.ServeMux = http.NewServeMux()

	for _, route := range u.routes {
		u.HandleFunc(route.method+" "+route.path, u.defaultHandler(route.handler))
	}
}

func defaultOptions() UmiOptions {
	return UmiOptions{
		PrintRoutes:   true,
		LogEvents:     true,
		Cors:          false,
		ServerOptions: nil,
	}
}

func (u *Umi) withLogger() *Umi {
	u.middlewares = append(u.middlewares, middleware.Logger())
	return u
}

// Creates a new instance of Umi
func New() *Umi {
	mux := http.NewServeMux()
	umi := &Umi{
		ServeMux:   mux,
		Context:    Context{},
		UmiOptions: defaultOptions(),
		routes:     make([]route, 0),
	}

	if umi.LogEvents {
		umi.withLogger()
	}

	if umi.Cors {
		umi.WithFlexibleCors()
	}

	return umi
}

// Allows the user to modify CORS rules,
// by default Umi follows a zero trust approach,
// this means that it will block any request from
// different origins
// This function expects the user to specify each part of its Cors rules
func (u *Umi) WithCors(rules Cors) *Umi {
	u.middlewares = append(u.middlewares, middleware.Cors(rules.ORIGIN, rules.ALLOW_CREDENTIALS, rules.ALLOW_HEADERS, rules.METHODS))
	return u
}

// Defaults to Origin *
func (u *Umi) WithFlexibleCors() *Umi {
	u.middlewares = append(u.middlewares, middleware.FlexibleCors())
	return u
}

func (u *Umi) WithOptions(opts UmiOptions) *Umi {
	u.UmiOptions = opts
	return u
}

func (u *Umi) WithServerOptions(opts *http.Server) *Umi {
	u.ServerOptions = opts
	return u
}

func (u *Umi) Start(addr string) {
	var server *http.Server
	if u.ServerOptions != nil {
		server = u.ServerOptions
		server.Addr = addr
		server.Handler = u
	} else {
		server = &http.Server{
			Handler:        u,
			Addr:           addr,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
	}

	if u.PrintRoutes {
		log.Println("Registered routes")
		for _, r := range u.routes {
			log.Printf("%s at %s", r.method, r.path)
		}
	}

	log.Printf("Umi is listening on addr %s", addr)
	log.Fatal(server.ListenAndServe())
}

func (u *Umi) Use(middleware middleware.Middleware) {
	u.middlewares = append(u.middlewares, middleware)
}

// Methods follow the specificaton at:
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Methods

// The GET method requests a representation of the specified resource.
// Requests using GET should only retrieve data and should not contain
// a request content.
func (u *Umi) GET(path string, handler types.HandlerFunc) {
	u.registerRoute("GET", path, handler)
	u.HandleFunc("GET "+path, u.defaultHandler(handler))
	u.HandleFunc("OPTIONS "+path, u.defaultHandler(handler))
}

// The POST method submits an entity to the specified resource,
// often causing a change in state or side effects on the server.
func (u *Umi) POST(path string, handler types.HandlerFunc) {
	u.registerRoute("POST", path, handler)
	u.HandleFunc("POST "+path, u.defaultHandler(handler))
	u.HandleFunc("OPTIONS "+path, u.defaultHandler(handler))
}

// The PUT method replaces all current representations of the target
// resource with the request content.
func (u *Umi) PUT(path string, handler types.HandlerFunc) {
	u.registerRoute("PUT", path, handler)
	u.HandleFunc("PUT "+path, u.defaultHandler(handler))
	u.HandleFunc("OPTIONS "+path, u.defaultHandler(handler))
}

// The DELETE method deletes the specified resource.
func (u *Umi) DELETE(path string, handler types.HandlerFunc) {
	u.registerRoute("DELETE", path, handler)
	u.HandleFunc("DELETE "+path, u.defaultHandler(handler))
	u.HandleFunc("OPTIONS "+path, u.defaultHandler(handler))
}

// The PATCH method applies partial modifications to a resource.
func (u *Umi) PATCH(path string, handler types.HandlerFunc) {
	u.registerRoute("PATCH", path, handler)
	u.HandleFunc("PATCH "+path, u.defaultHandler(handler))
	u.HandleFunc("OPTIONS "+path, u.defaultHandler(handler))
}

// The HEAD method asks for a response identical to a GET request,
// but without a response body.
func (u *Umi) HEAD(path string, handler types.HandlerFunc) {
	u.registerRoute("HEAD", path, handler)
	u.HandleFunc("", u.defaultHandler(handler))
}

// The CONNECT method establishes a tunnel to the server identified by
// the target resource.
func (u *Umi) CONNECT(path string, handler types.HandlerFunc) {
	u.registerRoute("CONNECT", path, handler)
	u.HandleFunc("", u.defaultHandler(handler))
}

// The OPTIONS method describes the communication options for the
// target resource.
func (u *Umi) OPTIONS(path string, handler types.HandlerFunc) {
	u.registerRoute("OPTIONS", path, handler)
	u.HandleFunc("OPTIONS "+path, u.defaultHandler(handler))
}

// The TRACE method performs a message loop-back test along the path to the
// target resource.
func (u *Umi) TRACE(path string, handler types.HandlerFunc) {
	u.registerRoute("TRACE", path, handler)
	u.HandleFunc("TRACE "+path, u.defaultHandler(handler))
}
