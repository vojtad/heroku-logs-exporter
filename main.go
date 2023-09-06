package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	herokuLog "heroku-logs-exporter/heroku_log"
	"heroku-logs-exporter/metrics"
)

var (
	listenAddress       = flag.String("web.listen-address", ":9841", "Address to listen on for telemetry")
	metricsPath         = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics")
	logsPath            = flag.String("web.logs-path", "/logs", "Path under which to accept Heroku Log Drain")
	logsTokenParamName  = flag.String("web.logs-token-param-name", "token", "Parameter name to check against token parameter value in Heroku Log Drain requests")
	logsTokenParamValue = flag.String("web.logs-token-param-value", "", "Token to check against token parameter in Heroku Log Drain requests")
)

var (
	exportedMetrics = []metrics.HerokuMetricGroup{
		metrics.NewHerokuSystemMetrics(),
		metrics.NewHerokuRuntimeMetrics(),
		metrics.NewHerokuPostgresMetrics(),
		metrics.NewHerokuPgbouncerMetrics(),
		metrics.NewHerokuRouterMetrics(),
		metrics.NewRackTimeoutMetrics(),
	}
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "heroku-logs-exporter")
}

func logsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if *logsTokenParamName != "" && *logsTokenParamValue != "" {
		if *logsTokenParamValue != r.URL.Query().Get(*logsTokenParamName) {
			log.Printf("Token mismatch: %s\n", r.URL)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
	}

	appName := r.URL.Query().Get("app_name")

	count := 0
	scanner := bufio.NewScanner(r.Body)
	for scanner.Scan() {
		hLog := herokuLog.ParseHerokuLog(appName, scanner.Text())

		for _, metric := range exportedMetrics {
			metric.UpdateFromLog(hLog)
		}

		count = count + 1
	}

	log.Printf("Processed %d log lines from %s\n", count, appName)
}

func main() {
	flag.Parse()

	http.HandleFunc("/", helloHandler)
	http.HandleFunc(*logsPath, logsHandler)
	http.Handle(*metricsPath, promhttp.Handler())

	log.Printf("Starting heroku-logs-exporter on %s\n", *listenAddress)

	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
