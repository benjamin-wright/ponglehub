package api

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
)

// AddRealm add the realm to the auth server
func (k *KeycloakAPI) AddRealm(name string) error {
	logrus.Debugf("Adding realm %s", name)
	requestBody, err := json.Marshal(map[string]string{
		"realm": name,
	})
	if err != nil {
		return err
	}

	err = k.post("/", requestBody, nil)
	if err != nil {
		return fmt.Errorf("Failed to create realm: %+v", err)
	}

	logrus.Infof("Added realm %s", name)
	return nil
}

// RemoveRealm delete the realm from the auth server
func (k *KeycloakAPI) RemoveRealm(name string) error {
	logrus.Debugf("Removing realm %s", name)
	err := k.delete("/"+name, nil)
	if err != nil {
		return fmt.Errorf("Failed to delete realm: %+v", err)
	}

	logrus.Infof("Removed realm %s", name)
	return nil
}
