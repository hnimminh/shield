package api

import (
	"encoding/json"
	"io"
	"net/http"

	zlog "github.com/rs/zerolog/log"
)

func NftCommand(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		zlog.Error().Err(err).Str("function", "Shield:API:NftCommand").Msg("Fail to read request")
	}
	zlog.Info().Str("function", "Shield:API:NftCommand").Msgf("A new request to setting up %s", body)

	// response
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(SuccessResponse{"Success"})
}

func NftShow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// response
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(SuccessResponse{"Success"})
}
