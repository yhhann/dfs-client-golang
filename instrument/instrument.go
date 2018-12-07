package instrument

import (
	"flag"
	"fmt"
	"math"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

const (
	WRITE = "Write"
	READ  = "Read"
)

type Measurements struct {
	Name  string
	Value float64
}

var (
	metricsAddr    = flag.String("metrics-address", ":7070", "The address to listen on for metrics.")
	metricsPath    = flag.String("metrics-path", "/dfs-metrics", "The path of metrics.")
	metricsBufSize = flag.Int("metrics-buf-size", 100, "Size of metrics buffer")
)

var (
	RRate = 0.0 // kbit/s
	WRate = 0.0 // kbit/s
)

var (
	// InProcess counts the running routine.
	InProcess chan *Measurements

	// SuccessDuration instruments duration of method called successfully.
	SuccessDuration chan *Measurements

	// FailedCounter instruments number of failed.
	FailedCounter chan *Measurements

	// TimeoutHistogram instruments timeout of method.
	TimeoutHistogram chan *Measurements

	// TransferRate instruments rate of file transfer.
	TransferRate chan *Measurements

	// FileSize instruments size of file.
	FileSize chan *Measurements

	// NoDeadlineCounter instruments number of method which without deadline.
	NoDeadlineCounter chan *Measurements

	// PrejudgeExceed instruments number of method which will exceed deadline.
	PrejudgeExceed chan *Measurements

	inProcessGauge        *prometheus.GaugeVec
	sucDuration           *prometheus.SummaryVec
	failCounter           *prometheus.CounterVec
	timeoutHistogram      *prometheus.HistogramVec
	transferRate          *prometheus.SummaryVec
	fileSize              *prometheus.HistogramVec
	noDeadlineCounter     *prometheus.CounterVec
	prejudgeExceedCounter *prometheus.CounterVec
)

func init() {
	flag.Float64Var(&RRate, "initial-read-rate", 0.0, "Initial transfer rate for reading in kbit/s.")
	flag.Float64Var(&WRate, "initial-write-rate", 0.0, "Initial transfer rate for writing in kbit/s.")

	inProcessGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "dfs2_0",
			Subsystem: "client",
			Name:      "in_process_counter",
			Help:      "Method in process.",
		},
		[]string{"service"},
	)
	InProcess = make(chan *Measurements, *metricsBufSize)

	sucDuration = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: "dfs2_0",
			Subsystem: "client",
			Name:      "suc_durations_nanoseconds",
			Help:      "Successful RPC latency distributions.",
		},
		[]string{"service"},
	)
	SuccessDuration = make(chan *Measurements, *metricsBufSize)

	failCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "dfs2_0",
			Subsystem: "client",
			Name:      "fail_counter",
			Help:      "Failed RPC counter.",
		},
		[]string{"service"},
	)
	FailedCounter = make(chan *Measurements, *metricsBufSize)

	timeoutHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "dfs2_0",
			Subsystem: "client",
			Name:      "timeout_nanoseconds",
			Help:      "timeout distributions.",
			Buckets:   prometheus.ExponentialBuckets(100000, 10, 6),
		},
		[]string{"service"},
	)
	TimeoutHistogram = make(chan *Measurements, *metricsBufSize)

	transferRate = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: "dfs2_0",
			Subsystem: "client",
			Name:      "rate_in_kbit_per_sec",
			Help:      "transfer rate distributions.",
		},
		[]string{"service"},
	)
	TransferRate = make(chan *Measurements, *metricsBufSize)

	fileSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "dfs2_0",
			Subsystem: "client",
			Name:      "size_in_bytes",
			Help:      "file size distributions.",
			Buckets:   prometheus.ExponentialBuckets(100*1024, 2, 6),
		},
		[]string{"service"},
	)
	FileSize = make(chan *Measurements, *metricsBufSize)

	noDeadlineCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "dfs2_0",
			Subsystem: "client",
			Name:      "no_deadline_counter",
			Help:      "no deadline counter.",
		},
		[]string{"service"},
	)
	NoDeadlineCounter = make(chan *Measurements, *metricsBufSize)

	prejudgeExceedCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "dfs2_0",
			Subsystem: "client",
			Name:      "prejudge_exceed_counter",
			Help:      "prejudge exceed counter.",
		},
		[]string{"service"},
	)
	PrejudgeExceed = make(chan *Measurements, *metricsBufSize)

	prometheus.MustRegister(inProcessGauge)
	prometheus.MustRegister(sucDuration)
	prometheus.MustRegister(failCounter)
	prometheus.MustRegister(timeoutHistogram)
	prometheus.MustRegister(transferRate)
	prometheus.MustRegister(fileSize)
	prometheus.MustRegister(noDeadlineCounter)
	prometheus.MustRegister(prejudgeExceedCounter)

	StartMetrics()
}

func StartMetrics() {
	go func() {
		http.Handle(*metricsPath, prometheus.UninstrumentedHandler())
		http.ListenAndServe(*metricsAddr, nil)
	}()

	go func() {
		for {
			select {
			case m := <-InProcess:
				inProcessGauge.WithLabelValues(m.Name).Add(m.Value)
			case m := <-FailedCounter:
				failCounter.WithLabelValues(m.Name).Inc()
			case m := <-NoDeadlineCounter:
				noDeadlineCounter.WithLabelValues(m.Name).Inc()
			case m := <-PrejudgeExceed:
				prejudgeExceedCounter.WithLabelValues(m.Name).Inc()
			case m := <-TimeoutHistogram:
				timeoutHistogram.WithLabelValues(m.Name).Observe(m.Value)
			case m := <-TransferRate:
				transferRate.WithLabelValues(m.Name).Observe(m.Value)
			case m := <-FileSize:
				fileSize.WithLabelValues(m.Name).Observe(m.Value)
			case m := <-SuccessDuration:
				sucDuration.WithLabelValues(m.Name).Observe(m.Value)
			}
		}
	}()
}

func GetTransferRateQuantile(method string, quantile float64) (float64, error) {
	/*
		sum, err := transferRate.GetMetricWith(prometheus.Labels{"service": method})
		if err != nil {
			return 0, err
		}
	*/

	m := &dto.Metric{}
	/*
		if err = sum.Write(m); err != nil {
			return 0, err
		}
	*/

	qs := m.GetSummary().GetQuantile()
	for _, q := range qs {
		if math.Abs(q.GetQuantile()-quantile) < 1e-4 {
			return q.GetValue(), nil
		}
	}

	return 0, fmt.Errorf("Not Found")
}

func startRateCheckRoutine() {
	go func() {
		// refresh rate every 5 seconds.
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				r, err := GetTransferRateQuantile(READ, 0.99)
				if err == nil && r > 0 { // err != nil ignored
					RRate = r
					glog.Infof("Succeeded to update read rate to %f", RRate)
				}

				w, err := GetTransferRateQuantile(WRITE, 0.99)
				if err == nil && w > 0 { // err != nil ignored
					WRate = w
					glog.Infof("Succeeded to update write rate to %f", WRate)
				}
			}
		}
	}()
}
