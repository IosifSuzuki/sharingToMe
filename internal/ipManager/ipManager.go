package ipManager

import (
	"IosifSuzuki/sharingToMe/internal/configuration"
	"IosifSuzuki/sharingToMe/internal/models"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	ipHost = "api.ipstack.com"
)

func GetIpInfo(ip string) (*models.IpInfo, error) {
	var (
		client = &http.Client{
			Timeout: 10 * time.Second,
		}
		pathURL = url.URL{
			Scheme: "http",
			Host: ipHost,
			Path: ip,
		}
	)
	var pathQuery = pathURL.Query()
	pathQuery.Add("access_key", configuration.Configuration.AppInfo.IpStackAccessKey)
	pathQuery.Add("fields", strings.Join([]string{ "latitude", "longitude", "location.country_flag"}, ","))
	pathURL.RawQuery = pathQuery.Encode()
	//TODO print info
	req, err := http.NewRequest("GET", pathURL.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("content-type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var ipInfo *models.IpInfo
	err = json.Unmarshal(bodyData, &ipInfo)
	return ipInfo, err
}
