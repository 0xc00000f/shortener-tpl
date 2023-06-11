package handlers

import (
	"encoding/json"
	"net"
	"net/http"

	"go.uber.org/zap"

	"github.com/0xc00000f/shortener-tpl/internal/shortener"
)

func GetStats(sa *shortener.NaiveShortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !isIPInSubnet(r.Header.Get("X-Real-IP"), sa.TrustedSubnet) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		stats, err := sa.Encoder().GetStats(r.Context())
		if err != nil {
			http.Error(w, "400 page not found", http.StatusBadRequest)
			return
		}

		result, err := json.MarshalIndent(stats, "", " ")
		if err != nil {
			http.Error(w, "400 page not found", http.StatusBadRequest)
			return
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)

		if _, err := w.Write(result); err != nil {
			sa.L.Error("writing body failure", zap.Error(err))
		}
	}
}

func isIPInSubnet(ipAddress, subnet string) bool {
	_, subnetIP, err := net.ParseCIDR(subnet)
	if err != nil {
		return false
	}

	ip := net.ParseIP(ipAddress)
	return subnetIP.Contains(ip)
}
