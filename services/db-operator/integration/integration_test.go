package integration

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/scheme"
	"ponglehub.co.uk/lib/postgres/pkg/connect"
	"ponglehub.co.uk/operators/db/internal/crds"
	"ponglehub.co.uk/operators/db/internal/types"
)

func Test(t *testing.T) {
	crds.AddToScheme(scheme.Scheme)

	t.Run("Runs a database", simpleSetup)
	t.Run("Can connect", canConnect)
}

func simpleSetup(t *testing.T) {
	testDB := types.Database{
		Name:      "test-db",
		Namespace: "test-namespace",
		Storage:   "2G",
	}

	helper := newHelper(t)
	helper.ensureNoDB(t, testDB)
	helper.createDb(t, testDB)
	helper.waitForRunning(t, testDB)
	result := helper.getDb(t, testDB)

	assert.Equal(
		t,
		types.Database{
			Name:      "test-db",
			Namespace: "test-namespace",
			Storage:   "2G",
			Ready:     true,
		},
		result,
	)
}

func canConnect(t *testing.T) {
	testDB := types.Database{
		Name:      "test-db",
		Namespace: "test-namespace",
		Storage:   "2G",
	}

	testClient := types.Client{
		Name:       "test-client",
		Deployment: "test-db",
		Username:   "test_client",
		Namespace:  "test-namespace",
		Database:   "auth_test",
		Secret:     "test-secret",
		Ready:      false,
	}

	helper := newHelper(t)
	helper.ensureNoDB(t, testDB)
	helper.createDb(t, testDB)
	helper.waitForRunning(t, testDB)

	helper.ensureNoClient(t, testClient)
	helper.createClient(t, testClient)
	helper.waitForClientSecret(t, testClient)

	_, err := connect.Connect(connect.ConnectConfig{
		Host:     fmt.Sprintf("%s.%s.svc.cluster.local", testDB.Name, testDB.Namespace),
		Port:     26257,
		Username: testClient.Username,
		Database: testClient.Database,
	})
	assert.NoError(t, err)
}
