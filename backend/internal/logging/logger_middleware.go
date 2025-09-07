package logging

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

// List of sensitive keys to mask
var sensitiveKeys = []string{"password", "token", "authorization", "secret", "apikey", "api_key", "cookie"}

// LoggerMiddleware logs structured request/response details with masking
func LoggerMiddleware(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			reqID := middleware.GetReqID(r.Context())
			clientIP := r.RemoteAddr
			method := r.Method
			path := r.URL.Path

			// Mask query parameters
			maskedQuery := maskQueryParams(r.URL.Query())

			// Mask headers
			maskedHeaders := maskHeaders(r.Header)

			// Read request body
			var reqBody []byte
			if r.Body != nil && (method == http.MethodPost || method == http.MethodPut) {
				bodyBytes, _ := io.ReadAll(r.Body)
				reqBody = maskSensitiveRecursive(bodyBytes)
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // reset body
			}

			// Capture response
			rec := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK, body: &bytes.Buffer{}}
			next.ServeHTTP(rec, r)

			latency := time.Since(start)
			resBody := maskSensitiveRecursive(rec.body.Bytes())

			logger.Info("HTTP request",
				zap.String("request_id", reqID),
				zap.String("method", method),
				zap.String("path", path),
				zap.Reflect("query_params", maskedQuery),
				zap.Reflect("headers", maskedHeaders),
				zap.String("client_ip", clientIP),
				zap.Int("status", rec.statusCode),
				zap.Duration("latency", latency),
				zap.ByteString("request_body", reqBody),
				zap.ByteString("response_body", resBody),
			)
		})
	}
}

// Mask sensitive query parameters
func maskQueryParams(values url.Values) map[string][]string {
	masked := make(map[string][]string)
	for key, vals := range values {
		if isSensitiveKey(key) {
			masked[key] = []string{"*****"}
		} else {
			masked[key] = vals
		}
	}
	return masked
}

// Mask sensitive headers
func maskHeaders(headers http.Header) map[string][]string {
	masked := make(map[string][]string)
	for key, vals := range headers {
		if isSensitiveKey(strings.ToLower(key)) {
			masked[key] = []string{"*****"}
		} else {
			masked[key] = vals
		}
	}
	return masked
}

// Mask sensitive fields recursively in JSON
func maskSensitiveRecursive(data []byte) []byte {
	if len(data) == 0 {
		return data
	}

	var obj interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		// Not JSON, return as-is
		return data
	}

	maskObject(obj)

	masked, err := json.Marshal(obj)
	if err != nil {
		return data
	}
	return masked
}

// maskObject recursively masks sensitive fields in maps and slices
func maskObject(obj interface{}) {
	switch v := obj.(type) {
	case map[string]interface{}:
		for key, val := range v {
			if isSensitiveKey(strings.ToLower(key)) {
				v[key] = "*****"
			} else {
				maskObject(val)
			}
		}
	case []interface{}:
		for i := range v {
			maskObject(v[i])
		}
	}
}

func isSensitiveKey(key string) bool {
	for _, s := range sensitiveKeys {
		if key == s {
			return true
		}
	}
	return false
}

// responseRecorder wraps http.ResponseWriter to capture status code and body
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (r *responseRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}
