package api

import (
	"github.com/vmware/go-vmware-nsxt/administration"
	"github.com/vmware/go-vmware-nsxt/manager"
)

type NodeInfo struct {
	Interfaces []Interface
	Config     administration.ClusterNodeConfig
	Status     administration.ClusterNodeStatus
}

type Interface struct {
	Config manager.NodeInterfaceProperties
	Stats  manager.NodeInterfaceStatisticsProperties
}

func (a *NSXApi) GetClusterNodeInfo(nodeID string) (*NodeInfo, error) {
	var err error

	a.log.Debugf("fetching cluster node '%s' informations", nodeID)
	res := NodeInfo{}

	// nolint: bodyclose
	res.Status, _, err = a.client.NsxComponentAdministrationApi.ReadClusterNodeStatus(a.client.Context, nodeID, nil)
	if err != nil {
		a.log.WithError(err).Errorf("could not fetch cluster node '%s' status", nodeID)
		return nil, err
	}

	// nolint: bodyclose
	res.Config, _, err = a.client.NsxComponentAdministrationApi.ReadClusterNodeConfig(a.client.Context, nodeID)
	if err != nil {
		a.log.WithError(err).Errorf("could not fetch cluster node '%s' configuration", nodeID)
		return nil, err
	}

	// nolint: bodyclose
	interfaces, _, err := a.client.NsxComponentAdministrationApi.ListClusterNodeInterfaces(a.client.Context, nodeID, nil)
	if err != nil {
		a.log.WithError(err).Errorf("could not fetch cluster node '%s' interfaces", nodeID)
		return nil, err
	}

	for _, cInterface := range interfaces.Results {
		iface := Interface{}
		// nolint: bodyclose
		iface.Config, _, _ = a.client.NsxComponentAdministrationApi.ReadClusterNodeInterface(a.client.Context, nodeID, cInterface.InterfaceId, nil)
		if err != nil {
			a.log.WithError(err).Errorf("could not fetch cluster node '%s' interface '%s' configuration", nodeID, cInterface.InterfaceId)
			return nil, err
		}
		// nolint: bodyclose
		iface.Stats, _, _ = a.client.NsxComponentAdministrationApi.ReadClusterNodeInterfaceStatistics(a.client.Context, nodeID, cInterface.InterfaceId, nil)
		if err != nil {
			a.log.WithError(err).Errorf("could not fetch cluster node '%s' interface '%s' statistics", nodeID, cInterface.InterfaceId)
			return nil, err
		}
		res.Interfaces = append(res.Interfaces, iface)
	}

	return &res, nil
}
