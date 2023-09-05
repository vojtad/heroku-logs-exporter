package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	herokuLog "heroku-logs-exporter/heroku_log"
)

type HerokuMetric interface {
	HerokuName() string
	Update(value string, labels []string)
	Delete(labels []string)
}

type HerokuMetricGroup interface {
	UpdateFromLog(log *herokuLog.HerokuLog)
}

func updateMetricsFromLog(metrics []HerokuMetric, labels []string, hLog *herokuLog.HerokuLog) {
	for _, metric := range metrics {
		if value, ok := hLog.Value(metric.HerokuName()); ok {
			metric.Update(value, labels)
		}
	}
}

func deleteMetrics(metrics []HerokuMetric, labels []string) {
	for _, metric := range metrics {
		metric.Delete(labels)
	}
}

func updateMetricFromLog(metrics []HerokuMetric, metricHerokuName string, labels []string, value string) {
	for _, metric := range metrics {
		if metric.HerokuName() == metricHerokuName {
			metric.Update(value, labels)
		}
	}
}

type HerokuCounterMetric struct {
	herokuName string
	metric     *prometheus.CounterVec
}

func NewHerokuCounterMetric(herokuName string, prometheusName string, help string, labels []string) *HerokuCounterMetric {
	m := new(HerokuCounterMetric)
	m.herokuName = herokuName
	m.metric = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheusName,
			Help: help,
		},
		labels,
	)

	return m
}

func (m HerokuCounterMetric) HerokuName() string {
	return m.herokuName
}

func (m HerokuCounterMetric) Update(value string, labels []string) {
	m.metric.WithLabelValues(labels...).Inc()
}

func (m HerokuCounterMetric) Delete(labels []string) {
	m.metric.DeleteLabelValues(labels...)
}

type HerokuGaugeMetric struct {
	herokuName string
	metric     *prometheus.GaugeVec
	parser     func(value string) float64
}

func NewHerokuGaugeMetric(herokuName string, prometheusName string, help string, labels []string, parser func(value string) float64) *HerokuGaugeMetric {
	m := new(HerokuGaugeMetric)
	m.herokuName = herokuName
	m.metric = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: prometheusName,
			Help: help,
		},
		labels,
	)

	if parser == nil {
		m.parser = herokuLog.ParseSimpleNumber
	} else {
		m.parser = parser
	}

	return m
}

func (m HerokuGaugeMetric) HerokuName() string {
	return m.herokuName
}

func (m HerokuGaugeMetric) Update(value string, labels []string) {
	m.metric.WithLabelValues(labels...).Set(m.parser(value))
}

func (m HerokuGaugeMetric) Delete(labels []string) {
	m.metric.DeleteLabelValues(labels...)
}

type HerokuSummaryMetric struct {
	herokuName string
	metric     *prometheus.SummaryVec
	parser     func(value string) float64
}

func NewHerokuSummaryMetric(herokuName string, prometheusName string, help string, labels []string, parser func(value string) float64) *HerokuSummaryMetric {
	m := new(HerokuSummaryMetric)
	m.herokuName = herokuName
	m.metric = promauto.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       prometheusName,
			Help:       help,
			Objectives: map[float64]float64{0.01: 0.001, 0.1: 0.01, 0.5: 0.05, 0.9: 0.01, 0.95: 0.001, 0.99: 0.001},
		},
		labels,
	)

	if parser == nil {
		m.parser = herokuLog.ParseSimpleNumber
	} else {
		m.parser = parser
	}

	return m
}

func (m HerokuSummaryMetric) HerokuName() string {
	return m.herokuName
}

func (m HerokuSummaryMetric) Update(value string, labels []string) {
	m.metric.WithLabelValues(labels...).Observe(m.parser(value))
}

func (m HerokuSummaryMetric) Delete(labels []string) {
	m.metric.DeleteLabelValues(labels...)
}

type HerokuHistogramMetric struct {
	herokuName string
	metric     *prometheus.HistogramVec
	parser     func(value string) float64
}

func NewHerokuHistogramMetric(herokuName string, prometheusName string, help string, labels []string, buckets []float64, parser func(value string) float64) *HerokuHistogramMetric {
	if buckets == nil {
		buckets = []float64{.005, .01, .02, 0.04, .06, .08, 0.1, .125, 0.15, 0.175, 0.2, 0.3, 0.4, .5, 1, 2.5, 5, 10, 15, 20}
	}

	m := new(HerokuHistogramMetric)
	m.herokuName = herokuName
	m.metric = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    prometheusName,
			Help:    help,
			Buckets: buckets,
		},
		labels,
	)

	if parser == nil {
		m.parser = herokuLog.ParseSimpleNumber
	} else {
		m.parser = parser
	}

	return m
}

func (m HerokuHistogramMetric) HerokuName() string {
	return m.herokuName
}

func (m HerokuHistogramMetric) Update(value string, labels []string) {
	m.metric.WithLabelValues(labels...).Observe(m.parser(value))
}

func (m HerokuHistogramMetric) Delete(labels []string) {
	m.metric.DeleteLabelValues(labels...)
}
