package integration

import (
	"testing"

	"k8s.io/client-go/kubernetes/scheme"
	"ponglehub.co.uk/operators/db/internal/crds"
)

func Test(t *testing.T) {
	crds.AddToScheme(scheme.Scheme)

	t.Run("Can connect", canConnect)
	t.Run("Out of order", outOfOrderConnect)
}

func canConnect(t *testing.T) {
	testDB := crds.Database{
		Name:      "test-1-db",
		Namespace: "test-namespace",
		Storage:   "2G",
	}

	testClient := crds.Client{
		Name:       "test-1-client",
		Deployment: "test-1-db",
		Username:   "test_client",
		Namespace:  "test-namespace",
		Database:   "auth_1_test",
		Secret:     "test-1-secret",
		Ready:      false,
	}

	helper := newHelper(t)
	helper.ensureNoDB(t, testDB)
	helper.ensureNoClient(t, testClient)

	helper.createDb(t, testDB)
	helper.waitForRunning(t, testDB)

	helper.createClient(t, testClient)
	helper.waitForClientSecret(t, testClient)
}

func outOfOrderConnect(t *testing.T) {
	testDB := crds.Database{
		Name:      "other-db",
		Namespace: "test-namespace",
		Storage:   "2G",
	}

	testClient := crds.Client{
		Name:       "other-client",
		Deployment: "other-db",
		Username:   "other_client",
		Namespace:  "test-namespace",
		Database:   "other_test",
		Secret:     "other-secret",
		Ready:      false,
	}

	helper := newHelper(t)
	helper.ensureNoClient(t, testClient)
	helper.ensureNoDB(t, testDB)

	helper.createClient(t, testClient)
	helper.createDb(t, testDB)
	helper.waitForRunning(t, testDB)
	helper.waitForClientSecret(t, testClient)
}
