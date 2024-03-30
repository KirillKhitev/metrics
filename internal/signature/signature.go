package signature

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"

	"github.com/KirillKhitev/metrics/internal/flags"
)

type signatureWriter struct {
	w http.ResponseWriter
}

func newSignatureWriter(w http.ResponseWriter) *signatureWriter {
	return &signatureWriter{
		w: w,
	}
}

func (s *signatureWriter) Header() http.Header {
	return s.w.Header()
}

func (s *signatureWriter) Write(p []byte) (int, error) {
	hashSum := GetHash(p, flags.Args.Key)

	s.w.Header().Set("HashSHA256", hashSum)

	return s.w.Write(p)
}

func (s *signatureWriter) WriteHeader(statusCode int) {
	s.w.WriteHeader(statusCode)
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if flags.Args.Key == "" {
			next.ServeHTTP(w, r)
			return
		}

		sw := newSignatureWriter(w)

		body, _ := io.ReadAll(r.Body)
		r.Body.Close()
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		hashSum := GetHash(body, flags.Args.Key)
		headerHash := r.Header.Get("HashSHA256")

		if headerHash != "" && headerHash != hashSum {
			sw.WriteHeader(http.StatusBadRequest)
			return
		}

		next.ServeHTTP(sw, r)
	})
}

func GetHash(data []byte, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(data)
	result := h.Sum(nil)

	return hex.EncodeToString(result)
}
