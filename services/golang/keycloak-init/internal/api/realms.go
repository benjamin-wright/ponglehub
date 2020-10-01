package api

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
)

// Realm a data representation of a Keycloak realm
type Realm struct {
	Name       string      `json:"realm"`
	Display    string      `json:"displayName"`
	SMTPServer *SMTPServer `json:"smtpServer"`
}

// SMTPServer the settings for SMTP in a realm
type SMTPServer struct {
	User               string `json:"user"`
	Password           string `json:"password"`
	ReplyToDisplayName string `json:"replyToDisplayName"`
	StartTLS           bool   `json:"starttls"`
	Auth               bool   `json:"auth"`
	Port               int    `json:"port"`
	Host               string `json:"host"`
	ReplyTo            string `json:"replyTo"`
	From               string `json:"from"`
	FromDisplayName    string `json:"fromDisplayName"`
	SSL                bool   `json:"ssl"`
}

// HasRealm check if the realm exists already
func (k *KeycloakAPI) HasRealm(realm Realm) (bool, error) {
	logrus.Debugf("Checking for realm %s", realm.Name)

	code, err := k.get("/"+realm.Name, nil)
	if code == 0 {
		return false, fmt.Errorf("Failed to check for realm: %+v", err)
	}

	return code >= 200 && code < 300, nil
}

// AddRealm add the realm to the auth server
func (k *KeycloakAPI) AddRealm(realm Realm) error {
	logrus.Debugf("Adding realm %s", realm.Name)
	requestBody, err := json.Marshal(realm)
	if err != nil {
		return err
	}

	_, err = k.post("/", requestBody, nil)
	if err != nil {
		return fmt.Errorf("Failed to create realm: %+v", err)
	}

	logrus.Infof("Added realm %s", realm.Name)
	return nil
}

// RemoveRealm delete the realm from the auth server
func (k *KeycloakAPI) RemoveRealm(realm Realm) error {
	logrus.Debugf("Removing realm %s", realm.Name)
	_, err := k.delete("/"+realm.Name, nil)
	if err != nil {
		return fmt.Errorf("Failed to delete realm: %+v", err)
	}

	logrus.Infof("Removed realm %s", realm.Name)
	return nil
}
