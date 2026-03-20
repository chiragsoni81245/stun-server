package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chiragsoni81245/stun-server/internal/stunserver"
	"github.com/chiragsoni81245/stun-server/pkg/config"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	cfg := config.Load()

	stunserver.InitMetrics()

	server := stunserver.New(cfg)

	if err := server.Start(); err != nil {
		log.Fatalf("failed to start: %v", err)
	}

	// Metrics endpoint
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Println("Metrics on :9090/metrics")
		http.ListenAndServe(":9090", nil)
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	server.Shutdown(ctx)
}
