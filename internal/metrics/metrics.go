package metrics

import (
	"net/http"
	"strconv"
	"sync"
	"time"
)

/*
   Prometheus-style in-memory metrics (no external dependency)
   Export via /metrics endpoint
   Counters for HTTP requests, errors, and credits created, approved, and rejected.
   Latency histogram for HTTP requests.
*/

var (
	mu sync.RWMutex

	// Counters
	httpRequestsTotal    map[string]int64
	httpRequestsErrors   map[string]int64
	creditsCreatedTotal  int64
	creditsApprovedTotal int64
	creditsRejectedTotal int64

	// Latency histogram
	httpRequestDuration map[string][]time.Duration
	maxSamples          = 1000
)

func init() {
	httpRequestsTotal = make(map[string]int64)
	httpRequestsErrors = make(map[string]int64)
	httpRequestDuration = make(map[string][]time.Duration)
}

// Increments request count for method_path
func IncHTTPRequest(methodPath string) {
	mu.Lock()
	defer mu.Unlock()
	httpRequestsTotal[methodPath]++
}

// Increments error count
func IncHTTPRequestError(methodPath string) {
	mu.Lock()
	defer mu.Unlock()
	httpRequestsErrors[methodPath]++
}

// Records request duration
func ObserveHTTPDuration(methodPath string, d time.Duration) {
	mu.Lock()
	defer mu.Unlock()
	list := httpRequestDuration[methodPath]

	if len(list) >= maxSamples {
		list = list[1:]
	}

	httpRequestDuration[methodPath] = append(list, d)
}

// Increments credits created
func IncCreditsCreated() {
	mu.Lock()
	defer mu.Unlock()
	creditsCreatedTotal++
}

// Increments credits approved
func IncCreditsApproved() {
	mu.Lock()
	defer mu.Unlock()
	creditsApprovedTotal++
}

// Increments credits rejected
func IncCreditsRejected() {
	mu.Lock()
	defer mu.Unlock()
	creditsRejectedTotal++
}

// Returns a copy of all metrics for exposition
func Snapshot() (total map[string]int64, errors map[string]int64, creditsCreated, creditsApproved, creditsRejected int64, durations map[string][]time.Duration) {
	mu.RLock()
	defer mu.RUnlock()
	total = make(map[string]int64)
	errors = make(map[string]int64)
	durations = make(map[string][]time.Duration)

	for k, v := range httpRequestsTotal {
		total[k] = v
	}

	for k, v := range httpRequestsErrors {
		errors[k] = v
	}

	for k, v := range httpRequestDuration {
		durations[k] = append([]time.Duration(nil), v...)
	}

	return total, errors, creditsCreatedTotal, creditsApprovedTotal, creditsRejectedTotal, durations
}

// Writes Prometheus-style text format for /metrics
func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	total, errors, created, approved, rejected, _ := Snapshot()
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	for k, v := range total {
		_, _ = w.Write([]byte("http_requests_total{path=\"" + k + "\"} " + formatInt64(v) + "\n"))
	}

	for k, v := range errors {
		_, _ = w.Write([]byte("http_requests_errors_total{path=\"" + k + "\"} " + formatInt64(v) + "\n"))
	}

	_, _ = w.Write([]byte("credits_created_total " + formatInt64(created) + "\n"))
	_, _ = w.Write([]byte("credits_approved_total " + formatInt64(approved) + "\n"))
	_, _ = w.Write([]byte("credits_rejected_total " + formatInt64(rejected) + "\n"))
}

func formatInt64(n int64) string {
	return strconv.FormatInt(n, 10)
}
