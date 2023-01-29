package mongobrocker

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"

	"github.com/phuhao00/broker"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Client struct {
	*broker.BaseComponent
	RealCli *mongo.Client
}

func NewClient(ctx context.Context, config *Config) *mongo.Client {
	opt := options.Client().ApplyURI(config.URI)
	opt.SetMinPoolSize(config.MinPoolSize).SetMaxPoolSize(config.MaxPoolSize)
	client, err := mongo.Connect(ctx, opt)
	if err != nil {
		panic(err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(err)
	}
	return client
}

func NewClientWithTLS() (*mongo.Client, error) {
	uri := "mongodb://user:password@localhost/?replicaSet=replset&authSource=admin"
	opts := options.Client().ApplyURI(uri)
	caBytes, _ := ioutil.ReadFile("/etc/ssl/certs/ca.pem")
	clientBytes, _ := ioutil.ReadFile("/etc/ssl/certs/client.pem")
	roots := x509.NewCertPool()
	if ok := roots.AppendCertsFromPEM(caBytes); !ok {
		return nil, errors.New("failed to parse root certificate")
	}
	certs, err := tls.X509KeyPair(clientBytes, clientBytes)
	if err != nil {
		return nil, err
	}
	cfg := &tls.Config{RootCAs: roots, Certificates: []tls.Certificate{certs}}
	opts.SetTLSConfig(cfg)
	client, err := mongo.Connect(context.Background(), opts)
	return client, err
}
