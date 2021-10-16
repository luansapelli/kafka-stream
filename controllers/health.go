package controllers

import (
	"fmt"
	"net/http"
)

// HealthCheck check the health of the service
func HealthCheck() {
	http.HandleFunc("/status-stream/health-check", healthCheckServer)
	_ = http.ListenAndServe(":80", nil)
}

func healthCheckServer(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "alive and kicking")
}
