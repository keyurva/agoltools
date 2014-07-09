package orgusers

import (
	"agolclient"
	"agoltools"
	"agoltools/auth"
	"strings"
	"time"
)

func init() {
	agoltools.HandleFunc("/registeredappsloginstats", auth.Authenticated(getRegisteredAppsLoginStats))
}

const (
	getRegisteredAppsLoginStatsTemplate = "agoltools/registeredappsloginstats/templates/getregisteredappsloginstats.html"
	registeredAppsLoginStatsTemplate    = "agoltools/registeredappsloginstats/templates/registeredappsloginstats.html"
)

func getRegisteredAppsLoginStats(r *agoltools.Request) (err error) {
	r.Data["PageTitle"] = "My Registered Apps Login Stats"

	f := strings.ToLower(strings.Trim(r.R.FormValue("f"), " "))
	if f == "" {
		return r.RenderUsingBaseTemplate(getRegisteredAppsLoginStatsTemplate)
	}

	period := strings.ToLower(strings.Trim(r.R.FormValue("period"), " "))
	req := createRegisteredAppLoginStatsReq(period)

	ris, err := agolclient.GetUserRegisteredAppsLoginStats(r.Transport(), r.Auth, req)
	if err != nil {
		return err
	}

	if f == "csv" {
		r.W.Header().Set("Content-Type", "text/csv")
		r.W.Header().Set("Content-Disposition", "inline;filename=registeredappsloginstats.csv")
		agolclient.RegisteredAppItemsLoginStatsCsv(r.W, ris, r.PortalHomeUrl())
		return
	}

	r.AddData(map[string]interface{}{
		"RegisteredAppLoginStatsItems": ris,
	})

	return r.RenderUsingBaseTemplate(registeredAppsLoginStatsTemplate)
}

func createRegisteredAppLoginStatsReq(period string) *agolclient.RegisteredAppLoginStatsReq {
	now := time.Now().UTC()
	midnight := now.Add(time.Duration(24-now.Hour()) * time.Hour).Truncate(time.Hour)
	req := &agolclient.RegisteredAppLoginStatsReq{
		EndTime: midnight,
	}

	if period == "6m" {
		req.StartTime = midnight.Add(-180 * 24 * time.Hour)
		req.Period = "180d"
	} else if period == "m" {
		req.StartTime = midnight.Add(-30 * 24 * time.Hour)
		req.Period = "30d"
	} else if period == "w" {
		req.StartTime = midnight.Add(-7 * 24 * time.Hour)
		req.Period = "1w"
	} else {
		req.StartTime = midnight.Add(-24 * time.Hour)
		req.Period = "24h"
	}

	return req
}
