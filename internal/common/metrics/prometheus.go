package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

type PrometheusMetricsClient struct {
	registry *prometheus.Registry
}

var dynamicCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "dynamic_counter",
		Help: "Count custom keys",
	},
	[]string{"key"},
)

type PrometheusMetricsClientConfig struct {
	Host        string
	ServiceName string
}

func NewPrometheusMetricsClient(config *PrometheusMetricsClientConfig) *PrometheusMetricsClient {
	client := &PrometheusMetricsClient{}
	client.initPrometheus(config)
	return client
}

func (p *PrometheusMetricsClient) initPrometheus(conf *PrometheusMetricsClientConfig) {
	p.registry = prometheus.NewRegistry()
	p.registry.MustRegister(collectors.NewGoCollector(), collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	// custom collectors:
	p.registry.Register(dynamicCounter)

	// metadata wrap
	prometheus.WrapRegistererWith(prometheus.Labels{"serviceName": conf.ServiceName}, p.registry)

	// export
	http.Handle("/metrics", promhttp.HandlerFor(p.registry, promhttp.HandlerOpts{}))
	go func() {
		logrus.Fatalf("failed to start prometheus metrics endpoint, err=%v", http.ListenAndServe(conf.Host, nil))
	}()
}

func (p *PrometheusMetricsClient) Inc(key string, value int) {
	dynamicCounter.WithLabelValues(key).Add(float64(value))
}
