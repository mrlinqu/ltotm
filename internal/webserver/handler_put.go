package webserver

import (
	"encoding/json"
	"net/http"
	"time"

	message_id "github.com/mrlinqu/ltotm/internal/message-id"
	"github.com/rs/zerolog/log"
)

type putRequest struct {
	Msg string `json:"msg"`
	Ttl int64  `json:"ttl"`
}

type putResponse struct {
	ExpirationTime string `json:"expiration_time"`
	Id             string `json:"id"`
}

func (s *WebServer) handlerPut(w http.ResponseWriter, r *http.Request) {
	var req putRequest

	err := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxPutRequestSize)).Decode(&req)
	if err != nil {
		log.Error().Msgf("decode put request: %v", err)
		s.errorHandler(w, r, http.StatusBadRequest, "request decode error")
		return
	}

	if req.Ttl == 0 || req.Ttl > 72 {
		req.Ttl = 24
	}

	//fileId := message_id.GenerateFromNow(req.Ttl)

	ttl := time.Now().Add(time.Duration(req.Ttl) * time.Hour)
	fileId := message_id.Generate(ttl)

	err = s.storage.Put(fileId, []byte(req.Msg))
	if err != nil {
		log.Error().Msgf("put message: %v", err)
		s.errorHandler(w, r, http.StatusInternalServerError, "save message error")
		return
	}

	resp := putResponse{
		ExpirationTime: ttl.Format(time.DateTime),
		Id:             fileId,
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Error().Msgf("encode put response: %v", err)
		s.errorHandler(w, r, http.StatusInternalServerError, "encode response error")
		return
	}
}
