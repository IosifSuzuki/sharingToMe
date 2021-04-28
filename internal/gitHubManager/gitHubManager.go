package gitHubManager

import (
	"IosifSuzuki/sharingToMe/internal/models"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	gitHubHost = "api.github.com"
)

func FetchCommitInfos() ([]models.CommitInfo, error) {
	var (
		client = &http.Client{
			Timeout: 10 * time.Second,
		}
		pathURL = url.URL{
			Scheme: "http",
			Host: gitHubHost,
			Path: "repos/IosifSuzuki/sharingToMe/commits",
		}
	)
	var pathQuery = pathURL.Query()
	pathQuery.Add("sha", "development")
	pathURL.RawQuery = pathQuery.Encode()

	req, err := http.NewRequest("GET", pathURL.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("content-type", "application/json")
	resp, err := client.Do(req)
	defer resp.Body.Close()
	bodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var commitInfos []models.CommitInfo
	err = json.Unmarshal(bodyData, &commitInfos)
	return commitInfos, nil
}
