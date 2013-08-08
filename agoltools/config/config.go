package config

import (
	"agolclient"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

func init() {
	Config, _ = fromFile("agoltools/config/config.json")
	Config.DerivedCfg = *NewDerivedCfg(Config)
	agolclient.SetConfig(&agolclient.Config{PortalAPIBaseUrl: Config.PortalAPIBaseUrl})
}

const (
	OrgKeyTemplate = "{$orgKey}"
)

var Config *Cfg

type Cfg struct {
	DerivedCfg
	PortalDomain    string
	PortalOrgDomain string
	AppId           string
	AppBaseUrl      string
}

type DerivedCfg struct {
	PortalAPIBaseUrl            string
	PortalOrgAPIBaseUrlTemplate string
	PortalHomeUrl               string
	PortalOrgHomeUrlTemplate    string
}

func NewDerivedCfg(cfg *Cfg) (derived *DerivedCfg) {
	return &DerivedCfg{
		PortalAPIBaseUrl:            fmt.Sprintf("https://%s/sharing/rest", cfg.PortalDomain),
		PortalOrgAPIBaseUrlTemplate: fmt.Sprintf("https://%s.%s/sharing/rest", OrgKeyTemplate, cfg.PortalOrgDomain),
		PortalHomeUrl:               fmt.Sprintf("https://%s/home", cfg.PortalDomain),
		PortalOrgHomeUrlTemplate:    fmt.Sprintf("https://%s.%s/home", OrgKeyTemplate, cfg.PortalOrgDomain),
	}
}

func PortalHomeUrl(orgUrlKey string) string {
	if orgUrlKey == "" {
		return Config.PortalHomeUrl
	}
	return strings.Replace(Config.PortalOrgHomeUrlTemplate, OrgKeyTemplate, orgUrlKey, 1)
}

func fromFile(filePath string) (cfg *Cfg, err error) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, &cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
