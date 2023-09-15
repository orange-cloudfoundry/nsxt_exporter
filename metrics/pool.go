package metrics

import (
	"fmt"

	"github.com/orange-cloudfoundry/nsxt_exporter/api"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

// pool_info{name, id, port, algorithm}
// pool_member_min{name, id}
// pool_member{name, id}
// pool_alarm{name, id, error_id, message}
// pool_status{name, id} 1 == up

// pool_member_failure{name, id, ip, port}
// pool_member_status{name, id, ip, port} 1 == up

type PoolMetrics struct {
	NetworkMetrics

	info        prometheus.GaugeVec
	alarm       prometheus.GaugeVec
	status      prometheus.GaugeVec
	memberCount prometheus.GaugeVec
	memberMin   prometheus.GaugeVec
	member      *MemberMetrics
}

type MemberMetrics struct {
	NetworkMetrics

	failure prometheus.GaugeVec
	status  prometheus.GaugeVec
}

func NewMemberMetrics(namespace string, labels []string) *MemberMetrics {
	labels = slice(labels, "ip", "port")
	return &MemberMetrics{
		NetworkMetrics: NewNetworkMetrics(namespace, "pool_member", labels),
		status: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "pool_member_status",
				Help:      "Gives status of pool, 1 is UP",
			}, slice(labels, "status")),
		failure: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "pool_member_failure",
				Help:      "Gives failure cause as label if any, value is always 1",
			}, slice(labels, "cause")),
	}
}

func (m *MemberMetrics) Reset() {
	m.NetworkMetrics.Reset()
	m.failure.Reset()
	m.status.Reset()
}

func (m *MemberMetrics) Populate(
	status *model.LBPoolMemberStatus,
	stats *model.LBPoolMemberStatistics,
	labels []string,
) {
	labels = slice(
		labels,
		zero(status.IpAddress),
		zero(status.Port),
	)

	setl(m.failure, labels, status.FailureCause)
	setv(m.status, labels, status.Status, StatusUp)
	m.NetworkMetrics.Populate(labels, stats.Statistics)
}

func NewPoolMetrics(namespace string) *PoolMetrics {
	labels := []string{"name", "id"}
	return &PoolMetrics{
		NetworkMetrics: NewNetworkMetrics(namespace, "pool", labels),
		member:         NewMemberMetrics(namespace, labels),
		status: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "pool_status",
				Help:      "Gives status of pool, 1 is UP",
			}, slice(labels, "status")),
		info: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "pool_info",
				Help:      "Give informations as label about pool, value is always 1",
			}, slice(labels, "port", "algorithm")),
		alarm: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "pool_alarm",
				Help:      "Give currently firing alarms if any on pool, value is always 1",
			}, slice(labels, "error_id", "message")),
		memberCount: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "pool_member",
				Help:      "Current number of member in pool",
			}, slice(labels)),
		memberMin: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "pool_member_min",
				Help:      "Minimum number of member to consider pool active",
			}, slice(labels)),
	}
}

func (p *PoolMetrics) Reset() {
	p.NetworkMetrics.Reset()
	p.info.Reset()
	p.alarm.Reset()
	p.status.Reset()
	p.member.Reset()
	p.memberMin.Reset()
}

func (p *PoolMetrics) Populate(info api.PoolInfo) {
	labels := []string{
		zero(info.Config.DisplayName),
		zero(info.Config.Id),
	}
	infoLabels := slice(
		labels,
		fmt.Sprintf("%d", zero(info.Config.MemberGroup.Port)),
		zero(info.Config.Algorithm),
	)

	p.NetworkMetrics.Populate(labels, info.Stats.Statistics)
	set(p.info, infoLabels, 1)
	setv(p.status, labels, info.Status.Status, StatusUp)

	if info.Status.Alarm != nil {
		alarmLabels := slice(
			labels,
			zero(info.Status.Alarm.ErrorId),
			zero(info.Status.Alarm.Message),
		)
		set(p.alarm, alarmLabels, 1)
	}
	set(p.memberCount, labels, len(info.Status.Members))
	setp(p.memberMin, labels, info.Config.MinActiveMembers)

	for i := range info.Members {
		p.member.Populate(&info.Members[i].Status, &info.Members[i].Stats, labels)
	}
}
