package runtime

import (
	"database/sql"
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
	services []*service.Service
	http.Handler
	group.Group
}

var RunGroup group.Group

func NewService(db *sql.DB, bucket *blob.Bucket, srv *server.Server, l requestlog.Logger, s *service.Service) *Service {
	var handl http.Handler
	m := mux.NewRouter()
	m.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	m.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	m.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	m.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	m.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
	m.Handle("/metrics", promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}))
	handl = requestlog.NewHandler(l, m)

	me := &Service{db: db, bucket: bucket, srv: srv, Handler: handl, Group: RunGroup}
	me.services = append(me.services, s)
	return me
}

func (s *Service) Services() []*service.Service {
	return s.services
}

func (s *Service) AddService(sv *service.Service) {
	s.services = append(s.services, sv)
}

func (a *Service) ResetServices() {
	a.services = nil
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
