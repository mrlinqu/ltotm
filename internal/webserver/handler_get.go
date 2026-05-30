package webserver

import (
	"encoding/json"
	"net/http"
	"time"

	message_id "github.com/mrlinqu/ltotm/internal/message-id"
	stor "github.com/mrlinqu/ltotm/internal/storage"
	"github.com/rs/zerolog/log"
)

type getRequest struct {
	Id string `json:"id"`
}

type getResponse struct {
	Msg string `json:"msg"`
}

func (s *WebServer) handlerGet(w http.ResponseWriter, r *http.Request) {
	var req getRequest

	err := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxGetRequestSize)).Decode(&req)
	if err != nil {
		log.Error().Msgf("decode get request: %v", err)
		s.errorHandler(w, r, http.StatusBadRequest, "request decode error")
		return
	}

	if req.Id == "" {
		log.Debug().Msgf("empty request: %v", req.Id)
		s.errorHandler(w, r, http.StatusBadRequest, "incorrect request")
		return
	}

	if len(req.Id) != 24 {
		log.Debug().Msgf("incorrect request id: %v", req.Id)
		s.errorHandler(w, r, http.StatusBadRequest, "incorrect request")
		return
	}

	ttl, err := message_id.GetTtl(req.Id)
	if err != nil {
		if err == message_id.ErrIncorrectMsgId {
			log.Debug().Msgf("incorrect message id lenght: %v", req.Id)
		} else {
			log.Error().Msgf("incorrect message id timestamp: %v", req.Id)
		}

		s.errorHandler(w, r, http.StatusNotFound, "storage error")
		return
	}

	if ttl.Before(time.Now()) {
		log.Debug().Msgf("request expirated id: %v", req.Id)
		s.errorHandler(w, r, http.StatusNotFound)
		return
	}

	msg, err := s.storage.Get(req.Id)
	if err != nil {
		if err == stor.ErrNotFound {
			log.Debug().Msgf("not found id: %v", req.Id)
			s.errorHandler(w, r, http.StatusNotFound)
			return
		}

		log.Error().Msgf("get message: %v", err)
		s.errorHandler(w, r, http.StatusInternalServerError, "get message error")
		return
	}

	err = s.storage.Del(req.Id)
	if err != nil {
		log.Error().Msgf("remove messge: %v", err)
		s.errorHandler(w, r, http.StatusInternalServerError, "storage error")
		return
	}

	resp := getResponse{
		Msg: string(msg),
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Error().Msgf("encode get response: %v", err)
		s.errorHandler(w, r, http.StatusInternalServerError, "encode response error")
		return
	}
}
