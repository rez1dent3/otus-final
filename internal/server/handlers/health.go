package handlers

import "net/http"

type Health struct{}

func (h *Health) Handle(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("ok"))
}
