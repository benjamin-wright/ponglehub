package api

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
)

// Realm a data representation of a Keycloak realm
type Realm struct {
	Name                  string      `json:"realm"`
	Display               string      `json:"displayName"`
	SMTPServer            *SMTPServer `json:"smtpServer"`
	Enabled               bool        `json:"enabled"`
	SSLRequired           string      `json:"sslRequired"`
	RegistrationAllowed   bool        `json:"registrationAllowed"`
	RememberMe            bool        `json:"rememberMe"`
	VerifyEmail           bool        `json:"verifyEmail"`
	LoginWithEmailAllowed bool        `json:"loginWithEmailAllowed"`
	ResetPasswordAllowed  bool        `json:"resetPasswordAllowed"`
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

func (r *Realm) String() string {
	return fmt.Sprintf(
		"{name: %s, display: %s, smptUser: %s, smptHost: %s, smtpPort: %d}",
		r.Name,
		r.Display,
		r.SMTPServer.User,
		r.SMTPServer.Host,
		r.SMTPServer.Port,
	)
}

// HasRealm checks if the realm exists already
func (k *KeycloakAPI) HasRealm(realm string) (bool, error) {
	logrus.Debugf("Checking for realm %s", realm)

	code, err := k.get("/"+realm, nil)
	if code == 0 {
		return false, fmt.Errorf("Failed to check for realm: %+v", err)
	}

	return code >= 200 && code < 300, nil
}

// AddRealm add the realm to the auth server
func (k *KeycloakAPI) AddRealm(realm Realm) error {
	logrus.Debugf("Adding realm %s", realm.String())
	requestBody, err := json.Marshal(realm)
	if err != nil {
		return err
	}

	_, err = k.post("/", requestBody, nil)
	if err != nil {
		return fmt.Errorf("Failed to create realm: %+v", err)
	}

	logrus.Infof("Added realm %s", realm.String())
	return nil
}

// RemoveRealm delete the realm from the auth server
func (k *KeycloakAPI) RemoveRealm(realm string) error {
	logrus.Debugf("Removing realm %s", realm)
	_, err := k.delete("/"+realm, nil)
	if err != nil {
		return fmt.Errorf("Failed to delete realm: %+v", err)
	}

	logrus.Infof("Removed realm %s", realm)
	return nil
}
