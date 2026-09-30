package security

import (
	"net/http"
	"os"
)

type SecurityPayload struct {
	ApiKey string
}

func NewSecurityPayload() *SecurityPayload {
	apiKey := os.Getenv("API_KEY")
	return &SecurityPayload{ApiKey: apiKey}
}

func (sp *SecurityPayload) ValidateApiKey(next http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		headerApiKey := req.Header.Get("X-API-KEY")

		if headerApiKey == "" || headerApiKey != sp.ApiKey {
			http.Error(rw, "unauthorized", http.StatusUnauthorized)
			return
		}

		next(rw, req)
	}
}
