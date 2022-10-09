package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/rs/cors"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Server struct {
	http.Server
	// Mux is the handler in http.Server
	mux *http.ServeMux
}

type In struct {
	fx.In

	LC         fx.Lifecycle
	Shutdowner fx.Shutdowner
	Logger     *zap.Logger
	Config     *Config
}

//nolint:gocritic // fx
func New(in In) (*Server, error) {
	var (
		logger = in.Logger.Named("http")
		mux    = http.NewServeMux()
	)
	addr := in.Config.Addr
	if addr == "" {
		addr = "127.0.0.1:8081"
	}
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://127.0.0.1:*",
			"http://localhost:*",
		},
	})
	srv := &Server{
		Server: http.Server{
			ReadHeaderTimeout: time.Second,
			Handler:           httpLog(logger, corsHandler.Handler(mux)),
			Addr:              addr,
		},
		mux: mux,
	}

	// wait OnStart goroutine inside OnStop
	var wg sync.WaitGroup

	in.LC.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer in.Shutdowner.Shutdown() //nolint:errcheck // shutdown
				logger.Info(fmt.Sprintf("Starting http server at: http://%s", addr))
				if err := srv.ListenAndServe(); err != nil {
					if !errors.Is(err, http.ErrServerClosed) {
						logger.Error("http server returned error", zap.Error(err))
					}
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			defer wg.Wait()
			err := srv.Shutdown(ctx)
			if err != nil {
				logger.Error("http shutdown failed", zap.Error(err))
				return err
			}
			return nil
		},
	})
	return srv, nil
}

func (s *Server) RegisterHandler(path string, handler http.Handler) {
	s.mux.Handle(path, handler)
}

func httpLog(logger *zap.Logger, handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("http request",
			zap.String("host", r.Host),
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("query", r.URL.Query().Encode()),
			zap.String("remote", r.RemoteAddr),
		)
		handler.ServeHTTP(w, r)
	}
}
