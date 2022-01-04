package router

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type RouteInput struct {
	Filter string
	URL    string
}

func TestRouterAdd(t *testing.T) {
	for _, test := range []struct {
		Name     string
		Input    RouteInput
		Existing []Route
		Expected []Route
	}{
		{
			Name:  "empty",
			Input: RouteInput{Filter: "new-filter", URL: "new-url"},
			Expected: []Route{
				{Filter: []string{"new-filter"}, URL: "new-url"},
			},
		},
		{
			Name:  "split",
			Input: RouteInput{Filter: "new.filter", URL: "new-url"},
			Expected: []Route{
				{Filter: []string{"new", "filter"}, URL: "new-url"},
			},
		},
		{
			Name:  "added",
			Input: RouteInput{Filter: "new-filter", URL: "new-url"},
			Existing: []Route{
				{Filter: []string{"old-filter"}, URL: "old-url"},
			},
			Expected: []Route{
				{Filter: []string{"old-filter"}, URL: "old-url"},
				{Filter: []string{"new-filter"}, URL: "new-url"},
			},
		},
		{
			Name:  "duplicate",
			Input: RouteInput{Filter: "new-filter", URL: "new-url"},
			Existing: []Route{
				{Filter: []string{"new-filter"}, URL: "new-url"},
			},
			Expected: []Route{
				{Filter: []string{"new-filter"}, URL: "new-url"},
				{Filter: []string{"new-filter"}, URL: "new-url"},
			},
		},
	} {
		t.Run(test.Name, func(u *testing.T) {
			router := Router{
				routes: test.Existing,
			}

			router.Add(test.Input.Filter, test.Input.URL)
			assert.Equal(u, test.Expected, router.routes)
		})
	}
}

func TestRouterRemove(t *testing.T) {
	for _, test := range []struct {
		Name     string
		Input    RouteInput
		Existing []Route
		Expected []Route
		Error    string
	}{
		{
			Name:     "empty",
			Input:    RouteInput{Filter: "new-filter", URL: "new-url"},
			Existing: []Route{},
			Error:    "failed to remove new-filter: new-url, not found",
		},
		{
			Name:  "not matching filter",
			Input: RouteInput{Filter: "wrong-filter", URL: "old-url"},
			Existing: []Route{
				{Filter: []string{"old-filter"}, URL: "old-url"},
				{Filter: []string{"new-filter"}, URL: "new-url"},
			},
			Expected: []Route{
				{Filter: []string{"old-filter"}, URL: "old-url"},
				{Filter: []string{"new-filter"}, URL: "new-url"},
			},
			Error: "failed to remove wrong-filter: old-url, not found",
		},
		{
			Name:  "not matching url",
			Input: RouteInput{Filter: "old-filter", URL: "wrong-url"},
			Existing: []Route{
				{Filter: []string{"old-filter"}, URL: "old-url"},
				{Filter: []string{"new-filter"}, URL: "new-url"},
			},
			Expected: []Route{
				{Filter: []string{"old-filter"}, URL: "old-url"},
				{Filter: []string{"new-filter"}, URL: "new-url"},
			},
			Error: "failed to remove old-filter: wrong-url, not found",
		},
		{
			Name:  "first",
			Input: RouteInput{Filter: "old-filter", URL: "old-url"},
			Existing: []Route{
				{Filter: []string{"old-filter"}, URL: "old-url"},
				{Filter: []string{"new-filter"}, URL: "new-url"},
			},
			Expected: []Route{
				{Filter: []string{"new-filter"}, URL: "new-url"},
			},
		},
		{
			Name:  "last",
			Input: RouteInput{Filter: "new-filter", URL: "new-url"},
			Existing: []Route{
				{Filter: []string{"old-filter"}, URL: "old-url"},
				{Filter: []string{"new-filter"}, URL: "new-url"},
			},
			Expected: []Route{
				{Filter: []string{"old-filter"}, URL: "old-url"},
			},
		},
		{
			Name:  "only",
			Input: RouteInput{Filter: "new-filter", URL: "new-url"},
			Existing: []Route{
				{Filter: []string{"new-filter"}, URL: "new-url"},
			},
			Expected: []Route{},
		},
		{
			Name:  "split",
			Input: RouteInput{Filter: "new.filter", URL: "new-url"},
			Existing: []Route{
				{Filter: []string{"new", "filter"}, URL: "new-url"},
			},
			Expected: []Route{},
		},
	} {
		t.Run(test.Name, func(u *testing.T) {
			router := Router{
				routes: test.Existing,
			}

			err := router.Remove(test.Input.Filter, test.Input.URL)

			if test.Expected != nil {
				assert.Equal(u, test.Expected, router.routes)
			}

			if test.Error != "" {
				assert.EqualError(u, err, test.Error)
			} else {
				assert.NoError(u, err)
			}
		})
	}
}

func TestRouterGetURLs(t *testing.T) {
	for _, test := range []struct {
		Name     string
		Input    string
		Routes   []Route
		Expected []string
	}{
		{
			Name:     "empty",
			Input:    "event-type",
			Routes:   []Route{},
			Expected: []string{},
		},
		{
			Name:  "match",
			Input: "event-type",
			Routes: []Route{
				{Filter: []string{"event-type"}, URL: "test-url"},
				{Filter: []string{"other-type"}, URL: "wrong-url"},
			},
			Expected: []string{"test-url"},
		},
		{
			Name:  "match multiple",
			Input: "event-type",
			Routes: []Route{
				{Filter: []string{"event-type"}, URL: "test-url"},
				{Filter: []string{"other-type"}, URL: "wrong-url"},
				{Filter: []string{"event-type"}, URL: "other-url"},
			},
			Expected: []string{"test-url", "other-url"},
		},
		{
			Name:  "wildcard match",
			Input: "event-type.will.resp",
			Routes: []Route{
				{Filter: []string{"event-type", "*"}, URL: "too-little"},
				{Filter: []string{"event-type", "*", "resp"}, URL: "just-right"},
				{Filter: []string{"event-type", "*", "response"}, URL: "late mismatch"},
				{Filter: []string{"event-type", "*", "resp", "another"}, URL: "too-much"},
				{Filter: []string{"other-type"}, URL: "wrong-url"},
			},
			Expected: []string{"just-right"},
		},
	} {
		t.Run(test.Name, func(u *testing.T) {
			router := Router{
				routes: test.Routes,
			}

			urls := router.GetURLs(test.Input)
			assert.Equal(u, test.Expected, urls)
		})
	}
}
