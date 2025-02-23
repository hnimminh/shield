package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/hnimminh/shield/internal/nftcli"
	zlog "github.com/rs/zerolog/log"
)

func NftCommand(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		zlog.Error().Err(err).Str("function", "Shield:API:NftCommand").Msg("Fail to read request")
	}
	zlog.Info().Str("function", "Shield:API:NftCommand").Msgf("A new request to setting up %s", body)

	_, err = nftcli.NftOverShell(string(body))
	if err != nil {
		zlog.Error().Err(err).Str("function", "Shield:API:NftCommand").Msg("Fail to executed")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(FailureResponse{"Failure", "InvalidCommand"})
		return
	}

	// response
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(SuccessResponse{"Success"})
}

func NftShow(w http.ResponseWriter, r *http.Request) {
	result, err := nftcli.NftOverShell("nft list ruleset")
	if err != nil {
		zlog.Error().Err(err).Str("function", "Shield:API:NftCommand").Msg("Fail to executed")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(403)
		json.NewEncoder(w).Encode(FailureResponse{"Failure", "Prohibited"})
		return
	}

	// response
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(200)
	w.Write([]byte(result))
}
