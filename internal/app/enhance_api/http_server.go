package enhance_api

import (
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hellofresh/health-go/v5"

	"go-nimeth/internal/api"
)

func NewHttpServer(
	enhanceApiHandler *api.EnhanceApiHandler) *HttpServer {
	globalMux := chi.NewRouter()
	globalMux.Mount("/debug", middleware.Profiler())

	mux := globalMux.Group(nil)

	// A good base middleware stack
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Timeout(60 * time.Second))
	mux.Use(CorsPublic)

	// add some checks on instance creation
	h, _ := health.New(health.WithComponent(health.Component{
		Name:    "enhance_api",
		Version: "v1.0",
	}),
	)
	mux.Get("/", h.HandlerFunc)
	mux.Mount("/api/v1", enhanceApiHandler.Route())

	server := &http.Server{
		Handler: globalMux,
	}
	return &HttpServer{server}
}

func CorsPublic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS,HEAD,GET,POST,PUT,DELETE,PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		if r.Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, r)
	})
}

type HttpServer struct {
	*http.Server
}
