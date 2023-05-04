package api

import (
	"github.com/vmware/go-vmware-nsxt/administration"
)

func (a *NSXApi) GetClusterStatus() (*administration.ClusterStatus, error) {
	a.log.Debugf("fetching cluster status")
	// nolint: bodyclose
	status, _, err := a.client.NsxComponentAdministrationApi.ReadClusterStatus(a.client.Context, nil)
	if err != nil {
		a.log.WithError(err).Errorf("could not fetch cluster status")
		return nil, err
	}
	return &status, nil
}
