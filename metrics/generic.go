package metrics

import (
	"fmt"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type SessionMetrics struct {
	rate    prometheus.GaugeVec
	current prometheus.GaugeVec
	total   prometheus.GaugeVec
	max     prometheus.GaugeVec
}

type NetworkMetrics struct {
	http         *RateMetrics
	inPacket     *RateMetrics
	outPacket    *RateMetrics
	inByte       *RateMetrics
	outByte      *RateMetrics
	session      *RateMetrics
	sessionTotal *TotalMetrics
}

type TotalMetrics struct {
	total prometheus.GaugeVec
	max   prometheus.GaugeVec
}

type RateMetrics struct {
	current prometheus.GaugeVec
	rate    prometheus.GaugeVec
}

func NewTotalMetrics(namespace string, object string, kind string, labels []string) *TotalMetrics {
	return &TotalMetrics{
		total: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      fmt.Sprintf("%s_%s_total", object, kind),
				Help:      fmt.Sprintf("Total number of %s in %s", kind, strings.ReplaceAll(object, "_", " ")),
			}, labels),
		max: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      fmt.Sprintf("%s_%s_max", object, kind),
				Help:      fmt.Sprintf("Maximum number of %s in %s", kind, strings.ReplaceAll(object, "_", " ")),
			}, labels),
	}
}

func (v *TotalMetrics) Reset() {
	v.total.Reset()
	v.max.Reset()
}

func NewRateMetrics(namespace string, object string, kind string, labels []string) *RateMetrics {
	return &RateMetrics{
		current: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      fmt.Sprintf("%s_%s", object, kind),
				Help:      fmt.Sprintf("Current number of %s in %s", kind, strings.ReplaceAll(object, "_", " ")),
			}, labels),
		rate: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      fmt.Sprintf("%s_%s_rate", object, kind),
				Help:      fmt.Sprintf("Number of %s per second in %s", kind, strings.ReplaceAll(object, "_", " ")),
			}, labels),
	}
}

func (v *RateMetrics) Reset() {
	v.current.Reset()
	v.rate.Reset()
}

func NewSessionMetrics(namespace string, object string, kind string, labels []string) *SessionMetrics {
	return &SessionMetrics{
		rate: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      fmt.Sprintf("%s_session_%s_rate", object, kind),
				Help:      fmt.Sprintf("Number of new %s session per second for %s", kind, object),
			}, labels),
		current: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      fmt.Sprintf("%s_session_%s_current", object, kind),
				Help:      fmt.Sprintf("Current number of %s session for %s", kind, object),
			}, labels),
		total: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      fmt.Sprintf("%s_session_%s_total", object, kind),
				Help:      fmt.Sprintf("Total number of %s session for %s", kind, object),
			}, labels),
		max: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      fmt.Sprintf("%s_session_%s_max", object, kind),
				Help:      fmt.Sprintf("Maximum number of %s session for %s", kind, object),
			}, labels),
	}
}

func (s *SessionMetrics) Reset() {
	s.rate.Reset()
	s.current.Reset()
	s.total.Reset()
	s.max.Reset()
}

func NewNetworkMetrics(namespace string, object string, labels []string) NetworkMetrics {
	return NetworkMetrics{
		http:         NewRateMetrics(namespace, object, "http_request", labels),
		inByte:       NewRateMetrics(namespace, object, "in_byte", labels),
		inPacket:     NewRateMetrics(namespace, object, "in_packet", labels),
		outByte:      NewRateMetrics(namespace, object, "out_byte", labels),
		outPacket:    NewRateMetrics(namespace, object, "out_packet", labels),
		session:      NewRateMetrics(namespace, object, "session", labels),
		sessionTotal: NewTotalMetrics(namespace, object, "session", labels),
	}
}

func (n *NetworkMetrics) Reset() {
	n.http.Reset()
	n.inByte.Reset()
	n.inPacket.Reset()
	n.outByte.Reset()
	n.outPacket.Reset()
	n.session.Reset()
	n.sessionTotal.Reset()
}

func (n *NetworkMetrics) Populate(labels []string, c *model.LBStatisticsCounter) {
	setp(n.http.current, labels, c.HttpRequests)
	setp(n.http.rate, labels, c.HttpRequestRate)
	setp(n.inPacket.current, labels, c.PacketsIn)
	setp(n.inPacket.rate, labels, c.PacketsInRate)
	setp(n.outPacket.current, labels, c.PacketsOut)
	setp(n.outPacket.rate, labels, c.PacketsOutRate)
	setp(n.inByte.current, labels, c.BytesIn)
	setp(n.inByte.rate, labels, c.BytesInRate)
	setp(n.outByte.current, labels, c.BytesOut)
	setp(n.outByte.rate, labels, c.BytesOutRate)
	setp(n.session.current, labels, c.CurrentSessions)
	setp(n.session.rate, labels, c.CurrentSessionRate)
	setp(n.sessionTotal.max, labels, c.MaxSessions)
	setp(n.sessionTotal.total, labels, c.TotalSessions)
}
