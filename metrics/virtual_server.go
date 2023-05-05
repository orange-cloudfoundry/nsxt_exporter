package metrics

import (
	"github.com/orange-cloudfoundry/nsxt_exporter/api"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// virtual_server_enable{name, id}
// virtual_server_status{name, id} 1 == UP
// virtual_server_info{name, id, ip, pool_id, lb_id}
// virtual_server_cnx_max{name, id}
// virtual_server_cnx_max_rate{name, id}
// virtual_server_alarm{name, id}
// virtual_server_session{name, id}
// virtual_server_session_rate{name, id}
// virtual_server_http_request{name, id}
// virtual_server_http_request_rate{name, id}
// virtual_server_in_packet{name, id}
// virtual_server_in_packet_rate{name, id}
// virtual_server_out_packet{name, id}
// virtual_server_out_packet_rate{name, id}
// virtual_server_in_byte{name, id}
// virtual_server_in_byte_rate{name, id}
// virtual_server_out_bytes{name, id}
// virtual_server_out_bytes_rate{name, id}
// virtual_server_session_max{name, id}
// virtual_server_session_total{name, id}
// virtual_server_source_ip{name, id}

type VSMetrics struct {
	NetworkMetrics

	enable prometheus.GaugeVec
	status prometheus.GaugeVec
	info   prometheus.GaugeVec
	alarm  prometheus.GaugeVec
	ip     prometheus.GaugeVec
}

func NewVSMetrics(namespace string) *VSMetrics {
	labels := []string{"name", "id"}
	return &VSMetrics{
		NetworkMetrics: NewNetworkMetrics(namespace, "virtual_server", labels),
		enable: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "virtual_server_enable",
				Help:      "Tells if virtual server is enabled, 1 is enabled",
			}, labels),
		status: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "virtual_server_status",
				Help:      "Gives status of virtual server, 1 is UP",
			}, slice(labels, "status")),
		info: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "virtual_server_info",
				Help:      "Give informations as label about virtual server, value is always 1",
			}, slice(labels, "ip", "pool_id", "lb_id")),
		alarm: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "virtual_server_alarm",
				Help:      "Give currently firing alarms if any on virtual server, value is always 1",
			}, slice(labels, "error_id", "message")),
		ip: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "virtual_server_source_ip",
				Help:      "Number of source IP persistence entries in virtual server",
			}, labels),
	}
}

func (v *VSMetrics) Reset() {
	v.NetworkMetrics.Reset()
	v.enable.Reset()
	v.status.Reset()
	v.info.Reset()
	v.alarm.Reset()
	v.ip.Reset()
	v.http.Reset()
}

func (v *VSMetrics) Populate(info api.VSInfo) {
	labels := []string{
		zero(info.Config.DisplayName),
		zero(info.Config.Id),
	}

	infoLabels := slice(
		labels,
		zero(info.Config.IpAddress),
		api.PathToID(zero(info.Config.PoolPath)),
		api.PathToID(zero(info.Config.LbServicePath)),
	)

	if info.Status.Alarm != nil {
		alarmLabels := slice(
			labels,
			zero(info.Status.Alarm.ErrorId),
			zero(info.Status.Alarm.Message),
		)
		set(v.alarm, alarmLabels, 1)
	}

	setb(v.enable, labels, info.Config.Enabled)
	setv(v.status, labels, info.Status.Status, StatusUp)
	set(v.info, infoLabels, 1)
	setp(v.ip, labels, info.Stats.Statistics.SourceIpPersistenceEntrySize)
	v.NetworkMetrics.Populate(labels, info.Stats.Statistics)
}
