package config

import (
	"crypto/tls"
	"crypto/x509"
	"embed"
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

type Configuration struct {
	App  *AppConf
	Log  *logConf
	Db   *Db
	Nats *NatsConf
	Grpc *GrpcConf
}

type AppConf struct {
	VendorCloudCode string
}
type GrpcConf struct {
	Host       string
	Port       int
	Middleware struct {
		Log      bool
		Recovery bool
		Metrics  bool
	}
	Size struct {
		Send    int
		Receive int
	}
	UseTLS     bool
	ServerName string
}

func (r *AppConf) validate() error {
	return nil
}

type NatsConf struct {
	Enable   bool
	Host     string
	Port     int
	User     string
	Password string
	Tls      bool
	TLSConf  *tls.Config
}

func (r *NatsConf) validate() error {
	if !r.Enable {
		return nil
	}
	if r.Port == 0 {
		r.Port = 4222
	}
	if r.Host == "" {
		r.Host = "localhost"
	}
	if r.Tls {
		var err error
		r.TLSConf, err = createTLSConfig()
		if err != nil {
			return err
		}
	}

	// read username and password from environment variables
	r.User = os.Getenv("NATS_USER")
	r.Password = os.Getenv("NATS_PASSWORD")
	return nil
}

//go:embed mtls/*
var mtlsFS embed.FS

func createTLSConfig() (*tls.Config, error) {
	mustReadFile := func(path string) ([]byte, error) {
		data, err := mtlsFS.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read embedded file %s: %v", path, err)
		}
		return data, nil
	}

	clientCrt, err := mustReadFile("mtls/client.crt")
	if err != nil {
		return nil, err
	}
	clientKey, err := mustReadFile("mtls/client.key")
	if err != nil {
		return nil, err
	}
	caCert, err := mustReadFile("mtls/ca.crt")
	if err != nil {
		return nil, err
	}

	clientCert, err := tls.X509KeyPair(clientCrt, clientKey)
	if err != nil {
		return nil, fmt.Errorf("could not load client key pair: %w", err)
	}

	// caCertPool := x509.NewCertPool()
	caCertPool, err := x509.SystemCertPool()
	if err != nil || caCertPool == nil {
		caCertPool = x509.NewCertPool()
	}
	if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
		return nil, fmt.Errorf("failed to append CA cert to pool")
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      caCertPool,
	}
	return tlsConfig, nil
}

type logConf struct {
	Level     string
	Grpc      bool
	OnConsole bool
	IsJson    bool
	Trace     bool
	MaxAge    int
	MaxSize   int
}

type Db struct {
	Name           string
	Host           string
	Port           int
	User           string
	Password       string
	IdleConnection int
	OpenConnection int
}

func loadConf() *Configuration {
	c := Configuration{}
	viper.SetConfigFile("config.yaml")
	viper.AddConfigPath(".")

	if e := viper.ReadInConfig(); e != nil {
		log.Fatal(e)
	}
	if e := viper.Unmarshal(&c); e != nil {
		log.Fatal(e)
	}
	return &c
}
