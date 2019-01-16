package app

import (
	"crypto/ecdsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"fmt"
	"github.com/gofunct/goscript/app/service"
	"github.com/google/wire"
	"go.opencensus.io/trace"
	"gocloud.dev/blob"
	"gocloud.dev/health"
	"gocloud.dev/health/sqlhealth"
	"gocloud.dev/server"
	"google.golang.org/grpc"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// applicationSet is the Wire provider set for the Guestbook application that
// does not depend on the underlying platform.
var ApplicationSet = wire.NewSet(
	http.NewServeMux,
	NewApplication,
	newRuntime,
	appHealthChecks,
	trace.AlwaysSample,
)

type Application struct {
	services []*service.Service
	mux      http.Handler
	*runtime
}

func NewApplication(mux http.Handler, runtime *runtime) *Application {
	m := http.NewServeMux()
	m.Handle("/", mux)
	return &Application{mux: m, runtime: runtime}
}

func (a *Application) Services() []*service.Service {
	return a.services
}

func (a *Application) AddService(s *service.Service) {
	a.services = append(a.services, s)
}

type runtime struct {
	db     *sql.DB
	bucket *blob.Bucket
	srv    *server.Server
}

func newRuntime(db *sql.DB, bucket *blob.Bucket, srv *server.Server) *runtime {
	return &runtime{db: db, bucket: bucket, srv: srv}
}

func appHealthChecks(db *sql.DB) ([]health.Checker, func()) {
	dbCheck := sqlhealth.New(db)
	list := []health.Checker{dbCheck}
	return list, func() {
		dbCheck.Stop()
	}
}

func (a *Application) HasServices() bool {
	return len(a.services) > 0
}

func (a *Application) WalkProtoDir(d string) error {
	return filepath.Walk(d, protoWalkFunc)
}

var protoWalkFunc = func(path string, info os.FileInfo, err error) error {
	// skip vendor directory
	if info.IsDir() && info.Name() == "vendor" {
		return filepath.SkipDir
	}
	// find all protobuf files
	if filepath.Ext(path) == ".proto" {
		// args
		args := []string{
			"--go_out=plugins=grpc:.",
			path,
		}
		cmd := exec.Command("protoc", args...)
		cmd.Env = os.Environ()
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

func (a *Application) HandleGrpc(grpcServer *grpc.Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			a.mux.ServeHTTP(w, r)
		}
	})
}

// LoadECKeyFromFile loads EC key from unencrypted PEM file.
func (a *Application) LoadECKeyFromFile(fileName string) (*ecdsa.PrivateKey, error) {
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
