package domain

import (
	"context"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/client"
)

var (
	ApiHost = "https://api.qiniu.com"
)

type DomainDescriber struct {
	OperationType  string `json:"operationType"`
	OperatingState string `json:"operatingState"`
}

type DomainCacheControl struct {
	Time     int    `json:"time,omitempty"`
	Timeunit int    `json:"timeunit,omitempty"`
	Type     string `json:"type,omitempty"`
	Rule     string `json:"rule,omitempty"`
}

type DomainCacheInfo struct {
	IgnoreParam   bool                 `json:"ignoreParam"`
	CacheControls []DomainCacheControl `json:"cacheControls,omitempty"`
}

type DomainSourceInfo struct {
	Type        string   `json:"sourceType,omitempty"`
	Host        string   `json:"sourceHost,omitempty"`
	IPs         []string `json:"sourceIPs,omitempty"`
	Domain      string   `json:"sourceDomain,omitempty"`
	QiniuBucket string   `json:"sourceQiniuBucket,omitempty"`
	URLScheme   string   `json:"sourceURLScheme,omitempty"`
	TestURLPath string   `json:"testURLPath,omitempty"`
	Advanced    []struct {
		Addr   string `json:"addr,omitempty"`
		Weight int    `json:"weight,omitempty"`
		Backup bool   `json:"backup,omitempty"`
	} `json:"advancedSources,omitempty"`
}

type DomainHttpsInfo struct {
	CertID      string `json:"certid,omitempty"`
	ForceHttps  bool   `json:"forceHttps,omitempty"`
	Http2Enable bool   `json:"http2Enable,omitempty"`
}

type DomainInfo struct {
	Name        string           `json:"name,omitempty"`
	CName       string           `json:"cname,omitempty"`
	Type        string           `json:"type,omitempty"`
	Platform    string           `json:"platform,omitempty"`
	GeoCover    string           `json:"geoCover,omitempty"`
	Protocol    string           `json:"protocol,omitempty"`
	TestURLPath string           `json:"testURLPath,omitempty"`
	Source      DomainSourceInfo `json:"source,omitempty"`
	Https       DomainHttpsInfo  `json:"https,omitempty"`
	Cache       DomainCacheInfo  `json:"cache,omitempty"`
}

type DomainManager struct {
	Client *client.Client
	Mac    *auth.Credentials
}

func NewDomainManager(mac *auth.Credentials) *DomainManager {
	return &DomainManager{
		Client: &client.DefaultClient,
		Mac:    mac,
	}
}

func (m *DomainManager) GetDomainsInfo() (domainInfos []DomainInfo, err error) {
	type Response struct {
		Marker  string       `json:"marker"`
		Domains []DomainInfo `json:"domains"`
	}

	marker := ""

	for {
		response := Response{}
		reqURL := fmt.Sprintf("%s/domain?marker=%s&limit=1000", ApiHost, marker)
		err = m.Client.CredentialedCall(context.Background(), m.Mac, auth.TokenQiniu, &response, "GET", reqURL, nil)

		if err != nil {
			return nil, err
		}

		if response.Marker == "" {
			break
		}

		marker = response.Marker
		domainInfos = append(domainInfos, response.Domains...)
	}

	return domainInfos, err
}

func (m *DomainManager) GetDomainInfo(domain string) (domainInfo DomainInfo, err error) {
	reqURL := fmt.Sprintf("%s/domain/%s", ApiHost, domain)
	err = m.Client.CredentialedCall(context.Background(), m.Mac, auth.TokenQiniu, &domainInfo, "GET", reqURL, nil)
	return domainInfo, err
}

func (m *DomainManager) DescribeDomain(domain string) (response DomainDescriber, err error) {
	reqURL := fmt.Sprintf("%s/domain/%s", ApiHost, domain)
	err = m.Client.CredentialedCall(context.Background(), m.Mac, auth.TokenQiniu, &response, "GET", reqURL, nil)
	return response, err
}

func (m *DomainManager) CreateDomain(domain string, body DomainInfo) (domainInfo DomainInfo, err error) {
	reqURL := fmt.Sprintf("%s/domain/%s", ApiHost, domain)
	err = m.Client.CredentialedCallWithJson(context.Background(), m.Mac, auth.TokenQiniu, &domainInfo, "POST", reqURL, nil, body)
	return domainInfo, err
}

func (m *DomainManager) OfflineDomain(domain string) (err error) {
	reqURL := fmt.Sprintf("%s/domain/%s/offline", ApiHost, domain)
	err = m.Client.CredentialedCall(context.Background(), m.Mac, auth.TokenQiniu, nil, "POST", reqURL, nil)
	return err
}

func (m *DomainManager) DeleteDomain(domain string) (err error) {
	reqURL := fmt.Sprintf("%s/domain/%s", ApiHost, domain)
	err = m.Client.CredentialedCall(context.Background(), m.Mac, auth.TokenQiniu, nil, "DELETE", reqURL, nil)
	return err
}

func (m *DomainManager) UnsslizeDomain(domain string) (err error) {
	reqURL := fmt.Sprintf("%s/domain/%s/unsslize", ApiHost, domain)
	err = m.Client.CredentialedCall(context.Background(), m.Mac, auth.TokenQiniu, nil, "PUT", reqURL, nil)
	return err
}

func (m *DomainManager) SslizeDomain(domain string, body DomainHttpsInfo) (err error) {
	reqURL := fmt.Sprintf("%s/domain/%s/sslize", ApiHost, domain)
	err = m.Client.CredentialedCallWithJson(context.Background(), m.Mac, auth.TokenQiniu, nil, "PUT", reqURL, nil, body)
	return err
}
