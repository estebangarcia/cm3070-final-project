package api

import (
	"context"
	"net/http"
	"regexp"

	"github.com/go-chi/chi/v5"
)

type customRoute struct {
	Pattern *regexp.Regexp
	Handler http.HandlerFunc
}

type customMux struct {
	routes map[string][]customRoute
}

func NewCustomMux() customMux {
	return customMux{
		routes: map[string][]customRoute{},
	}
}

func (c customMux) Handle(w http.ResponseWriter, r *http.Request) {
	uri := chi.URLParam(r, "*")

	methodRoutes, exists := c.routes[r.Method]
	if !exists {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	for _, route := range methodRoutes {
		matches := route.Pattern.FindStringSubmatch(uri)
		if len(matches) > 0 {
			ctx := r.Context()
			for i, name := range route.Pattern.SubexpNames() {
				if i != 0 && name != "" {
					ctx = context.WithValue(ctx, name, matches[i])
				}
			}
			route.Handler(w, r.WithContext(ctx))
			return
		}
	}

	http.Error(w, "Not Found", http.StatusNotFound)
}

func (c customMux) addRoute(method string, route string, handler http.HandlerFunc) {
	r := customRoute{
		Pattern: regexp.MustCompile(route),
		Handler: handler,
	}
	methodRoutes, exists := c.routes[method]
	if !exists {
		methodRoutes = []customRoute{}
	}
	methodRoutes = append(methodRoutes, r)
	c.routes[method] = methodRoutes
}

func (c customMux) Post(route string, handler http.HandlerFunc) {
	c.addRoute(http.MethodPost, route, handler)
}

func (c customMux) Patch(route string, handler http.HandlerFunc) {
	c.addRoute(http.MethodPatch, route, handler)
}

func (c customMux) Delete(route string, handler http.HandlerFunc) {
	c.addRoute(http.MethodDelete, route, handler)
}

func (c customMux) Put(route string, handler http.HandlerFunc) {
	c.addRoute(http.MethodPut, route, handler)
}

func (c customMux) Get(route string, handler http.HandlerFunc) {
	c.addRoute(http.MethodGet, route, handler)
}

func (c customMux) Head(route string, handler http.HandlerFunc) {
	c.addRoute(http.MethodHead, route, handler)
}
