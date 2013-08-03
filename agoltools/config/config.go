package config

import (
	"agolclient"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
}

func NewDerivedCfg(cfg *Cfg) (derived *DerivedCfg) {
	return &DerivedCfg{
		PortalAPIBaseUrl:            fmt.Sprintf("https://%s/sharing/rest", cfg.PortalDomain),
		PortalOrgAPIBaseUrlTemplate: fmt.Sprintf("https://%s.%s/sharing/rest", OrgKeyTemplate, cfg.PortalOrgDomain),
	}
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
