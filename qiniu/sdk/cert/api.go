package cert

import (
	"context"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/client"
)

var (
	ApiHost = "https://api.qiniu.com"
)

type CertBody struct {
	Name string
	Ca   string
	Pri  string
}

type CertInfo struct {
	Id         string   `json:"certid"`
	Name       string   `json:"name"`
	CommonName string   `json:"common_name"`
	DnsNames   []string `json:"dnsnames"`
	Pri        string   `json:"pri"`
	Ca         string   `json:"ca"`
}

type CertsInfo struct {
	Marker string     `json:"marker"`
	Certs  []CertInfo `json:"certs"`
}

type CertManager struct {
	Client *client.Client
	Mac    *auth.Credentials
}

func NewCertManager(mac *auth.Credentials) *CertManager {
	return &CertManager{
		Client: &client.DefaultClient,
		Mac:    mac,
	}
}

func (m *CertManager) GetCertInfo(id string) (certInfo CertInfo, err error) {
	type CertResponse struct {
		Cert CertInfo `json:"cert"`
	}
	certResponse := CertResponse{}
	reqURL := fmt.Sprintf("%s/sslcert/%s", ApiHost, string(id))
	err = m.Client.CredentialedCallWithJson(context.Background(), m.Mac, auth.TokenQiniu, &certResponse, "GET", reqURL, nil, nil)
	certInfo = certResponse.Cert
	return
}

func (m *CertManager) GetCertsInfo() (certsInfo CertsInfo, err error) {
	reqURL := fmt.Sprintf("%s/sslcert?marker=0&limit=100", ApiHost)
	err = m.Client.CredentialedCallWithJson(context.Background(), m.Mac, auth.TokenQiniu, &certsInfo, "GET", reqURL, nil, nil)
	return
}

func (m *CertManager) DeleteCert(id string) (err error) {
	reqURL := fmt.Sprintf("%s/sslcert/%s", ApiHost, string(id))
	err = m.Client.CredentialedCallWithJson(context.Background(), m.Mac, auth.TokenQiniu, nil, "DELETE", reqURL, nil, nil)
	return
}

func (m *CertManager) CreateCert(body CertBody) (certInfo CertInfo, err error) {
	reqURL := fmt.Sprintf("%s/sslcert", ApiHost)
	err = m.Client.CredentialedCallWithJson(context.Background(), m.Mac, auth.TokenQiniu, &certInfo, "POST", reqURL, nil, body)
	return
}
