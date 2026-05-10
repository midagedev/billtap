package stripecompat

import (
	"fmt"
	"strings"
)

type Route struct {
	Method string
	Path   string
	Source string
}

type RouteCatalog struct {
	routes  map[string]Route
	ordered []Route
}

func NewRouteCatalog(routes []Route) (RouteCatalog, error) {
	catalog := RouteCatalog{routes: map[string]Route{}}
	for _, route := range routes {
		route.Method = strings.ToUpper(strings.TrimSpace(route.Method))
		route.Path = NormalizePath(route.Path)
		route.Source = strings.TrimSpace(route.Source)
		if route.Method == "" {
			return RouteCatalog{}, fmt.Errorf("route for %s has empty method", route.Path)
		}
		if route.Path == "" {
			return RouteCatalog{}, fmt.Errorf("route for %s has empty path", route.Method)
		}
		key := RouteKey(route.Method, route.Path)
		if _, exists := catalog.routes[key]; exists {
			return RouteCatalog{}, fmt.Errorf("duplicate known Stripe route: %s", key)
		}
		catalog.routes[key] = route
		catalog.ordered = append(catalog.ordered, route)
	}
	return catalog, nil
}

func MustRouteCatalog(routes []Route) RouteCatalog {
	catalog, err := NewRouteCatalog(routes)
	if err != nil {
		panic(err)
	}
	return catalog
}

func DefaultRouteCatalog() RouteCatalog {
	return MustRouteCatalog(DefaultKnownRoutes())
}

func (c RouteCatalog) Lookup(method string, path string) (Route, bool) {
	route, ok := c.routes[RouteKey(method, NormalizePath(path))]
	if ok {
		return route, true
	}
	normalizedMethod := strings.ToUpper(strings.TrimSpace(method))
	for _, route := range c.ordered {
		if route.Method == normalizedMethod && pathMatches(route.Path, path) {
			return route, true
		}
	}
	return Route{}, false
}

func (c RouteCatalog) Routes() []Route {
	routes := make([]Route, len(c.ordered))
	copy(routes, c.ordered)
	return routes
}
