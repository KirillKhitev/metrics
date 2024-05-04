// Пакет для работы с подсетью.
package subnet

import (
	"net"
	"net/http"

	"github.com/KirillKhitev/metrics/internal/flags"
)

// Middleware проверяет, что запрос от агента пришел из доверенной подсети.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if flags.Args.TrustedSubnet == "" {
			next.ServeHTTP(w, r)
			return
		}

		ipStr := r.Header.Get("X-Real-IP")

		if ipStr == "" {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		ip := net.ParseIP(ipStr)
		if ip == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		_, ipNet, err := net.ParseCIDR(flags.Args.TrustedSubnet)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !ipNet.Contains(ip) {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
