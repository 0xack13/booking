package api

import (
	"context"
	"fmt"
	"github.com/xsolrac87/booking/booking"
	"log"
	"net/http"
	"time"
)

const (
	defaultHTTPPort = 8055
	defaultTimeout  = 10 * time.Second
)

type HttpServer struct {
	*http.Server
	timeout time.Duration
}

func NewHTTPServer(addr string, opts ...Option) (*HttpServer, error) {
	var options options
	for _, opt := range opts {
		err := opt(&options)
		if err != nil {
			return nil, err
		}
	}

	var port int
	if options.port != nil {
		port = *options.port
	} else {
		port = defaultHTTPPort
	}

	var timeout time.Duration
	if options.timeout != nil {
		timeout = *options.timeout
	} else {
		timeout = defaultTimeout
	}

	s := &HttpServer{
		Server: &http.Server{
			Addr: fmt.Sprintf("%s:%d", addr, port),
		},
		timeout: timeout,
	}

	s.registerRoutes()
	return s, nil
}

func (s *HttpServer) Run(ctx context.Context) error {
	go func() {
		log.Println("HTTP Server running on", s.Addr)
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("HTTP server shut down", err)
		}
	}()

	<-ctx.Done()
	ctxShutDown, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	log.Println("shutting down HTTP server")
	return s.Shutdown(ctxShutDown)
}

func (s *HttpServer) registerRoutes() {
	// Booking
	http.HandleFunc("/stats", booking.HandleR.HandlerStats)
	http.HandleFunc("/maximize", booking.HandleR.HandlerMaximize)
}
