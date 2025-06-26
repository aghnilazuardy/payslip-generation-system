// request_id.go
package middleware

import (
	"net/http"
)

type ctxKey string

const requestIDKey ctxKey = "request_id"

func GetRequestID(r *http.Request) string {
	if v := r.Context().Value(requestIDKey); v != nil {
		if rid, ok := v.(string); ok {
			return rid
		}
	}
	return ""
}
