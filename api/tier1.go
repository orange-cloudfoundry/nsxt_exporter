package api

import (
	"github.com/orange-cloudfoundry/nsxt_exporter/config"
	log "github.com/sirupsen/logrus"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/tier_1s"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func (a *NSXApi) ListT1() ([]model.Tier1, error) {
	var cursor *string

	a.log.Debugf("fetching T1 gateways list")
	res := []model.Tier1{}
	cli := infra.NewTier1sClient(a.connector)

	for {
		lbs, err := cli.List(cursor, &config.False, nil, nil, nil, nil)
		if err != nil {
			a.log.WithError(err).Errorf("could not list T1 gateways")
			return nil, err
		}

		for _, cRes := range lbs.Results {
			empty := len(a.config.T1Filters) == 0
			hasName := containsString(a.config.T1Filters, *cRes.DisplayName)
			hasID := containsString(a.config.T1Filters, *cRes.Id)
			if empty || hasName || hasID {
				log.Debugf("found tier1 gateway '%s' (%s)", *cRes.DisplayName, *cRes.Id)
				res = append(res, cRes)
			}
		}

		cursor = lbs.Cursor
		if cursor == nil {
			break
		}
	}

	return res, nil
}

func (a *NSXApi) GetT1Status(tierID string) (*model.Tier1GatewayState, error) {
	a.log.Debugf("fetching T1 gateway '%s' status", tierID)
	cli := tier_1s.NewStateClient(a.connector)
	statuses, err := cli.Get(tierID, nil, nil, nil, nil, nil, nil, nil, nil)
	if err != nil {
		a.log.WithError(err).Errorf("could not fetch t1 gateway '%s' status", tierID)
		return nil, err
	}
	return &statuses, nil
}
