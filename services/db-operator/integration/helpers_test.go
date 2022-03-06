package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/api/errors"
	"ponglehub.co.uk/operators/db/internal/crds"
	"ponglehub.co.uk/operators/db/internal/deployments"
)

type helpers struct {
	client *crds.DBClient
	depl   *deployments.DeploymentsClient
}

func newHelper(t *testing.T) *helpers {
	cli, err := crds.New()
	if err != nil {
		assert.NoError(t, err)
		t.FailNow()
	}

	depl, err := deployments.New()
	if err != nil {
		assert.NoError(t, err)
		t.FailNow()
	}

	return &helpers{client: cli, depl: depl}
}

func (h *helpers) ensureNoDB(t *testing.T, db crds.Database) {
	_, err := h.client.DBGet(db.Name, db.Namespace)
	if errors.IsNotFound(err) {
		return
	}

	err = h.client.DBDelete(db.Name, db.Namespace)
	if err != nil {
		assert.NoError(t, err)
		t.FailNow()
	}
}

func (h *helpers) createDb(t *testing.T, db crds.Database) {
	if err := h.client.DBCreate(db); err != nil {
		assert.NoError(t, err)
		t.FailNow()
	}
}

func (h *helpers) createClient(t *testing.T, client crds.Client) {
	if err := h.client.ClientCreate(client); err != nil {
		assert.NoError(t, err)
		t.FailNow()
	}
}

func waitFor(t *testing.T, name string, f func() bool) {
	ctx, _ := context.WithTimeout(context.Background(), 120*time.Second)
	finished := ctx.Done()

	for {
		if f() {
			return
		}

		select {
		case <-finished:
			t.Errorf("Timed out waiting for %s", name)
			t.FailNow()
		default:
			time.Sleep(time.Millisecond * 100)
		}
	}
}

func (h *helpers) getDb(t *testing.T, db crds.Database) crds.Database {
	current, err := h.client.DBGet(db.Name, db.Namespace)
	if err != nil {
		assert.NoError(t, err)
		t.FailNow()
	}

	return current
}

func (h *helpers) waitForRunning(t *testing.T, db crds.Database) {
	waitFor(t, fmt.Sprintf("database %s (%s)", db.Name, db.Namespace), func() bool {
		db, err := h.client.DBGet(db.Name, db.Namespace)
		if err != nil {
			assert.NoError(t, err)
			t.FailNow()
		}

		return db.Ready
	})
}

func (h *helpers) ensureNoClient(t *testing.T, client crds.Client) {
	_, err := h.client.ClientGet(client.Name, client.Namespace)
	if errors.IsNotFound(err) {
		return
	}

	err = h.client.ClientDelete(client.Name, client.Namespace)
	if err != nil {
		assert.NoError(t, err)
		t.FailNow()
	}
}

func (h *helpers) waitForClientSecret(t *testing.T, client crds.Client) {
	waitFor(t, fmt.Sprintf("client ready %s (%s)", client.Name, client.Namespace), func() bool {
		c, err := h.client.ClientGet(client.Name, client.Namespace)
		if err != nil {
			assert.NoError(t, err)
			t.FailNow()
		}

		return c.Ready
	})
}
