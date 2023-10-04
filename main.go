package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type ServiceInfoResponse struct {
	ANInfo    *CfZoneInfo         `json:"AN,omitempty"`
	AtlasInfo *CfZoneInfo         `json:"ATLAS,omitempty"`
	GESANInfo *ServiceInfoWrapper `json:"GES,omitempty"`
}

type ServiceInfoWrapper struct {
	GESProxy *CfZoneInfo `json:"GES_PROXY,omitempty"`
	GESWaf   *CfZoneInfo `json:"GES_WAF,omitempty"`
}

type CfZoneInfo struct {
	CfZoneIPs        []string `json:"ip_addresses,omitempty"`
	CfZoneCName      string   `json:"cname,omitempty"`
	CfZoneRootDomain string   `json:"root_domain,omitempty"`
	CfZoneID         string   `json:"cf_zone_id,omitempty"`
}

type CommaSeparatedRes struct {
	Value string
}

// Render renders a CustomHostnameRecord to the http.ResponseWriter.
func (r *ServiceInfoResponse) Render(w http.ResponseWriter, req *http.Request) error {
	return nil
}

func (r *CommaSeparatedRes) Render(w http.ResponseWriter, req *http.Request) error {
	return nil
}

func main() {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// RESTy routes for "articles" resource
	r.Route("/service_info", func(r chi.Router) {
		r.Get("/", GetAllServiceInfo)
		r.Get("/{zone_name}", GetZoneServiceInfo) // GET /articles/search
	})

	r.Get("/array_test", GetArrayTest)

	fmt.Println("listening.....")
	http.ListenAndServe(":3333", r)

}

func GetZoneServiceInfo(writer http.ResponseWriter, request *http.Request) {
	zoneInfo := createZoneInfo()

	zoneName := chi.URLParam(request, "zone_name")

	switch strings.ToUpper(zoneName) {
	case "GES":
		res := &ServiceInfoResponse{
			GESANInfo: getGESServiceInfo(zoneInfo),
		}
		render.Render(writer, request, res)
	case "AN":
		res := &ServiceInfoResponse{
			ANInfo: getANServiceInfo(zoneInfo),
		}
		render.Render(writer, request, res)
	case "ATLAS":
		res := &ServiceInfoResponse{
			AtlasInfo: getAtlasServiceInfo(zoneInfo),
		}
		render.Render(writer, request, res)
	}
}

func GetArrayTest(writer http.ResponseWriter, request *http.Request) {
	status := request.URL.Query()["status"]

	res := &CommaSeparatedRes{Value: strings.Join(status, ",")}
	render.Render(writer, request, res)
}

func GetAllServiceInfo(writer http.ResponseWriter, request *http.Request) {
	zoneInfo := createZoneInfo()
	res := &ServiceInfoResponse{
		ANInfo:    getANServiceInfo(zoneInfo),
		AtlasInfo: getAtlasServiceInfo(zoneInfo),
		GESANInfo: getGESServiceInfo(zoneInfo),
	}

	render.Render(writer, request, res)
}

func getANServiceInfo(zoneInfo map[string]CfZoneInfo) *CfZoneInfo {
	return getZoneInfo(zoneInfo, "AN")
}

func getAtlasServiceInfo(zoneInfo map[string]CfZoneInfo) *CfZoneInfo {
	return getZoneInfo(zoneInfo, "ATLAS")
}

func getGESServiceInfo(zoneInfo map[string]CfZoneInfo) *ServiceInfoWrapper {
	return &ServiceInfoWrapper{
		GESProxy: getZoneInfo(zoneInfo, "GES_PROXY"),
		GESWaf:   getZoneInfo(zoneInfo, "GES_WAF"),
	}
}

func getZoneInfo(zoneInfo map[string]CfZoneInfo, zoneName string) *CfZoneInfo {
	return &CfZoneInfo{
		CfZoneIPs:        zoneInfo[zoneName].CfZoneIPs,
		CfZoneCName:      zoneInfo[zoneName].CfZoneCName,
		CfZoneRootDomain: zoneInfo[zoneName].CfZoneRootDomain,
		CfZoneID:         zoneInfo[zoneName].CfZoneID,
	}
}

func createZoneInfo() map[string]CfZoneInfo {
	return map[string]CfZoneInfo{
		"AN": {
			CfZoneID:         "11111",
			CfZoneCName:      "wp",
			CfZoneRootDomain: "wpenginepowered.com",
			CfZoneIPs:        createStringArray(strings.Split("1.1.1.1,2.2.2.2", ",")),
		},
		"ATLAS": {
			CfZoneID:         "22222",
			CfZoneCName:      "js.wp",
			CfZoneRootDomain: "wpenginepowered.com",
			CfZoneIPs:        createStringArray(strings.Split("3.3.3.3,4.4.4.4", ",")),
		},
		"GES_WAF": {
			CfZoneID:         "33333",
			CfZoneRootDomain: "wpewaf.com",
			CfZoneIPs:        createStringArray(strings.Split("", ",")),
		},
		"GES_PROXY": {
			CfZoneID:         "44444",
			CfZoneRootDomain: "wpeproxy.com",
			CfZoneIPs:        createStringArray(strings.Split("", ",")),
		},
	}
}

func createStringArray(config interface{}) []string {
	result := config.([]string)
	if len(result) == 0 || result[0] == "" {
		return []string{}
	}

	return result
}
