package models

import "encoding/json"

type IpInfo struct {
	IP 			string 	`json:"ip"`
	Longitude   float32 `json:"longitude"`
	Latitude    float32	`json:"latitude"`
	CountryFlag string	`json:"countryFlag"`
}

func (i *IpInfo) UnmarshalJSON(data []byte) error {
	var ipInfo struct {
		Longitude float32 `json:"longitude"`
		Latitude  float32 `json:"latitude"`
		Location  struct {
			CountryFlag string `json:"country_flag"`
		}
	}
	if err := json.Unmarshal(data, &ipInfo); err != nil {
		return err
	}
	i.Longitude = ipInfo.Longitude
	i.Latitude = ipInfo.Latitude
	i.CountryFlag = ipInfo.Location.CountryFlag
	return nil
}
