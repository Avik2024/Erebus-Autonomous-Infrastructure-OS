package main


import (
"context"
"fmt"
"log"
"net/http"
"os"
"os/signal"
"syscall"
"time"


"github.com/Avik2024/erebus/backend/internal/health"
"github.com/Avik2024/erebus/backend/internal/version"
)


func main() {
mux := http.NewServeMux()
mux.HandleFunc("/healthz", health.Handler)
mux.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
fmt.Fprintln(w, version.Version)
})


srv := &http.Server{
Addr: ":8080",
Handler: mux,
}


go func() {
log.Printf("Erebus backend starting on %s", srv.Addr)
if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
log.Fatalf("listen: %s\n", err)
}
}()


// Graceful shutdown
quit := make(chan os.Signal, 1)
signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
<-quit
log.Println("Shutting down server...")


ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
if err := srv.Shutdown(ctx); err != nil {
log.Fatalf("Server forced to shutdown: %v", err)
}


log.Println("Server exiting")
}