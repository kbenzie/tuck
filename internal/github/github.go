package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"tuck/internal/config"
	"tuck/internal/log"
)

type ReleaseAsset struct {
	Id                 int            `json:"id"`
	Name               string         `json:"name"`
	ContentType        string         `json:"content_type"`
	Size               int            `json:"size"`
	Digest             string         `json:"digest"`
	State              string         `json:"state"`
	Url                string         `json:"url"`
	NodeId             string         `json:"node_id"`
	DownloadCount      int            `json:"download_count"`
	Label              string         `json:"label"`
	Uploader           map[string]any `json:"uploader"`
	BrowserDownloadUrl string         `json:"browser_download_url"`
	CreatedAt          string         `json:"created_at"`
	UpdatedAt          string         `json:"updated_at"`
}

type Release struct {
	Assets          []ReleaseAsset `json:"assets"`
	AssetsUrl       string         `json:"assets_url"`
	Author          map[string]any `json:"author"`
	CreatedAt       string         `json:"created_at"`
	Draft           bool           `json:"draft"`
	HtmlUrl         string         `json:"html_url"`
	Id              int            `json:"id"`
	Name            string         `json:"name"`
	NodeId          string         `json:"node_id"`
	Prerelease      bool           `json:"prerelease"`
	PublishedAt     string         `json:"published_at"`
	TagName         string         `json:"tag_name"`
	TarballUrl      string         `json:"tarball_url"`
	TargetCommitish string         `json:"target_commitish"`
	UploadUrl       string         `json:"upload_url"`
	ZipballUrl      string         `json:"zipball_url"`
}

func isGhLoggedIn() bool {
	cmd := exec.Command("gh", "auth", "status")
	response, err := cmd.Output()
	if err != nil {
		return false
	}
	type AuthStatus struct {
		Hosts map[string]any
	}
	ghAuthStatus := AuthStatus{}
	err = json.Unmarshal(response, &ghAuthStatus)
	if len(ghAuthStatus.Hosts) == 0 {
		return false
	}
	return true
}

func getRelease(url string) (Release, error) {
	release := Release{}
	if isGhLoggedIn() {
		// Prefer using gh when is authenticated
		cmd := exec.Command("gh", "api", url)
		cmd.Stderr = os.Stderr
		response, err := cmd.Output()
		if err != nil {
			return release, err
		}
		err = json.Unmarshal(response, &release)
		if err != nil {
			return release, err
		}
	} else {
		// Use raw http request as gh doesn't allow unauthenticated api
		// requests, however this might get rate limited
		response, err := http.Get(fmt.Sprintf("https://api.github.com/%s", url))
		if err != nil {
			return release, err
		}
		defer response.Body.Close()
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return release, err
		}
		err = json.Unmarshal(body, &release)
		if err != nil {
			return release, err
		}
	}
	return release, nil
}

func GetRelease(repo string, release string) (Release, error) {
	if release == "latest" {
		return getRelease(fmt.Sprintf("repos/%s/releases/latest", repo))
	} else {
		return getRelease(fmt.Sprintf("repos/%s/releases/tags/%s", repo, release))
	}
}

func makeRegexFilters(filters []string) []*regexp.Regexp {
	regexFilters := []*regexp.Regexp{}
	for _, filter := range filters {
		regexFilters = append(regexFilters, regexp.MustCompile(filter))
	}
	return regexFilters
}

func matchAllFilters(assets []ReleaseAsset, regexFilters []*regexp.Regexp) []ReleaseAsset {
	candidates := []ReleaseAsset{}
	// TODO: not sure this is actually matching all the filters...
	for _, asset := range assets {
		matchCount := 0
		for _, re := range regexFilters {
			if re.MatchString(asset.Name) {
				matchCount++
			}
		}
		if matchCount == len(regexFilters) {
			candidates = append(candidates, asset)
		}
	}
	return candidates
}

func matchAnyFilter(assets []ReleaseAsset, regexFilters []*regexp.Regexp) []ReleaseAsset {
	candidates := []ReleaseAsset{}
	for _, asset := range assets {
		matched := false
		for _, re := range regexFilters {
			if re.MatchString(asset.Name) {
				matched = true
			}
		}
		if matched {
			candidates = append(candidates, asset)
		}
	}
	return candidates
}

func SelectAsset(release Release, filters config.ConfigFilters) (ReleaseAsset, error) {
	candidate := ReleaseAsset{}
	candidates := matchAllFilters(release.Assets,
		makeRegexFilters(filters.Required))
	log.Infof("found %d candiates matching required filters:\n", len(candidates))
	for _, cand := range candidates {
		log.Infof("  %s\n", cand.Name)
	}
	switch len(candidates) {
	case 0:
		return ReleaseAsset{}, fmt.Errorf(
			"no assets found matching the filters '%v'", filters)
	case 1:
		candidate = candidates[0]
	default:
		candidates = matchAnyFilter(candidates,
			makeRegexFilters(filters.Optional))
		switch len(candidates) {
		case 1:
			candidate = candidates[0]
		case 0:
			return ReleaseAsset{}, fmt.Errorf("multiple assets matched the " +
				"required filters but non matched the optional filters")
		default:
			names := []string{}
			for _, cand := range candidates {
				names = append(names, cand.Name)
			}
			return ReleaseAsset{}, fmt.Errorf("multiple assets matched both "+
				"the required and optional filters:\n  %s\n",
				strings.Join(names, "\n  "))
		}
	}
	log.Infoln("selected release asset:", candidate.Name)
	return candidate, nil
}
