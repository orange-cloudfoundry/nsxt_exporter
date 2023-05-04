package metrics

import (
	"fmt"

	"github.com/orange-cloudfoundry/nsxt_exporter/api"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type InterfaceMetrics struct {
	info      prometheus.GaugeVec
	rxByte    prometheus.GaugeVec
	rxDropped prometheus.GaugeVec
	rxError   prometheus.GaugeVec
	rxFrame   prometheus.GaugeVec
	rxPacket  prometheus.GaugeVec
	txByte    prometheus.GaugeVec
	txCarrier prometheus.GaugeVec
	txColl    prometheus.GaugeVec
	txDropped prometheus.GaugeVec
	txError   prometheus.GaugeVec
	txPacket  prometheus.GaugeVec
}

type NodeMetrics struct {
	status       prometheus.GaugeVec
	cpu          prometheus.GaugeVec
	fsTotal      prometheus.GaugeVec
	fsUsed       prometheus.GaugeVec
	load1        prometheus.GaugeVec
	load5        prometheus.GaugeVec
	load15       prometheus.GaugeVec
	memTotal     prometheus.GaugeVec
	memUsed      prometheus.GaugeVec
	memCache     prometheus.GaugeVec
	uptime       prometheus.GaugeVec
	version      prometheus.GaugeVec
	certificates prometheus.GaugeVec
	interfaces   *InterfaceMetrics
}

func NewInterfaceMetrics(namespace string) *InterfaceMetrics {
	return &InterfaceMetrics{
		info: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "cluster_node_interface",
				Help:      "Information about cluster node interface, value is always 1",
			}, []string{"uuid", "ip", "name", "dev", "admin", "link", "mtu"}),
		rxByte: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "cluster_node_interface_rx_byte",
				Help:      "Number of bytes received",
			}, []string{"uuid", "ip", "name", "dev"}),
		rxDropped: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "cluster_node_interface_rx_dropped",
				Help:      "Number of packets dropped",
			}, []string{"uuid", "ip", "name", "dev"}),
		rxError: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "cluster_node_interface_rx_error",
				Help:      "Number of receive errors",
			}, []string{"uuid", "ip", "name", "dev"}),
		rxFrame: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "cluster_node_interface_rx_frame",
				Help:      "Number of framing errors",
			}, []string{"uuid", "ip", "name", "dev"}),
		rxPacket: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "cluster_node_interface_rx_packet",
				Help:      "Number of packets received",
			}, []string{"uuid", "ip", "name", "dev"}),
		txByte: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "cluster_node_interface_tx_byte",
				Help:      "Number of bytes transmitted",
			}, []string{"uuid", "ip", "name", "dev"}),
		txCarrier: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "cluster_node_interface_tx_carrier",
				Help:      "Number of carrier losses detected",
			}, []string{"uuid", "ip", "name", "dev"}),
		txColl: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "cluster_node_interface_tx_coll",
				Help:      "Number of collisions detected",
			}, []string{"uuid", "ip", "name", "dev"}),
		txDropped: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "cluster_node_interface_tx_dropped",
				Help:      "Number of packets dropped",
			}, []string{"uuid", "ip", "name", "dev"}),
		txError: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "cluster_node_interface_tx_error",
				Help:      "Number of transmit errors",
			}, []string{"uuid", "ip", "name", "dev"}),
		txPacket: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "cluster_node_interface_tx_packet",
				Help:      "Number of packets transmitted",
			}, []string{"uuid", "ip", "name", "dev"}),
	}
}

func (m *InterfaceMetrics) Reset() {
	m.info.Reset()
	m.rxByte.Reset()
	m.rxDropped.Reset()
	m.rxError.Reset()
	m.rxFrame.Reset()
	m.rxPacket.Reset()
	m.txByte.Reset()
	m.txCarrier.Reset()
	m.txColl.Reset()
	m.txDropped.Reset()
	m.txError.Reset()
	m.txPacket.Reset()
}

func NewNodeMetrics(namespace string) *NodeMetrics {
	return &NodeMetrics{
		interfaces: NewInterfaceMetrics(namespace),
		status: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "cluster_node_status",
				Help:      "Cluster node status, 1 means CONNECTED",
			}, []string{"uuid", "ip", "name"}),
		cpu: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "cluster_node_cpu",
				Help:      "Number of CPU core",
			}, []string{"uuid", "ip", "name"}),
		fsTotal: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "cluster_node_fs_total",
				Help:      "Total filesystem space in kB",
			}, []string{"uuid", "ip", "name", "type", "mount"}),
		fsUsed: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "cluster_node_fs_used",
				Help:      "Used filesystem space in kB",
			}, []string{"uuid", "ip", "name", "type", "mount"}),
		load1: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "cluster_node_load1",
				Help:      "Current load average (load 1 minute)",
			}, []string{"uuid", "ip", "name"}),
		load5: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "cluster_node_load5",
				Help:      "Current load average (load 5 minutes)",
			}, []string{"uuid", "ip", "name"}),
		load15: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "cluster_node_load15",
				Help:      "Current load average (load 15 minutes)",
			}, []string{"uuid", "ip", "name"}),
		memTotal: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "cluster_node_mem_total",
				Help:      "Total available memory in kB",
			}, []string{"uuid", "ip", "name"}),
		memUsed: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "cluster_node_mem_used",
				Help:      "Used memory in kB",
			}, []string{"uuid", "ip", "name"}),
		memCache: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "cluster_node_mem_cache",
				Help:      "Cached memory in kB",
			}, []string{"uuid", "ip", "name"}),
		uptime: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "cluster_node_uptime",
				Help:      "Uptime expressed in millisecond since start",
			}, []string{"uuid", "ip", "name"}),
		version: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "cluster_node_version",
				Help:      "Node current version, value always 1",
			}, []string{"uuid", "ip", "name", "version"}),
		certificates: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "cluster_node_certificates",
				Help:      "Node SSL certificate validity end date expressed in number of second since EPOCH",
			}, []string{"uuid", "ip", "name", "type", "index"}),
	}
}

