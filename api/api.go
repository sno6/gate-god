package api

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/sno6/gate-god/relay"
)

type API struct {
	mux    *http.ServeMux
	logger *log.Logger
	r      *relay.Relay
}

func New(r *relay.Relay) *API {
	api := &API{
		mux:    http.NewServeMux(),
		r:      r,
		logger: log.New(os.Stdout, "[API]: ", log.LstdFlags),
	}
	api.mux.Handle("/relay", http.HandlerFunc(api.HandleRelay))
	return api
}

func (api *API) Serve(port int) error {
	api.logger.Printf("Starting server on port: %d\n", port)
	return http.ListenAndServe(":"+strconv.Itoa(port), api.mux)
}

func (api *API) HandleRelay(_ http.ResponseWriter, _ *http.Request) {
	api.logger.Println("Running relay manually")
	api.r.Toggle()
}
