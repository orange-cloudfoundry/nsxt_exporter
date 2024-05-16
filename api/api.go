package api

import (
	"crypto/tls"
	"crypto/x509"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/orange-cloudfoundry/nsxt_exporter/config"
	log "github.com/sirupsen/logrus"
	nsxt "github.com/vmware/go-vmware-nsxt"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client/middleware/retry"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/security"
)

var (
	retryCodes = []int{429, 503}
	True       = true
	False      = false
	RealTime   = "realtime"
	Cached     = "cached"
)

type NSXApi struct {
	sync.Mutex
	config    *config.NSXConfig
	connector *client.RestConnector
	client    *nsxt.APIClient
	log       *log.Entry
}

func NewNSXApi(config *config.NSXConfig) (*NSXApi, error) {
	api := &NSXApi{
		config: config,
		log:    log.WithField("module", "api"),
	}

	if err := api.initNSXPolicyConnector(); err != nil {
		api.log.WithError(err).Error("unable to create nsx policy client")
		return nil, err
	}

	retriesConfig := nsxt.ClientRetriesConfiguration{
		MaxRetries:      config.MaxRetries,
		RetryMinDelay:   0,
		RetryMaxDelay:   500,
		RetryOnStatuses: retryCodes,
	}

	host, err := config.NSXHost()
	if err != nil {
		return nil, err
	}

	clientConfig := &nsxt.Configuration{
		BasePath:             "/api/v1",
		Host:                 host,
		Scheme:               "https",
		UserAgent:            "nsxt_exporter",
		UserName:             config.Username,
		Password:             config.Password,
		RemoteAuth:           false,
		ClientAuthCertFile:   config.ClientCertPath,
		ClientAuthKeyFile:    config.ClientKeyPath,
		CAFile:               config.CaCertPath,
		Insecure:             config.SkipSslVerify,
		RetriesConfiguration: retriesConfig,
		SkipSessionAuth:      true,
	}

	api.client, err = nsxt.NewAPIClient(clientConfig)
	if err != nil {
		return nil, err
	}

	return api, nil
}

func (a *NSXApi) initNSXPolicyConnector() error {
	retryFn := a.getNSXPolicyRetryFunc()
	httpClient, err := a.getNSXPolicyHTTPClient()
	if err != nil {
		return err
	}
	a.connector = client.NewRestConnector(
		a.config.URL,
		*httpClient,
		client.WithDecorators(retry.NewRetryDecorator(uint(a.config.MaxRetries), retryFn)),
	)
	a.connector.SetSecurityContext(a.getNSXPolicySecurityContext())
	return nil
}

func (a *NSXApi) getNSXPolicyTLSConfig() (*tls.Config, error) {
	tlsConfig := &tls.Config{
		// nolint:gosec
		InsecureSkipVerify: a.config.SkipSslVerify,
	}

	if len(a.config.ClientCertPath) != 0 {
		cert, err := tls.LoadX509KeyPair(a.config.ClientCertPath, a.config.ClientKeyPath)
		if err != nil {
			a.log.WithError(err).Error("invalid client certificates")
			return nil, err
		}
		tlsConfig.GetClientCertificate = func(*tls.CertificateRequestInfo) (*tls.Certificate, error) {
			return &cert, nil
		}
	}

	if len(a.config.CaCertPath) != 0 {
		caCert, err := os.ReadFile(a.config.CaCertPath)
		if err != nil {
			a.log.WithError(err).Errorf("invalid ca-certificate file '%s'", a.config.CaCertPath)
			return nil, err
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		tlsConfig.RootCAs = caCertPool
	}

	return tlsConfig, nil
}

func (a *NSXApi) getNSXPolicyHTTPClient() (*http.Client, error) {
	tlsConfig, err := a.getNSXPolicyTLSConfig()
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy:           http.ProxyFromEnvironment,
			TLSClientConfig: tlsConfig,
		},
	}
	return client, nil
}

func (a *NSXApi) getNSXPolicyRetryFunc() retry.RetryFunc {
	return func(retryContext retry.RetryContext) bool {
		shouldRetry := false
		if retryContext.Response != nil {
			for _, code := range retryCodes {
				if retryContext.Response.StatusCode == code {
					a.log.Debugf("retrying request due to error code %d", code)
					shouldRetry = true
					break
				}
			}
		} else {
			shouldRetry = true
			a.log.Debugf("retrying request due to error")
		}
		if !shouldRetry {
			return false
		}
		min := 500
		max := 5000
		if max > 0 {
			//nolint:gosec
			interval := rand.Intn(max-min) + min
			time.Sleep(time.Duration(interval) * time.Millisecond)
			a.log.Debugf("waited %d ms before retrying", interval)
		}
		return true
	}
}

func (a *NSXApi) getNSXPolicySecurityContext() core.SecurityContext {
	securityCtx := core.NewSecurityContextImpl()
	if a.config.NeedPasswordLogin() {
		securityCtx.SetProperty(security.AUTHENTICATION_SCHEME_ID, security.USER_PASSWORD_SCHEME_ID)
		securityCtx.SetProperty(security.USER_KEY, a.config.Username)
		securityCtx.SetProperty(security.PASSWORD_KEY, a.config.Password)
	}
	return securityCtx
}
