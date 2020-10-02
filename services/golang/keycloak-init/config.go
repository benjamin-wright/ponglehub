package main

import (
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/sirupsen/logrus"
)

type config struct {
	URL          string
	Realm        string
	Username     string
	Password     string
	SMTPEmail    string
	SMTPPassword string
	SMTPHost     string
	SMTPPort     int
	SMTPFrom     string
}

func setProp(cfg *config, property string, envvar string) error {
	if cfg == nil {
		return fmt.Errorf("Cannot set property %s on nil config", property)
	}

	configValue := reflect.ValueOf(cfg).Elem()
	configPropertyValue := configValue.FieldByName(property)

	if !configPropertyValue.IsValid() {
		return fmt.Errorf("No such property %s in config", property)
	}

	if !configPropertyValue.CanSet() {
		return fmt.Errorf("Cannot set %s property value", property)
	}

	value, ok := os.LookupEnv(envvar)
	if !ok {
		return fmt.Errorf("Value required for %s", envvar)
	}

	switch t := configPropertyValue.Interface().(type) {
	case string:
		configPropertyValue.Set(reflect.ValueOf(value))
	case int:
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("Failed to convert env var %s to int: %+v", envvar, err)
		}

		configPropertyValue.Set(reflect.ValueOf(intValue))
	default:
		return fmt.Errorf("Failed to convert env var %s to unknown type %t", envvar, t)
	}

	return nil
}

type propBinding struct {
	prop   string
	envVar string
}

func setProps(cfg *config, props []propBinding) error {
	for _, prop := range props {
		err := setProp(cfg, prop.prop, prop.envVar)
		if err != nil {
			return err
		}
	}

	return nil
}

func newConfig() (*config, error) {
	cfg := &config{}

	err := setProps(
		cfg,
		[]propBinding{
			{prop: "URL", envVar: "KEYCLOAK_INIT_URL"},
			{prop: "Realm", envVar: "KEYCLOAK_INIT_REALM"},
			{prop: "Username", envVar: "KEYCLOAK_INIT_USER"},
			{prop: "Password", envVar: "KEYCLOAK_INIT_PASSWORD"},
			{prop: "SMTPEmail", envVar: "KEYCLOAK_SMTP_EMAIL"},
			{prop: "SMTPPassword", envVar: "KEYCLOAK_SMTP_PASSWORD"},
			{prop: "SMTPHost", envVar: "KEYCLOAK_SMTP_HOST"},
			{prop: "SMTPPort", envVar: "KEYCLOAK_SMTP_PORT"},
			{prop: "SMTPFrom", envVar: "KEYCLOAK_SMTP_FROM"},
		},
	)

	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *config) print() {
	logrus.Infof("Config:\n - url: %s\n - realm: %s\n - username: %s\n - password: %t", c.URL, c.Realm, c.Username, c.Password != "")
}
