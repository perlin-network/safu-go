package api

import (
	"github.com/perlin-network/safu-go/database"
	"github.com/perlin-network/safu-go/etherscan"
	"github.com/perlin-network/safu-go/ledger"
	"gopkg.in/go-playground/validator.v9"
	"net/http"

	"github.com/perlin-network/safu-go/log"
	"github.com/rs/cors"
)

const (
	RoutePostScamRepot = "/post_scam_report"
	RouteQueryAddress  = "/query_address"
)

var (
	validate = validator.New()
)

// service represents a service.
type service struct {
	esClient *etherscan.ESClient
	store    *database.TieDotStore
	ledger   *ledger.Ledger
}

// init registers routes to the HTTP serve mux.
func (s *service) init(mux *http.ServeMux) {
	mux.HandleFunc(RoutePostScamRepot, s.wrap(s.postScamReport))
	mux.HandleFunc(RouteQueryAddress, s.wrap(s.queryAddress))
}

// Run runs the API server with a specified set of options.
func Run(serverAddr string, esClient *etherscan.ESClient, store *database.TieDotStore, ledger *ledger.Ledger) {
	mux := http.NewServeMux()

	service := &service{
		esClient: esClient,
		store:    store,
		ledger:   ledger,
	}

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
