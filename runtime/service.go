package runtime

import (
	"crypto/ecdsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"fmt"
	"github.com/gofunct/goscript/runtime/health"
	"github.com/gofunct/goscript/service"
	"github.com/google/wire"
	"github.com/gorilla/mux"
	"github.com/oklog/oklog/pkg/group"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opencensus.io/trace"
	"gocloud.dev/blob"
	"gocloud.dev/requestlog"
	"gocloud.dev/server"
	"google.golang.org/grpc"
	"io/ioutil"
	"net/http"
	"net/http/pprof"
	"strings"
)

var Set = wire.NewSet(
	NewService,
	trace.AlwaysSample,
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

func NewService(db *sql.DB, bucket *blob.Bucket, srv *server.Server, services *service.Service, l requestlog.Logger) *Service {
	s := &Service{}
	m := mux.NewRouter()
	m.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	m.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	m.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	m.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	m.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
	m.Handle("/metrics", promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}))
	handl := s.HandleGrpc(services.Server)
	handl = requestlog.NewHandler(l, m)

	return &Service{db: db, bucket: bucket, srv: srv, services: services, Handler: handl, group: RunGroup}
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

func (s *Service) HandleGrpc(grpcServer *grpc.Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			ctx, span := trace.StartSpan(r.Context(), r.URL.Host+r.URL.Path)
			defer span.End()

			r = r.WithContext(ctx)

			grpcServer.ServeHTTP(w, r)
		} else {
			ctx, span := trace.StartSpan(r.Context(), r.URL.Host+r.URL.Path)
			defer span.End()

			r = r.WithContext(ctx)
			s.Handler.ServeHTTP(w, r)
		}
	})
}
