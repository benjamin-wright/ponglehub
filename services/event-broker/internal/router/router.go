package router

import (
	"fmt"
	"strings"
)

type Route struct {
	Filter string
	Parts  []string
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
		Filter: filter,
		Parts:  strings.Split(filter, "."),
		URL:    url,
	})
}

func (r *Router) Remove(filter string, url string) error {
	for idx, route := range r.routes {
		if route.Filter == filter && route.URL == url {
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
		lenFilterParts := len(route.Parts)

		matchIndex := 0
		doubleWild := false
		mismatch := false

		for _, part := range typeParts {
			if matchIndex < lenFilterParts {
				if route.Parts[matchIndex] == "*" {
					matchIndex += 1
					doubleWild = false
					continue
				}

				if route.Parts[matchIndex] == "**" {
					matchIndex += 1
					doubleWild = true
					continue
				}

				if route.Parts[matchIndex] == part {
					matchIndex += 1
					doubleWild = false
					continue
				}
			}

			if doubleWild {
				continue
			}

			mismatch = true
			break
		}

		if mismatch || matchIndex != lenFilterParts {
			continue
		}

		urls = append(urls, route.URL)
	}

	return urls
}
