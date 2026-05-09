package api

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const idempotencyKeyHeader = "Idempotency-Key"

type idempotencyStore struct {
	mu      sync.Mutex
	entries map[string]idempotencyEntry
}

type idempotencyEntry struct {
	Fingerprint string
	Status      int
	Header      http.Header
	Body        []byte
}

func newIdempotencyStore() *idempotencyStore {
	return &idempotencyStore{entries: map[string]idempotencyEntry{}}
}

func (h *Handler) serveWithIdempotency(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.mux.ServeHTTP(w, r)
		return
	}
	key := strings.TrimSpace(r.Header.Get(idempotencyKeyHeader))
	if key == "" {
		h.mux.ServeHTTP(w, r)
		return
	}

	rawBody, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	r.Body = io.NopCloser(bytes.NewReader(rawBody))
	fingerprint := idempotencyFingerprint(r, rawBody)

	if entry, ok := h.idem.get(key); ok {
		if entry.Fingerprint != fingerprint {
			writeStripeError(w, http.StatusConflict, stripeAPIError{
				Type:    stripeErrorIdempotency,
				Message: "Keys for idempotent requests can only be reused with the same method, path, and parameters.",
				Code:    "idempotency_key_in_use",
			})
			return
		}
		writeRecordedResponse(w, entry)
		return
	}

	rec := &recordingResponseWriter{ResponseWriter: w, status: http.StatusOK}
	h.mux.ServeHTTP(rec, r)
	if shouldStoreIdempotentResponse(rec.status) {
		h.idem.set(key, idempotencyEntry{
			Fingerprint: fingerprint,
			Status:      rec.status,
			Header:      rec.Header().Clone(),
			Body:        rec.body.Bytes(),
		})
	}
}

func (s *idempotencyStore) get(key string) (idempotencyEntry, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry, ok := s.entries[key]
	return entry, ok
}

func (s *idempotencyStore) set(key string, entry idempotencyEntry) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.entries[key]; exists {
		return
	}
	s.entries[key] = entry
}

func idempotencyFingerprint(r *http.Request, rawBody []byte) string {
	canonicalBody := canonicalRequestBody(r.Header.Get("Content-Type"), rawBody)
	sum := sha256.Sum256([]byte(strings.Join([]string{
		r.Method,
		r.URL.Path,
		r.URL.RawQuery,
		canonicalBody,
	}, "\n")))
	return hex.EncodeToString(sum[:])
}

func canonicalRequestBody(contentType string, rawBody []byte) string {
	if strings.Contains(contentType, "application/x-www-form-urlencoded") {
		values, err := url.ParseQuery(string(rawBody))
		if err == nil {
			return values.Encode()
		}
	}
	if strings.Contains(contentType, "application/json") {
		var payload any
		if err := json.Unmarshal(rawBody, &payload); err == nil {
			canonical, err := json.Marshal(payload)
			if err == nil {
				return string(canonical)
			}
		}
	}
	return string(rawBody)
}

func shouldStoreIdempotentResponse(status int) bool {
	return status >= http.StatusOK && status < http.StatusBadRequest || status >= http.StatusInternalServerError
}

func writeRecordedResponse(w http.ResponseWriter, entry idempotencyEntry) {
	for key, values := range entry.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(entry.Status)
	_, _ = w.Write(entry.Body)
}

type recordingResponseWriter struct {
	http.ResponseWriter
	status      int
	body        bytes.Buffer
	wroteHeader bool
}

func (w *recordingResponseWriter) WriteHeader(status int) {
	if w.wroteHeader {
		return
	}
	w.wroteHeader = true
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *recordingResponseWriter) Write(body []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(w.status)
	}
	w.body.Write(body)
	return w.ResponseWriter.Write(body)
}
