package api

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	vapiBindings_ "github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/lb_services"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"golang.org/x/exp/slices"
)

type LBInfo struct {
	Config         *model.LBService
	Status         *model.LBServiceStatus
	Stats          *model.LBServiceStatistics
	VirtualServers []VSInfo
	Pools          []PoolInfo
}

type VSInfo struct {
	Config *model.LBVirtualServer
	Status model.LBVirtualServerStatus
	Stats  model.LBVirtualServerStatistics
}

type PoolInfo struct {
	Config  *model.LBPool
	Status  model.LBPoolStatus
	Stats   model.LBPoolStatistics
	Members []MemberInfo
}

type MemberInfo struct {
	Status model.LBPoolMemberStatus
	Stats  model.LBPoolMemberStatistics
}

func (a *NSXApi) ListLoadBalancers() ([]model.LBService, error) {
	var cursor *string

	a.log.Debugf("fetching LBService list")
	res := []model.LBService{}
	cli := infra.NewLbServicesClient(a.connector)

	for {
		lbs, err := cli.List(cursor, &False, nil, nil, nil, nil)
		if err != nil {
			a.log.WithError(err).Errorf("could not list LBs")
			return nil, err
		}

		for _, cRes := range lbs.Results {
			empty := len(a.config.LBFilters) == 0
			hasName := slices.Contains(a.config.LBFilters, *cRes.DisplayName)
			hasID := slices.Contains(a.config.LBFilters, *cRes.Id)
			log.Debugf("found LB service '%s' (%s)", *cRes.DisplayName, *cRes.Id)
			if empty || hasName || hasID {
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

func (a *NSXApi) getLBService(lbID string) (*model.LBService, error) {
	a.log.Debugf("fetching LBService '%s'", lbID)

	cli := infra.NewLbServicesClient(a.connector)
	config, err := cli.Get(lbID)
	if err != nil {
		a.log.WithError(err).Errorf("could not fetch LB '%s'", lbID)
		return nil, err
	}

	return &config, nil
}

func (a *NSXApi) getLBServiceStatus(lbID string) (*model.LBServiceStatus, error) {
	a.log.Debugf("fetching status of LBService '%s'", lbID)

	cli := lb_services.NewDetailedStatusClient(a.connector)
	statuses, err := cli.Get(lbID, nil, &False, &RealTime, nil)
	if err != nil {
		a.log.WithError(err).Errorf("could not fetch LB '%s' status", lbID)
		return nil, err
	}

	if len(statuses.Results) == 0 {
		a.log.WithError(fmt.Errorf("not found")).Errorf("could not fetch LB '%s' status", lbID)
		return nil, err
	}

	if len(statuses.Results) > 1 {
		a.log.WithError(fmt.Errorf("too many results")).Errorf("could not fetch LB '%s' status", lbID)
		return nil, err
	}

	t := a.connector.TypeConverter()
	s, errs := t.ConvertToGolang(statuses.Results[0], vapiBindings_.NewReferenceType(model.LBServiceStatusBindingType))
	if len(errs) != 0 {
		a.log.WithError(errs[0]).Errorf("could not fetch LB '%s' status", lbID)
		return nil, err
	}
	val := s.(model.LBServiceStatus)
	return &val, nil
}

func (a *NSXApi) getLBServiceStats(lbID string) (*model.LBServiceStatistics, error) {
	a.log.Debugf("fetching status of LBService '%s'", lbID)

	cli := lb_services.NewStatisticsClient(a.connector)
	stats, err := cli.Get(lbID, nil, &RealTime)
	if err != nil {
		a.log.WithError(err).Errorf("could not fetch LB '%s' status", lbID)
		return nil, err
	}

	if len(stats.Results) == 0 {
		a.log.WithError(fmt.Errorf("not found")).Errorf("could not fetch LB '%s' status", lbID)
		return nil, err
	}

	if len(stats.Results) > 1 {
		a.log.WithError(fmt.Errorf("too many results")).Errorf("could not fetch LB '%s' status", lbID)
		return nil, err
	}

	t := a.connector.TypeConverter()
	s, errs := t.ConvertToGolang(stats.Results[0], vapiBindings_.NewReferenceType(model.LBServiceStatisticsBindingType))
	if len(errs) != 0 {
		a.log.WithError(errs[0]).Errorf("could not fetch LB '%s' status", lbID)
		return nil, err
	}
	val := s.(model.LBServiceStatistics)
	return &val, nil
}

func (a *NSXApi) GetLBServiceInfo(lbID string) (*LBInfo, error) {
	var err error

	res := LBInfo{}
	res.Config, err = a.getLBService(lbID)
	if err != nil {
		return nil, err
	}
	res.Status, err = a.getLBServiceStatus(lbID)
	if err != nil {
		return nil, err
	}
	res.Stats, err = a.getLBServiceStats(lbID)
	if err != nil {
		return nil, err
	}

	for _, cStatus := range res.Status.VirtualServers {
		config, err := a.getVirtualServer(PathToID(*cStatus.VirtualServerPath))
		if err != nil {
			return nil, err
		}
		stats, err := search(func(l model.LBVirtualServerStatistics) string {
			return *l.VirtualServerPath
		}, *cStatus.VirtualServerPath, res.Stats.VirtualServers)
		if err != nil {
			a.log.WithError(err).Errorf("could not associate pool statistics to pool status")
			return nil, err
		}
		res.VirtualServers = append(res.VirtualServers, VSInfo{
			Config: config,
			Status: cStatus,
			Stats:  *stats,
		})
	}

	for _, cStatus := range res.Status.Pools {
		config, err := a.getPool(PathToID(*cStatus.PoolPath))
		if err != nil {
			return nil, err
		}
		stats, err := search(func(l model.LBPoolStatistics) string {
			return *l.PoolPath
		}, *cStatus.PoolPath, res.Stats.Pools)
		if err != nil {
			a.log.WithError(err).Errorf("could not associate pool statistics to pool status")
			return nil, err
		}

		members := []MemberInfo{}
		for _, cMember := range cStatus.Members {
			id := fmt.Sprintf("%s:%s", *cMember.IpAddress, *cMember.Port)
			mStat, err := search(func(m model.LBPoolMemberStatistics) string {
				return fmt.Sprintf("%s:%s", *m.IpAddress, *m.Port)
			}, id, stats.Members)
			if err != nil {
				a.log.WithError(err).Errorf("could not associate member statistics to member status")
				return nil, err
			}
			members = append(members, MemberInfo{
				Stats:  *mStat,
				Status: cMember,
			})
		}

		res.Pools = append(res.Pools, PoolInfo{
			Config:  config,
			Status:  cStatus,
			Stats:   *stats,
			Members: members,
		})
	}

	return &res, nil
}

func (a *NSXApi) getVirtualServer(vsID string) (*model.LBVirtualServer, error) {
	a.log.Debugf("fetching virtual server '%s'", vsID)

	cli := infra.NewLbVirtualServersClient(a.connector)
	config, err := cli.Get(vsID)
	if err != nil {
		a.log.WithError(err).Errorf("could not fetch virtual server '%s'", vsID)
		return nil, err
	}

	return &config, nil
}

func (a *NSXApi) getPool(poolID string) (*model.LBPool, error) {
	a.log.Debugf("fetching pool load balancer '%s'", poolID)

	cli := infra.NewLbPoolsClient(a.connector)
	config, err := cli.Get(poolID)
	if err != nil {
		a.log.WithError(err).Errorf("could not fetch load balancer pool '%s'", poolID)
		return nil, err
	}

	return &config, nil
}
