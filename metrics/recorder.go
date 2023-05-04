package metrics

import (
	"time"

	"github.com/orange-cloudfoundry/nsxt_exporter/api"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Recorder struct {
	manager               *api.NSXApi
	scrapeError           prometheus.Gauge
	scrapeDurationSeconds prometheus.Gauge
	clusterControlStatus  prometheus.Gauge
	clusterMgmtStatus     prometheus.Gauge
	node                  *NodeMetrics
	lb                    *LBMetrics
	vs                    *VSMetrics
	pool                  *PoolMetrics
	tier0                 *Tier0Metrics
	tier1                 *Tier1Metrics
}

func NewRecorder(manager *api.NSXApi, namespace string) *Recorder {
	return &Recorder{
		manager: manager,
		node:    NewNodeMetrics(namespace),
		lb:      NewLBMetrics(namespace),
		vs:      NewVSMetrics(namespace),
		pool:    NewPoolMetrics(namespace),
		tier0:   NewTier0Metrics(namespace),
		tier1:   NewTier1Metrics(namespace),
		scrapeError: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "scrape_error",
				Help:      "last scrape status, 1 when error",
			}),
		scrapeDurationSeconds: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "scrape_duration_seconds",
				Help:      "Duration of Vsphere scraping in milliseconds",
			}),
		clusterControlStatus: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "cluster_control_status",
				Help:      "Cluster control status, 1 means STABLE",
			}),
		clusterMgmtStatus: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "cluster_mgmt_status",
				Help:      "Cluster management status, 1 means STABLE",
			}),
	}
}

func (r *Recorder) Reset() {
	r.node.Reset()
	r.lb.Reset()
	r.vs.Reset()
	r.pool.Reset()
	r.scrapeError.Set(0)
}

func (r *Recorder) RecordMetrics() error {
	r.Reset()

	start := time.Now()

	// cluster
	cluster, err := r.manager.GetClusterStatus()
	if err != nil {
		r.scrapeError.Set(1)
		return err
	}
	r.clusterControlStatus.Set(statusToValue(cluster.ControlClusterStatus.Status, StatusStable))
	r.clusterMgmtStatus.Set(statusToValue(cluster.MgmtClusterStatus.Status, StatusStable))

	// cluster nodes
	nodes := []string{}
	for _, cNode := range cluster.MgmtClusterStatus.OnlineNodes {
		nodes = append(nodes, cNode.Uuid)
	}
	for _, cNode := range cluster.MgmtClusterStatus.OfflineNodes {
		nodes = append(nodes, cNode.Uuid)
	}
	for _, cNodeID := range nodes {
		info, err := r.manager.GetClusterNodeInfo(cNodeID)
		if err != nil {
			r.scrapeError.Set(1)
			break
		}
		err = r.node.Populate(info)
		if err != nil {
			r.scrapeError.Set(1)
		}
	}

	// lb
	lbs, err := r.manager.ListLoadBalancers()
	if err != nil {
		r.scrapeError.Set(1)
		return err
	}

	// load balancer
	for _, cLb := range lbs {
		info, err := r.manager.GetLBServiceInfo(*cLb.Id)
		if err != nil {
			r.scrapeError.Set(1)
			return err
		}
		r.lb.Populate(*cLb.DisplayName, *cLb.Id, info)
		// virtual server
		for _, cVS := range info.VirtualServers {
			r.vs.Populate(cVS)
		}
		// pool
		for _, cPool := range info.Pools {
			r.pool.Populate(cPool)
		}
	}

	t1GWs, err := r.manager.ListT1()
	if err != nil {
		r.scrapeError.Set(1)
		return err
	}
	for _, cT1 := range t1GWs {
		state, err := r.manager.GetT1Status(*cT1.Id)
		if err != nil {
			r.scrapeError.Set(1)
			return err
		}
		r.tier1.Populate(cT1, state.Tier1State, state.Tier1Status)
	}

	t0GWs, err := r.manager.ListT0()
	if err != nil {
		r.scrapeError.Set(1)
		return err
	}
	for _, cT0 := range t0GWs {
		state, err := r.manager.GetT0Status(*cT0.Id)
		if err != nil {
			r.scrapeError.Set(1)
			return err
		}
		r.tier0.Populate(cT0, state.Tier0State, state.Tier0Status)
	}

	r.scrapeDurationSeconds.Set(time.Since(start).Seconds())
	return nil
}
