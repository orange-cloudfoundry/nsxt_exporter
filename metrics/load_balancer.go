package metrics

import (
	"github.com/orange-cloudfoundry/nsxt_exporter/api"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// load_balancer_enable{"name", "id"} 0|1
// load_balancer_status{"name", "id"}  1==UP
// load_balancer_info{"name", "id", "size"} 1
// load_balancer_cpu_percent{"name", "id"} %
// load_balancer_errors{"name", "id", "message"} 1
// load_balancer_mem_percent{"name", "id"} %
// load_balancer_alarm{"name", "id", "error_id", "message"} 1
// load_balancer_virtual_server_count{"name", "id"}
// load_balancer_l4_session_rate{"name", "id"}
// load_balancer_l4_session_current{"name", "id"}
// load_balancer_l4_session_total{"name", "id"}
// load_balancer_l4_session_max{"name", "id"}
// load_balancer_l7_session_rate{"name", "id"}
// load_balancer_l7_session_current{"name", "id"}
// load_balancer_l7_session_total{"name", "id"}
// load_balancer_l7_session_max{"name", "id"}

type LBMetrics struct {
	enable    prometheus.GaugeVec
	status    prometheus.GaugeVec
	info      prometheus.GaugeVec
	cpu       prometheus.GaugeVec
	memory    prometheus.GaugeVec
	error     prometheus.GaugeVec
	alarm     prometheus.GaugeVec
	vsCount   prometheus.GaugeVec
	sessionL4 *SessionMetrics
	sessionL7 *SessionMetrics
}

func NewLBMetrics(namespace string) *LBMetrics {
	labels := []string{"name", "id"}
	return &LBMetrics{
		enable: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "load_balancer_enable",
				Help:      "Tells if load balancer is enabled, 1 is enabled",
			}, labels),
		status: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "load_balancer_status",
				Help:      "Gives status of load balancer, 1 is UP",
			}, slice(labels, "status")),
		info: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "load_balancer_info",
				Help:      "Give informations as label about load balancer, value is always 1",
			}, slice(labels, "size")),
		cpu: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "load_balancer_cpu",
				Help:      "CPU usage percentage of load balancer",
			}, labels),
		memory: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "load_balancer_mem",
				Help:      "Memory usage percentage of load balancer",
			}, labels),
		error: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "load_balancer_error",
				Help:      "Current error message for load balancer if any, value is always 1",
			}, slice(labels, "message")),
		alarm: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "load_balancer_alarm",
				Help:      "Give currently firing alarms if any on load balancer, value is always 1",
			}, slice(labels, "error_id", "message")),
		vsCount: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "load_balancer_virtual_server",
				Help:      "Give number of virtual server associated to load balancer",
			}, labels),
		sessionL4: NewSessionMetrics(namespace, "load_balancer", "l4", labels),
		sessionL7: NewSessionMetrics(namespace, "load_balancer", "l7", labels),
	}
}

func (l *LBMetrics) Reset() {
	l.enable.Reset()
	l.status.Reset()
	l.info.Reset()
	l.cpu.Reset()
	l.error.Reset()
	l.memory.Reset()
	l.alarm.Reset()
	l.vsCount.Reset()
	l.sessionL4.Reset()
	l.sessionL7.Reset()
}

func (l *LBMetrics) Populate(name string, id string, info *api.LBInfo) {
	labels := []string{name, id}
	infoLabels := slice(labels, zero(info.Config.Size))

	setb(l.enable, labels, info.Config.Enabled)
	setv(l.status, labels, info.Status.ServiceStatus, StatusUp)
	set(l.info, infoLabels, 1)
	setp(l.cpu, labels, info.Status.CpuUsage)
	setp(l.memory, labels, info.Status.MemoryUsage)
	setl(l.error, labels, info.Status.ErrorMessage)

	if info.Status.Alarm != nil {
		alarmLabels := slice(
			labels,
			zero(info.Status.Alarm.ErrorId),
			zero(info.Status.Alarm.Message),
		)
		set(l.alarm, alarmLabels, 1)
	}

	set(l.vsCount, labels, len(info.Status.VirtualServers))
	setp(l.sessionL4.rate, labels, info.Stats.Statistics.L4CurrentSessionRate)
	setp(l.sessionL4.current, labels, info.Stats.Statistics.L4CurrentSessions)
	setp(l.sessionL4.max, labels, info.Stats.Statistics.L4MaxSessions)
	setp(l.sessionL4.total, labels, info.Stats.Statistics.L4TotalSessions)
	setp(l.sessionL7.rate, labels, info.Stats.Statistics.L7CurrentSessionRate)
	setp(l.sessionL7.current, labels, info.Stats.Statistics.L7CurrentSessions)
	setp(l.sessionL7.max, labels, info.Stats.Statistics.L7MaxSessions)
	setp(l.sessionL7.total, labels, info.Stats.Statistics.L7TotalSessions)
}
