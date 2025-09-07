package health


import (
"net/http"
)


// Handler responds with a simple alive message
func Handler(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "text/plain; charset=utf-8")
w.WriteHeader(http.StatusOK)
_, _ = w.Write([]byte("Erebus is alive 🚀"))
}