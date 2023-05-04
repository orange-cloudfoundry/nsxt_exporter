package api

import (
	log "github.com/sirupsen/logrus"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/tier_0s"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"golang.org/x/exp/slices"
)

func (a *NSXApi) ListT0() ([]model.Tier0, error) {
	var cursor *string

	a.log.Debugf("fetching T0 gateways list")
	res := []model.Tier0{}
	cli := infra.NewTier0sClient(a.connector)

	for {
		lbs, err := cli.List(cursor, &False, nil, nil, nil, nil)
		if err != nil {
			a.log.WithError(err).Errorf("could not list T0 gateways")
			return nil, err
		}

		for _, cRes := range lbs.Results {
			empty := (len(a.config.T0Filters) == 0)
			hasName := slices.Contains(a.config.T0Filters, *cRes.DisplayName)
			hasID := slices.Contains(a.config.T0Filters, *cRes.Id)
			if empty || hasName || hasID {
				log.Debugf("found tier0 gateway '%s' (%s)", *cRes.DisplayName, *cRes.Id)
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

func (a *NSXApi) GetT0Status(tierID string) (*model.Tier0GatewayState, error) {
	a.log.Debugf("fetching T0 gateway '%s' status", tierID)
	cli := tier_0s.NewStateClient(a.connector)
	statuses, err := cli.Get(tierID, nil, nil, nil, nil, nil, nil, nil, nil)
	if err != nil {
		a.log.WithError(err).Errorf("could not fetch t0 gateway '%s' status", tierID)
		return nil, err
	}
	return &statuses, nil
}
