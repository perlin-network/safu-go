package api

import (
	"net/http"

	"github.com/perlin-network/safu-go/log"
	"github.com/rs/cors"
)

// service represents a service.
type service struct {
}

// init registers routes to the HTTP serve mux.
func (s *service) init(mux *http.ServeMux) {
	mux.Handle("/debug/vars", http.DefaultServeMux)
	mux.HandleFunc("/post_scam_report", s.wrap(s.postScamReport))
	mux.HandleFunc("/query_address", s.wrap(s.queryAddress))
}

// Run runs the API server with a specified set of options.
func Run(serverAddr string) {
	mux := http.NewServeMux()

	service := &service{}

	service.init(mux)

	handler := cors.AllowAll().Handler(mux)

	server := &http.Server{
		Addr:    serverAddr,
		Handler: handler,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal().Err(err).Msg(" ")
	}
}
