package web

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hnimminh/shield/config"
	"github.com/hnimminh/shield/web/api"
	zlog "github.com/rs/zerolog/log"
)

// @title           Horseman Agent
// @version         0.0.0
// @description     Horseman Agent Restful API
// @contact.name    Minh Minh
// @contact.email   hnimminh@outlook.com
// @BasePath
func Server() {
	router := mux.NewRouter()

	// HEALTHCHECK
	router.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("+OK"))
	}).Methods("GET")

	// NFT
	router.HandleFunc("/nft", api.NftCommand).Methods("PUT", "POST")
	router.HandleFunc("/nft", api.NftShow).Methods("GET")

	// SERVING
	bindaddr := fmt.Sprintf("%s:%d", config.HTTPListenIP, config.HTTPListenPort)
	zlog.Info().Str("function", "Shield:Web:Server").Msgf("Worker - server is listen on %s", bindaddr)
	err := http.ListenAndServe(bindaddr, router)
	if err != nil {
		zlog.Fatal().Err(err).Str("function", "Shield:Web:Server").Msg("Worker Server fail to start")
	}
}
