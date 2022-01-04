package router

import (
	"fmt"
	"strings"
)

type Route struct {
	Filter []string
	URL    string
}

type Router struct {
	routes []Route
}

func New() Router {
	return Router{
		routes: []Route{},
	}
}

func (r *Router) Add(filter string, url string) {
	r.routes = append(r.routes, Route{
		Filter: strings.Split(filter, "."),
		URL:    url,
	})
}

func (r *Router) Remove(filter string, url string) error {
	for idx, route := range r.routes {
		if strings.Join(route.Filter, ".") == filter && route.URL == url {
			r.routes[idx] = r.routes[len(r.routes)-1]
			r.routes = r.routes[:len(r.routes)-1]
			return nil
		}
	}

	return fmt.Errorf("failed to remove %s: %s, not found", filter, url)
}

func (r *Router) GetURLs(eventType string) []string {
	urls := []string{}

	typeParts := strings.Split(eventType, ".")

	for _, route := range r.routes {
		if len(route.Filter) != len(typeParts) {
			continue
		}

		mismatch := false
		for idx, _ := range route.Filter {
			if route.Filter[idx] != "*" && route.Filter[idx] != typeParts[idx] {
				mismatch = true
				break
			}
		}

		if mismatch {
			continue
		}

		urls = append(urls, route.URL)
	}

	return urls
}
