package root

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/coreos/go-semver/semver"
)

func IsVersionNewerSemver(givenVerStr, currentVerStr string) (bool, error) {
	givenVer, err := semver.NewVersion(givenVerStr)
	if err != nil {
		return false, err
	}

	currentVer, err := semver.NewVersion(currentVerStr)
	if err != nil {
		return false, err
	}

	result := givenVer.Compare(*currentVer)
	if result >= 0 {
		return true, nil
	}

	return false, nil
}

type ReleaseInfo struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
	Body    string `json:"body"`
	URL     string `json:"html_url"`
}

// GetLatestRelease Get latest release info from GitHub
func GetLatestRelease(repo string) (*ReleaseInfo, error) {
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get release info failed: %d", resp.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var releaseInfo ReleaseInfo
	err = json.Unmarshal(bodyBytes, &releaseInfo)
	if err != nil {
		return nil, err
	}

	return &releaseInfo, nil
}

func TrimPrefixV(s string) string {
	if strings.HasPrefix(s, "v") {
		return s[1:]
	}
	return s
}
