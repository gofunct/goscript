package runtime

import (
	"crypto/ecdsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"fmt"
	"github.com/gofunct/goscript/runtime/health"
	"github.com/gofunct/goscript/runtime/router"
	"github.com/gofunct/goscript/service"
	"github.com/google/wire"
	"github.com/oklog/oklog/pkg/group"
	"go.opencensus.io/trace"
	"gocloud.dev/blob"
	"gocloud.dev/server"
	"io/ioutil"
	"net/http"
)

var Set = wire.NewSet(
	NewService,
	trace.AlwaysSample,
	router.Set,
	health.Set,
	service.Set,
)

type Service struct {
	db       *sql.DB
	bucket   *blob.Bucket
	srv      *server.Server
	services *service.Service
	http.Handler
	group group.Group
}

var RunGroup group.Group

func NewService(db *sql.DB, bucket *blob.Bucket, srv *server.Server, h http.Handler, services *service.Service) *Service {
	return &Service{db: db, bucket: bucket, srv: srv, services: services, Handler: h, group: RunGroup}
}

func (s *Service) Services() *service.Service {
	return s.services
}

func (a *Service) ResetServices() {
	a.services = nil
}

// LoadECKeyFromFile loads EC key from unencrypted PEM file.
func (s *Service) LoadECKeyFromFile(fileName string) (*ecdsa.PrivateKey, error) {
	privateKeyBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to read signing key file: %v", err)
	}

	privateKeyPEM, _ := pem.Decode(privateKeyBytes)
	if privateKeyPEM == nil {
		return nil, fmt.Errorf("failed to decode pem signing key file: %v", err)
	}

	privateKey, err := x509.ParseECPrivateKey(privateKeyPEM.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse signing key file: %v", err)
	}

	return privateKey, nil
}
