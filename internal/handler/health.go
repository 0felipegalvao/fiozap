package handler

import (
	"net/http"
	"runtime"
	"time"

	"fiozap/internal/model"
)

var startTime = time.Now()

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) GetHealth(w http.ResponseWriter, r *http.Request) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	response := map[string]interface{}{
		"status":     "ok",
		"timestamp":  time.Now().UTC().Format(time.RFC3339),
		"uptime":     time.Since(startTime).String(),
		"goroutines": runtime.NumGoroutine(),
		"memory": map[string]interface{}{
			"alloc_mb":       memStats.Alloc / 1024 / 1024,
			"total_alloc_mb": memStats.TotalAlloc / 1024 / 1024,
			"sys_mb":         memStats.Sys / 1024 / 1024,
			"num_gc":         memStats.NumGC,
		},
	}

	model.RespondOK(w, response)
}