func (m *NodeMetrics) Reset() {
	m.status.Reset()
	m.cpu.Reset()
	m.fsTotal.Reset()
	m.fsUsed.Reset()
	m.load1.Reset()
	m.load5.Reset()
	m.load15.Reset()
	m.memTotal.Reset()
	m.memUsed.Reset()
	m.memCache.Reset()
	m.uptime.Reset()
	m.version.Reset()
	m.certificates.Reset()
	m.interfaces.Reset()
}

func (m *NodeMetrics) Populate(info *api.NodeInfo) error {
	labels := []string{
		info.Config.Id,
		info.Config.ApplianceMgmtListenAddr,
		info.Config.DisplayName,
	}

	m.status.WithLabelValues(labels...).Set(statusToValue(info.Status.MgmtClusterStatus.MgmtClusterStatus, StatusConnected))
	m.cpu.WithLabelValues(labels...).Set(float64(info.Status.SystemStatus.CpuCores))
	for _, cFS := range info.Status.SystemStatus.FileSystems {
		fsLabels := slice(labels, cFS.Type_, cFS.FileSystem)
		m.fsTotal.WithLabelValues(fsLabels...).Set(float64(cFS.Total))
		m.fsUsed.WithLabelValues(fsLabels...).Set(float64(cFS.Used))
	}

	if len(info.Status.SystemStatus.LoadAverage) == 3 {
		m.load1.WithLabelValues(labels...).Set(float64(info.Status.SystemStatus.LoadAverage[0]))
		m.load5.WithLabelValues(labels...).Set(float64(info.Status.SystemStatus.LoadAverage[1]))
		m.load15.WithLabelValues(labels...).Set(float64(info.Status.SystemStatus.LoadAverage[2]))
	}

	m.uptime.WithLabelValues(labels...).Set(float64(info.Status.SystemStatus.Uptime))

	versionLabels := slice(labels, info.Status.Version)
	m.version.WithLabelValues(versionLabels...).Set(float64(1))

	certs, err := processCertificates(info.Config.ManagerRole.ApiListenAddr.Certificate)
	if err != nil {
		return err
	}
	for _, cCert := range certs {
		certLabels := slice(labels, "api", fmt.Sprintf("%d", cCert.index))
		m.certificates.WithLabelValues(certLabels...).Set(float64(cCert.notAfter.Unix()))
	}

	certs, err = processCertificates(info.Config.ManagerRole.MgmtClusterListenAddr.Certificate)
	if err != nil {
		return err
	}
	for _, cCert := range certs {
		certLabels := slice(labels, "mgmt_cluster", fmt.Sprintf("%d", cCert.index))
		m.certificates.WithLabelValues(certLabels...).Set(float64(cCert.notAfter.Unix()))
	}

	certs, err = processCertificates(info.Config.ManagerRole.MgmtPlaneListenAddr.Certificate)
	if err != nil {
		return err
	}
	for _, cCert := range certs {
		certLabels := slice(labels, "mgmt_plane", fmt.Sprintf("%d", cCert.index))
		m.certificates.WithLabelValues(certLabels...).Set(float64(cCert.notAfter.Unix()))
	}

	for _, cIface := range info.Interfaces {
		iFaceLabels := slice(labels, cIface.Config.InterfaceId)
		infoLabels := slice(iFaceLabels, cIface.Config.AdminStatus, cIface.Config.LinkStatus, fmt.Sprintf("%d", cIface.Config.Mtu))
		m.interfaces.info.WithLabelValues(infoLabels...).Set(1)
		m.interfaces.rxByte.WithLabelValues(iFaceLabels...).Set(float64(cIface.Stats.RxBytes))
		m.interfaces.rxDropped.WithLabelValues(iFaceLabels...).Set(float64(cIface.Stats.RxDropped))
		m.interfaces.rxError.WithLabelValues(iFaceLabels...).Set(float64(cIface.Stats.RxErrors))
		m.interfaces.rxFrame.WithLabelValues(iFaceLabels...).Set(float64(cIface.Stats.RxFrame))
		m.interfaces.rxPacket.WithLabelValues(iFaceLabels...).Set(float64(cIface.Stats.RxPackets))
		m.interfaces.txByte.WithLabelValues(iFaceLabels...).Set(float64(cIface.Stats.TxBytes))
		m.interfaces.txCarrier.WithLabelValues(iFaceLabels...).Set(float64(cIface.Stats.TxCarrier))
		m.interfaces.txColl.WithLabelValues(iFaceLabels...).Set(float64(cIface.Stats.TxColls))
		m.interfaces.txDropped.WithLabelValues(iFaceLabels...).Set(float64(cIface.Stats.TxDropped))
		m.interfaces.txError.WithLabelValues(iFaceLabels...).Set(float64(cIface.Stats.TxErrors))
		m.interfaces.txPacket.WithLabelValues(iFaceLabels...).Set(float64(cIface.Stats.TxPackets))
	}

	return nil
}
