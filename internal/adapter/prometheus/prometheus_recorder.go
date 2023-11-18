package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type PrometheusRecorder struct {
	cacheClearCounter      prometheus.Counter
	mediaRequestCounter    prometheus.Counter
	bandwithUsageCounter   prometheus.Counter
	dataStorageUsageGauge  prometheus.Gauge
	cacheStorageUsageGauge prometheus.Gauge
}

func NewPrometheusRecorder() *PrometheusRecorder {
	cacheClearCounter := promauto.NewCounter(prometheus.CounterOpts{
		Name: "mindia_cache_cleared_processed_ops_total",
		Help: "The total number of processed cache clear",
	})

	mediaRequestCounter := promauto.NewCounter(prometheus.CounterOpts{
		Name: "mindia_media_request_processed_ops_total",
		Help: "The total number of processed media requests",
	})

	bandwithUsageCounter := promauto.NewCounter(prometheus.CounterOpts{
		Name: "mindia_media_bandwidth_usage_total",
		Help: "The total number of processed media requests",
	})

	dataStorageUsageGauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "mindia_data_storage_usage_total",
		Help: "The total storage usage ",
	})

	cacheStorageUsageGauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "mindia_cache_storage_usage_total",
		Help: "The total storage usage ",
	})

	return &PrometheusRecorder{
		cacheClearCounter:      cacheClearCounter,
		mediaRequestCounter:    mediaRequestCounter,
		bandwithUsageCounter:   bandwithUsageCounter,
		dataStorageUsageGauge:  dataStorageUsageGauge,
		cacheStorageUsageGauge: cacheStorageUsageGauge,
	}
}

func (r *PrometheusRecorder) RecordBandwithUsage(bytesLength int) {
	r.bandwithUsageCounter.Add(float64(bytesLength))
}

func (r *PrometheusRecorder) RecordMediaRequest() {
	r.mediaRequestCounter.Inc()
}

func (r *PrometheusRecorder) RecordMediaDelete() {

}

func (r *PrometheusRecorder) RecordDataStorageUsage(bytesUsage int64) {
	r.dataStorageUsageGauge.Set(float64(bytesUsage))
}

func (r *PrometheusRecorder) RecordCacheStorageUsage(bytesUsage int64) {
	r.cacheStorageUsageGauge.Set(float64(bytesUsage))
}

func (r *PrometheusRecorder) RecordCacheClear() {
	r.cacheClearCounter.Inc()
}

func (r *PrometheusRecorder) RecordTaskCreation() {

}
