package reconciler

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"ponglehub.co.uk/operators/db/internal/actions"
	"ponglehub.co.uk/operators/db/internal/crds"
	"ponglehub.co.uk/operators/db/internal/deployments"
)

func TestReconcile(t *testing.T) {
	for _, test := range []struct {
		name     string
		setup    func(u *testing.T, r *Reconciler)
		expected []interface{}
	}{
		{
			name:     "empty",
			setup:    func(u *testing.T, r *Reconciler) {},
			expected: []interface{}{},
		},
		{
			name: "create db",
			setup: func(u *testing.T, r *Reconciler) {
				r.SetDatabase(crds.Database{Name: "my-db", Namespace: "hi", Storage: "1Gi"})
			},
			expected: []interface{}{
				actions.AddStatefulSet{StatefulSet: deployments.StatefulSet{Name: "my-db", Namespace: "hi", Storage: "1Gi"}},
				actions.AddService{Service: deployments.Service{Name: "my-db", Namespace: "hi"}},
			},
		},
		{
			name: "missing service",
			setup: func(u *testing.T, r *Reconciler) {
				r.SetDatabase(crds.Database{Name: "my-db", Namespace: "hi", Storage: "1Gi"})
				r.SetStatefulSet(deployments.StatefulSet{Name: "my-db", Namespace: "hi", Storage: "1Gi"})
			},
			expected: []interface{}{
				actions.AddService{Service: deployments.Service{Name: "my-db", Namespace: "hi"}},
			},
		},
		{
			name: "missing statefulset",
			setup: func(u *testing.T, r *Reconciler) {
				r.SetDatabase(crds.Database{Name: "my-db", Namespace: "hi", Storage: "1Gi"})
				r.SetService(deployments.Service{Name: "my-db", Namespace: "hi"})
			},
			expected: []interface{}{
				actions.AddStatefulSet{StatefulSet: deployments.StatefulSet{Name: "my-db", Namespace: "hi", Storage: "1Gi"}},
			},
		},
		{
			name: "update statefulset",
			setup: func(u *testing.T, r *Reconciler) {
				r.SetDatabase(crds.Database{Name: "my-db", Namespace: "hi", Storage: "1Gi"})
				r.SetService(deployments.Service{Name: "my-db", Namespace: "hi"})
				r.SetStatefulSet(deployments.StatefulSet{Name: "my-db", Namespace: "hi", Storage: "2Gi"})
			},
			expected: []interface{}{
				actions.DeleteStatefulSet{StatefulSet: deployments.StatefulSet{Name: "my-db", Namespace: "hi", Storage: "2Gi"}},
				actions.AddStatefulSet{StatefulSet: deployments.StatefulSet{Name: "my-db", Namespace: "hi", Storage: "1Gi"}},
			},
		},
		{
			name: "doesn't re-issue create",
			setup: func(u *testing.T, r *Reconciler) {
				r.SetDatabase(crds.Database{Name: "my-db", Namespace: "hi", Storage: "1Gi"})
				r.Reconcile(make(chan interface{}, 2))
			},
			expected: []interface{}{},
		},
		{
			name: "doesn't re-issue delete",
			setup: func(u *testing.T, r *Reconciler) {
				r.SetService(deployments.Service{Name: "my-db", Namespace: "hi"})
				r.SetStatefulSet(deployments.StatefulSet{Name: "my-db", Namespace: "hi", Storage: "2Gi"})
				r.Reconcile(make(chan interface{}, 2))
			},
			expected: []interface{}{},
		},
		{
			name: "remove existing",
			setup: func(u *testing.T, r *Reconciler) {
				r.SetService(deployments.Service{Name: "my-db", Namespace: "hi"})
				r.SetStatefulSet(deployments.StatefulSet{Name: "my-db", Namespace: "hi", Storage: "2Gi"})
			},
			expected: []interface{}{
				actions.DeleteStatefulSet{StatefulSet: deployments.StatefulSet{Name: "my-db", Namespace: "hi", Storage: "2Gi"}},
				actions.DeleteService{Service: deployments.Service{Name: "my-db", Namespace: "hi"}},
			},
		},
		{
			name: "update a pending database",
			setup: func(u *testing.T, r *Reconciler) {
				r.SetDatabase(crds.Database{Name: "my-db", Namespace: "hi", Storage: "1Gi"})
				r.Reconcile(make(chan interface{}, 2))
				r.SetDatabase(crds.Database{Name: "my-db", Namespace: "hi", Storage: "2Gi"})
			},
			expected: []interface{}{
				actions.DeleteStatefulSet{StatefulSet: deployments.StatefulSet{Name: "my-db", Namespace: "hi", Storage: "1Gi"}},
				actions.AddStatefulSet{StatefulSet: deployments.StatefulSet{Name: "my-db", Namespace: "hi", Storage: "2Gi"}},
			},
		},
		{
			name: "delete after pending create",
			setup: func(u *testing.T, r *Reconciler) {
				r.SetDatabase(crds.Database{Name: "my-db", Namespace: "hi", Storage: "1Gi"})
				r.Reconcile(make(chan interface{}, 2))
				r.RemoveDatabase(crds.Database{Name: "my-db", Namespace: "hi", Storage: "1Gi"})
			},
			expected: []interface{}{
				actions.DeleteStatefulSet{StatefulSet: deployments.StatefulSet{Name: "my-db", Namespace: "hi", Storage: "1Gi"}},
				actions.DeleteService{Service: deployments.Service{Name: "my-db", Namespace: "hi"}},
			},
		},
		{
			name: "re-issue pending create after pending delete",
			setup: func(u *testing.T, r *Reconciler) {
				r.SetDatabase(crds.Database{Name: "my-db", Namespace: "hi", Storage: "1Gi"})
				r.Reconcile(make(chan interface{}, 2))
				r.RemoveDatabase(crds.Database{Name: "my-db", Namespace: "hi", Storage: "1Gi"})
				r.Reconcile(make(chan interface{}, 2))
				r.SetDatabase(crds.Database{Name: "my-db", Namespace: "hi", Storage: "1Gi"})
			},
			expected: []interface{}{
				actions.AddStatefulSet{StatefulSet: deployments.StatefulSet{Name: "my-db", Namespace: "hi", Storage: "1Gi"}},
				actions.AddService{Service: deployments.Service{Name: "my-db", Namespace: "hi"}},
			},
		},
	} {
		t.Run(test.name, func(u *testing.T) {
			r := New()

			test.setup(u, r)

			events := make(chan interface{})
			go func(r *Reconciler, events chan<- interface{}) {
				r.Reconcile(events)
				close(events)
			}(r, events)

			results := []interface{}{}
			for event := range events {
				results = append(results, event)
			}

			assert.Equal(u, test.expected, results)
		})
	}
}
