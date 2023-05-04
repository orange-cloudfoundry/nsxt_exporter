package metrics

import (
	"fmt"

	"github.com/orange-cloudfoundry/nsxt_exporter/api"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

// tier0_info{"name", "id", "mode"} 1
// tier0_status{"name", "id"} 1 == in_sync
// tier0_failure{"name", "id", "code", "message"} 1
// tier0_transport{"name", "id"} len(Details)
// tier0_edge{"name", "id", "index", "mode"} 1

type TierMetrics struct {
	kind      string
	info      prometheus.GaugeVec
	status    prometheus.GaugeVec
	failure   prometheus.GaugeVec
	transport prometheus.GaugeVec
	edge      prometheus.GaugeVec
}

type Tier1Metrics struct {
	TierMetrics
}

type Tier0Metrics struct {
	TierMetrics
}

func NewTier0Metrics(namespace string) *Tier0Metrics {
	return &Tier0Metrics{
		TierMetrics: NewTierMetrics(namespace, "tier0"),
	}
}

func NewTier1Metrics(namespace string) *Tier1Metrics {
	return &Tier1Metrics{
		TierMetrics: NewTierMetrics(namespace, "tier1"),
	}
}

func NewTierMetrics(namespace string, kind string) TierMetrics {
	labels := []string{"id", "name"}
	return TierMetrics{
		kind: kind,
		info: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      fmt.Sprintf("%s_info", kind),
				Help:      fmt.Sprintf("Give informations as label about %s, value is always 1", kind),
			}, slice(labels, "mode")),
		status: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      fmt.Sprintf("%s_status", kind),
				Help:      fmt.Sprintf("Give status of %s, 1 is is_sync", kind),
			}, labels),
		failure: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      fmt.Sprintf("%s_failure", kind),
				Help:      fmt.Sprintf("Give failure details for %s if any, value is always 1", kind),
			}, slice(labels, "code", "message")),
		transport: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      fmt.Sprintf("%s_transport", kind),
				Help:      fmt.Sprintf("Number of transport node associated to %s ", kind),
			}, labels),
		edge: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      fmt.Sprintf("%s_edge", kind),
				Help:      fmt.Sprintf("Current mode as label of edge node associated to %s, value is always 1", kind),
			}, slice(labels, "index", "mode")),
	}
}

func (t *TierMetrics) Reset() {
	t.info.Reset()
	t.status.Reset()
	t.failure.Reset()
	t.transport.Reset()
	t.edge.Reset()
}

func (t *TierMetrics) Populate(
	labels []string,
	mode string,
	state *model.LogicalRouterState,
	status *model.LogicalRouterStatus,
) {
	set(t.info, slice(labels, mode), 1)
	setv(t.status, labels, state.State, StatusInSync)

	if state.FailureCode != nil {
		failureLabels := slice(
			labels,
			fmt.Sprintf("%d", zero(state.FailureCode)),
			zero(state.FailureMessage),
		)
		set(t.failure, failureLabels, 1)
	}
	set(t.transport, labels, len(state.Details))
	for _, cNode := range status.PerNodeStatus {
		ID := api.PathToID(zero(cNode.EdgePath))
		HA := zero(cNode.HighAvailabilityStatus)
		set(t.edge, slice(labels, ID, HA), 1)
	}
}

func (t *Tier0Metrics) Populate(
	config model.Tier0,
	state *model.LogicalRouterState,
	status *model.LogicalRouterStatus,
) {
	labels := []string{
		zero(config.Id),
		zero(config.DisplayName),
	}
	mode := zero(config.HaMode)
	t.TierMetrics.Populate(labels, mode, state, status)
}

func (t *Tier1Metrics) Populate(
	config model.Tier1,
	state *model.LogicalRouterState,
	status *model.LogicalRouterStatus,
) {
	labels := []string{
		zero(config.Id),
		zero(config.DisplayName),
	}
	mode := zero(config.HaMode)
	t.TierMetrics.Populate(labels, mode, state, status)
}
