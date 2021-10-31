package connect

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"path"

	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
)

func getTlsConfig(certsDir string, username string) (*tls.Config, error) {
	clientCrt := path.Join(certsDir, fmt.Sprintf("client.%s.crt", username))
	clientKey := path.Join(certsDir, fmt.Sprintf("client.%s.key", username))
	caCrt := path.Join(certsDir, "ca.crt")

	cert, err := tls.LoadX509KeyPair(clientCrt, clientKey)
	if err != nil {
		return nil, fmt.Errorf("failed to load client certificate: %v", err)
	}

	CACert, err := ioutil.ReadFile(caCrt)
	if err != nil {
		return nil, fmt.Errorf("failed to load server certificate: %v", err)
	}

	CACertPool := x509.NewCertPool()
	CACertPool.AppendCertsFromPEM(CACert)

	return &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            CACertPool,
		InsecureSkipVerify: true,
	}, nil
}

func getConnection(config *pgx.ConnConfig) *pgx.Conn {
	finished := make(chan *pgx.Conn, 1)

	go func(finished chan<- *pgx.Conn) {
		attempts := 0
		limit := 10
		var connection *pgx.Conn
		var err error
		for attempts < limit {
			connection, err = pgx.ConnectConfig(context.Background(), config)
			if err != nil {
				logrus.Warnf("error connecting to the database: %+v", err)
			} else {
				break
			}
		}

		finished <- connection
	}(finished)

	return <-finished
}

func Connect(config ConnectConfig) (*pgx.Conn, error) {
	userPass := config.Username
	if config.Password != "" {
		userPass += ":" + config.Password
	}

	dbSuffix := ""
	if config.Database != "" {
		dbSuffix = "/" + config.Database
	}

	pgxConfig, err := pgx.ParseConfig(fmt.Sprintf("postgresql://%s@%s:%d%s", userPass, config.Host, config.Port, dbSuffix))
	if err != nil {
		return nil, err
	}

	tlsConfig, err := getTlsConfig(config.CertsDir, config.Username)
	if err != nil {
		return nil, err
	}

	pgxConfig.TLSConfig = tlsConfig
	conn := getConnection(pgxConfig)
	if conn == nil {
		return nil, errors.New("failed to create connection, exiting")
	}

	return conn, nil
}
