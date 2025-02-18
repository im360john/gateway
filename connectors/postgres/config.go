package postgres

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

type Config struct {
	Hosts     []string
	Database  string
	User      string
	Password  string
	Port      int
	TLSFile   string
	EnableTLS bool
}

func (c Config) TLSConfig() (*tls.Config, error) {
	if c.EnableTLS {
		rootCertPool := x509.NewCertPool()
		if len(c.TLSFile) > 0 {
			if ok := rootCertPool.AppendCertsFromPEM([]byte(c.TLSFile)); !ok {
				return nil, errors.New("unable to add TLS to cert pool")
			}
		}
		return &tls.Config{
			RootCAs:    rootCertPool,
			ServerName: c.Hosts[0],

			InsecureSkipVerify: len(c.TLSFile) == 0,
		}, nil
	}

	return nil, nil
}

func (c Config) MakeConfig() (*pgx.ConnConfig, error) {
	tlsConfig, err := c.TLSConfig()
	if err != nil {
		return nil, err
	}
	config, _ := pgx.ParseConfig("")
	config.Host = c.Hosts[0]
	config.Port = uint16(c.Port)
	config.Database = c.Database
	config.User = c.User
	config.Password = string(c.Password)
	config.TLSConfig = tlsConfig
	config.PreferSimpleProtocol = true
	return config, nil
}

func (c Config) Type() string {
	return "postgres"
}
