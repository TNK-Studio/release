package release

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// Check for an update. Takes in the current version and GitHub repo URL.
// Returns true or false if there is an update or not as well as the version
// value. Will return false if there is no network connection.
func Check(localVersion string, repoURL string) (bool, string, error) {
	requestURL := convertURL(repoURL)
	hasConnection := checkConnection(requestURL)
	if !hasConnection {
		return false, "", nil
	}
	currentVersion, err := getVersion(requestURL)
	if err != nil {
		return false, "", err
	}
	if localVersion != currentVersion {
		return true, currentVersion, nil
	}
	return false, currentVersion, nil
}

// Check for a network connection
func checkConnection(uri string) bool {
	resp, err := http.Get(uri)
	if err != nil || resp.StatusCode != 200 {
		return false
	}
	return true
}

// Convert repo url to api url
// From: https://github.com/Matt-Gleich/nuke
// To:   https://api.github.com/repos/Matt-Gleich/nuke/releases/latest
func convertURL(repoURL string) string {
	var fixedURL string
	fixedURL = strings.Replace(repoURL, "https://github.com/", "https://api.github.com/repos/", 1)
	if fixedURL[len(fixedURL)-1:] == "/" {
		fixedURL = fixedURL + "releases/latest"
	} else {
		fixedURL = fixedURL + "/releases/latest"
	}
	return fixedURL
}

// Make the actual get request to get the version
func getVersion(requestURL string) (string, error) {
	resp, err := http.Get(requestURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return "", err
	}
	version := fmt.Sprintf("%v", data["tag_name"])
	if version == "" {
		return "", errors.New("Version number for repo is blank")
	}
	if version == "<nil>" {
		return "", errors.New("Latest release not found for given repo URL")
	}
	return version, nil
}
