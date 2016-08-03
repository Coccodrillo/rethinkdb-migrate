package main

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/go-yaml/yaml"
	r "gopkg.in/dancannon/gorethink.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	Address  string `yaml:"address"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	CertFile string `yaml:"cert_file"`
}

func NewConfig() *Config {
	c := &Config{}
	data, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	err = yaml.Unmarshal([]byte(data), c)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return c
}

func GetSession(c *Config) (session *r.Session, err error) {
	var t *tls.Config
	if c.CertFile != "" {
		pem, err := ioutil.ReadFile(c.CertFile)
		if err != nil {
			log.Fatalf("Rethinkdb/SSL: %s", err)
		}
		t = &tls.Config{RootCAs: x509.NewCertPool()}
		t.RootCAs.AppendCertsFromPEM(pem)
	}
	session, err = r.Connect(r.ConnectOpts{
		Address:   c.Address,
		Database:  c.Database,
		Username:  c.Username,
		Password:  c.Password,
		TLSConfig: t,
	})
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return session, err
}
