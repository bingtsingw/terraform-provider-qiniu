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

type CertInfo struct {
	Id         string   `json:"certid"`
	Name       string   `json:"name"`
	CommonName string   `json:"common_name"`
	DnsNames   []string `json:"dnsnames"`
	Pri        string   `json:"pri"`
	Ca         string   `json:"ca"`
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

func (m *CertManager) GetCertsInfo() (certsInfo []CertInfo, err error) {
	type Response struct {
		Marker string     `json:"marker"`
		Certs  []CertInfo `json:"certs"`
	}

	marker := ""

	for {
		response := Response{}
		reqURL := fmt.Sprintf("%s/sslcert?marker=%s&limit=100", ApiHost, marker)
		err = m.Client.CredentialedCall(context.Background(), m.Mac, auth.TokenQiniu, &response, "GET", reqURL, nil)

		if err != nil {
			return nil, err
		}

		if response.Marker == "" {
			break
		}

		marker = response.Marker
		certsInfo = append(certsInfo, response.Certs...)
	}

	return certsInfo, err
}

func (m *CertManager) DeleteCert(id string) (err error) {
	reqURL := fmt.Sprintf("%s/sslcert/%s", ApiHost, string(id))
	err = m.Client.CredentialedCallWithJson(context.Background(), m.Mac, auth.TokenQiniu, nil, "DELETE", reqURL, nil, nil)
	return
}

func (m *CertManager) CreateCert(body CertInfo) (certInfo CertInfo, err error) {
	reqURL := fmt.Sprintf("%s/sslcert", ApiHost)
	err = m.Client.CredentialedCallWithJson(context.Background(), m.Mac, auth.TokenQiniu, &certInfo, "POST", reqURL, nil, body)
	return
}
