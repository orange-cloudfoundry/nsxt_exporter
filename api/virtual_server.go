package api

// import (
// // log "github.com/sirupsen/logrus"
// // "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
// // model "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
// // "golang.org/x/exp/slices"
// )

// func (a *NSXApi) listVirtualServers() ([]model.LBVirtualServer, error) {
// 	var cursor *string

// 	a.log.Debugf("listing virtual servers")
// 	res := []model.LBVirtualServer{}
// 	cli := infra.NewLbVirtualServersClient(a.connector)

// 	for {
// 		lbs, err := cli.List(cursor, &False, nil, nil, nil, nil)
// 		if err != nil {
// 			a.log.WithError(err).Errorf("could not list virtual servers")
// 			return nil, err
// 		}

// 		for _, cRes := range lbs.Results {
// 			empty := (len(a.config.VSFilters) == 0)
// 			hasName := slices.Contains(a.config.VSFilters, *cRes.DisplayName)
// 			hasID := slices.Contains(a.config.VSFilters, *cRes.Id)
// 			log.Debugf("found virtual server '%s' (%s)", *cRes.DisplayName, *cRes.Id)
// 			if empty || hasName || hasID {
// 				res = append(res, cRes)
// 			}
// 		}

// 		cursor = lbs.Cursor
// 		if cursor == nil {
// 			break
// 		}
// 	}
// 	return res, nil
// }
