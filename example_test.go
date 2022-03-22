package gin_monitor_test

import (
	"log"
	"net/http"
	"testing"
	"time"

	ginMonitor "github.com/bancodobrasil/gin-monitor"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type FakeDependencyChecker struct{}

func (m *FakeDependencyChecker) GetDependencyName() string {
	return "fake-dependency"
}

func (m *FakeDependencyChecker) Check() ginMonitor.DependencyStatus {
	return ginMonitor.DOWN
}

func YourHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("gin-monitor!\n"))
}

func TestMainHandler(t *testing.T) {
	// Creates gin-monitor instance
	monitor, err := ginMonitor.New("v1.0.0", ginMonitor.DefaultErrorMessageKey, ginMonitor.DefaultBuckets)
	if err != nil {
		panic(err)
	}

	dependencyChecker := &FakeDependencyChecker{}
	monitor.AddDependencyChecker(dependencyChecker, time.Second*30)

	r := gin.New()

	// Register gin-monitor middleware
	r.Use(monitor.Prometheus())
	// Register metrics endpoint
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	// Routes consist of a path and a handler function.
	r.GET("/", gin.WrapF(YourHandler))

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8000", r))
}
