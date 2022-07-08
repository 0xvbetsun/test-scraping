package scraping

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

type Handler struct {
	Limit    uint32
	connects uint32
	mu       sync.Mutex
}

func NewHandler(limit uint32) *Handler {
	if limit == 0 || limit >= 999 {
		limit = 999
	}
	return &Handler{Limit: limit}
}

func (h *Handler) Allow() bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.Limit > h.connects
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST, OPTIONS")
		msg := fmt.Sprintf("Method %s is not allowed.", r.Method)
		http.Error(w, msg, http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Content-Type") != "text/plain" {
		msg := "Content-Type header is not text/plain"
		http.Error(w, msg, http.StatusUnsupportedMediaType)
		return
	}
	if !h.Allow() {
		http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
		return
	}
	h.inc()
	defer h.dec()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	w.Write([]byte(body))
}

func (h *Handler) inc() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.connects++
}

func (h *Handler) dec() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.connects--
}
