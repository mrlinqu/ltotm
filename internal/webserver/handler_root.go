package webserver

import (
	_ "embed"
	"net/http"
)

//go:embed assets/bundle.html
var rootPageData []byte

func (s *WebServer) handlerRoot(w http.ResponseWriter, r *http.Request) {
	w.Write(rootPageData)
}
