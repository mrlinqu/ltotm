package webserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type WebServer struct {
	srv *http.Server

	addr    string
	sslCert string
	sslKey  string

	storage storage
}

func New(addr string, sslCert string, sslKey string, storage storage) *WebServer {
	return &WebServer{
		addr:    addr,
		sslCert: sslCert,
		sslKey:  sslKey,
		storage: storage,
	}
}

func (s *WebServer) Run(ctx context.Context) <-chan struct{} {
	done := make(chan struct{})

	s.srv = &http.Server{Addr: s.addr}

	http.HandleFunc("/", s.handlerRoot)
	http.HandleFunc("/put", s.handlerPut)
	http.HandleFunc("/get", s.handlerGet)

	go func() {
		defer close(done)

		var err error
		if s.sslCert != "" && s.sslKey != "" {
			err = s.srv.ListenAndServeTLS(s.sslCert, s.sslKey)
		} else {
			err = s.srv.ListenAndServe()
		}

		if err != nil && err != http.ErrServerClosed {
			log.Fatal().Msgf("failed to serve web: %v", err)
		}
	}()

	return done
}

func (s *WebServer) Shutdown(ctx context.Context) <-chan struct{} {
	done := make(chan struct{})

	tmCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	go func() {
		defer close(done)

		<-tmCtx.Done()
	}()

	go func() {
		defer cancel()

		log.Debug().Msg("shutdown web server...")

		if err := s.srv.Shutdown(tmCtx); err != nil && err != context.Canceled {
			log.Error().Msgf("fail shutdown web: %v", err)
		}
	}()

	return done
}

func (s *WebServer) errorHandler(w http.ResponseWriter, r *http.Request, status int, texts ...string) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(w, "Not Found")
	}

	for txt := range texts {
		fmt.Fprint(w, txt)
	}
}
