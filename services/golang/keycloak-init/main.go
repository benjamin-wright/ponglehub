package main

import (
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/keycloak-init/internal/api"
)

func main() {
	logrus.Info("Starting...")

	cfg, err := newConfig()

	if err != nil {
		logrus.Fatalf("Failed to load config: %+v", err)
	}

	logrus.Infof("Config: %s", cfg.String())

	keycloak, err := api.New(cfg.URL, cfg.Username, cfg.Password)
	if err != nil {
		logrus.Fatalf("Failed to create api instance: %+v", err)
	}

	realm := api.Realm{
		Name:    cfg.Realm,
		Display: cfg.Realm,
		SMTPServer: &api.SMTPServer{
			User:               cfg.SMTPEmail,
			Password:           cfg.SMTPPassword,
			ReplyToDisplayName: cfg.SMTPFrom,
			StartTLS:           false,
			Auth:               true,
			Port:               cfg.SMTPPort,
			Host:               cfg.SMTPHost,
			ReplyTo:            cfg.SMTPEmail,
			From:               cfg.SMTPEmail,
			FromDisplayName:    cfg.SMTPFrom,
			SSL:                true,
		},
		Enabled:               true,
		SSLRequired:           "all",
		RegistrationAllowed:   true,
		RememberMe:            true,
		VerifyEmail:           true,
		LoginWithEmailAllowed: true,
		ResetPasswordAllowed:  true,
	}

	exists, err := keycloak.HasRealm(cfg.Realm)
	if err != nil {
		logrus.Fatalf("Failed checking for realm: %+v", err)
	}

	if exists {
		err = keycloak.RemoveRealm(cfg.Realm)
		if err != nil {
			logrus.Fatalf("Failed to remove realm: %+v", err)
		}
	} else {
		logrus.Infof("%s realm not currently installed", cfg.Realm)
	}

	err = keycloak.AddRealm(realm)
	if err != nil {
		logrus.Fatalf("Failed to add realm: %+v", err)
	}

	logrus.Info("Finished")
}
