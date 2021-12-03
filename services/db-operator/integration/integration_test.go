package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/scheme"
	"ponglehub.co.uk/operators/db/internal/client"
	"ponglehub.co.uk/operators/db/internal/types"
)

func Test(t *testing.T) {
	client.AddToScheme(scheme.Scheme)

	cli, err := client.New()
	if err != nil {
		assert.NoError(t, err)
		t.FailNow()
	}

	if _, err := cli.DBGet("test-db", "test-namespace"); err == nil {
		if err := cli.DBDelete("test-db", "test-namespace"); err != nil {
			assert.NoError(t, err)
			t.FailNow()
		}
	}

	if err := cli.DBCreate(types.Database{
		Name:      "test-db",
		Namespace: "test-namespace",
		Storage:   "2G",
	}); err != nil {
		assert.NoError(t, err)
		t.FailNow()
	}

	t.Log("TBD")
}
