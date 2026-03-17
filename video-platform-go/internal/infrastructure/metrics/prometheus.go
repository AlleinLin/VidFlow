package metrics

import (
	"runtime"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	HTTPRequestsInFlight = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Number of HTTP requests currently being processed",
		},
		[]string{"method"},
	)

	VideoUploadsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "video_uploads_total",
			Help: "Total number of video uploads",
		},
	)

	VideoTranscodeDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "video_transcode_duration_seconds",
			Help:    "Duration of video transcoding in seconds",
			Buckets: []float64{10, 30, 60, 120, 300, 600, 1200, 1800, 3600},
		},
		[]string{"resolution"},
	)

	VideoTranscodeQueueSize = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "video_transcode_queue_size",
			Help: "Number of videos waiting to be transcoded",
		},
	)

	VideoViewsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "video_views_total",
			Help: "Total number of video views",
		},
		[]string{"video_id"},
	)

	CommentsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "comments_total",
			Help: "Total number of comments",
		},
	)

	LikesTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "likes_total",
			Help: "Total number of likes",
		},
	)

	DanmakusTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "danmakus_total",
			Help: "Total number of danmakus",
		},
	)

	DatabaseConnections = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "database_connections",
			Help: "Number of database connections",
		},
		[]string{"state"},
	)

	DatabaseQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "database_query_duration_seconds",
			Help:    "Duration of database queries in seconds",
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"query"},
	)

	CacheHits = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total number of cache hits",
		},
	)

	CacheMisses = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "Total number of cache misses",
		},
	)

	CacheOperations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_operations_total",
			Help: "Total number of cache operations",
		},
		[]string{"operation"},
	)

	KafkaMessagesProduced = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kafka_messages_produced_total",
			Help: "Total number of Kafka messages produced",
		},
		[]string{"topic"},
	)

	KafkaMessagesConsumed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kafka_messages_consumed_total",
			Help: "Total number of Kafka messages consumed",
		},
		[]string{"topic", "group"},
	)

	KafkaConsumerLag = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "kafka_consumer_lag",
			Help: "Kafka consumer lag in messages",
		},
		[]string{"topic", "group"},
	)

	UserRegistrations = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "user_registrations_total",
			Help: "Total number of user registrations",
		},
	)

	UserLogins = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "user_logins_total",
			Help: "Total number of user logins",
		},
	)

	ActiveUsers = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_users",
			Help: "Number of active users",
		},
	)

	RecommendationRequests = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "recommendation_requests_total",
			Help: "Total number of recommendation requests",
		},
	)

	RecommendationLatency = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "recommendation_latency_seconds",
			Help:    "Latency of recommendation generation in seconds",
			Buckets: []float64{0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5},
		},
	)

	SearchRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "search_requests_total",
			Help: "Total number of search requests",
		},
		[]string{"type"},
	)

	SearchLatency = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "search_latency_seconds",
			Help:    "Latency of search queries in seconds",
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1},
		},
	)

	StorageOperations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "storage_operations_total",
			Help: "Total number of storage operations",
		},
		[]string{"operation", "type"},
	)

	StorageBytes = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "storage_bytes_total",
			Help: "Total bytes transferred in storage operations",
		},
		[]string{"operation", "type"},
	)

	Goroutines = promauto.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "goroutines_count",
			Help: "Number of goroutines",
		},
		func() float64 {
			return float64(getGoroutines())
		},
	)

	MemoryAlloc = promauto.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "memory_alloc_bytes",
			Help: "Bytes of allocated heap objects",
		},
		func() float64 {
			return float64(getMemoryAlloc())
		},
	)

	MemorySys = promauto.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "memory_sys_bytes",
			Help: "Total bytes of memory obtained from the OS",
		},
		func() float64 {
			return float64(getMemorySys())
		},
	)
)

func getGoroutines() int {
	return runtime.NumGoroutine()
}

func getMemoryAlloc() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.Alloc
}

func getMemorySys() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.Sys
}

func RecordHTTPRequest(method, path, status string, duration float64) {
	HTTPRequestsTotal.WithLabelValues(method, path, status).Inc()
	HTTPRequestDuration.WithLabelValues(method, path).Observe(duration)
}

func RecordDatabaseQuery(query string, duration float64) {
	DatabaseQueryDuration.WithLabelValues(query).Observe(duration)
}

func RecordCacheHit() {
	CacheHits.Inc()
}

func RecordCacheMiss() {
	CacheMisses.Inc()
}

func RecordCacheOperation(operation string) {
	CacheOperations.WithLabelValues(operation).Inc()
}

func RecordVideoUpload() {
	VideoUploadsTotal.Inc()
}

func RecordVideoView(videoID int64) {
	VideoViewsTotal.WithLabelValues(string(rune(videoID))).Inc()
}

func RecordComment() {
	CommentsTotal.Inc()
}

func RecordLike() {
	LikesTotal.Inc()
}

func RecordDanmaku() {
	DanmakusTotal.Inc()
}

func RecordUserRegistration() {
	UserRegistrations.Inc()
}

func RecordUserLogin() {
	UserLogins.Inc()
}

func RecordSearchRequest(searchType string) {
	SearchRequests.WithLabelValues(searchType).Inc()
}

func RecordSearchLatency(duration float64) {
	SearchLatency.Observe(duration)
}

func RecordRecommendationRequest() {
	RecommendationRequests.Inc()
}

func RecordRecommendationLatency(duration float64) {
	RecommendationLatency.Observe(duration)
}

func SetDatabaseConnections(idle, open, max int) {
	DatabaseConnections.WithLabelValues("idle").Set(float64(idle))
	DatabaseConnections.WithLabelValues("open").Set(float64(open))
	DatabaseConnections.WithLabelValues("max").Set(float64(max))
}

func SetTranscodeQueueSize(size int) {
	VideoTranscodeQueueSize.Set(float64(size))
}

func SetActiveUsers(count int) {
	ActiveUsers.Set(float64(count))
}
