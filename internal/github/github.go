package github

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
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

func getRelease(url string) (Release, error) {
	cmd := exec.Command("gh", "api", url)
	cmd.Stderr = os.Stderr
	response, err := cmd.Output()
	if err != nil {
		return Release{}, err
	}
	release := Release{}
	err = json.Unmarshal(response, &release)
	if err != nil {
		return Release{}, err
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

func SelectAsset(release Release, filters []string) (ReleaseAsset, error) {
	candidates := []ReleaseAsset{}
	res := []*regexp.Regexp{}
	for _, filter := range filters {
		res = append(res, regexp.MustCompile(filter))
	}
	for _, asset := range release.Assets {
		for _, re := range res {
			if re.MatchString(asset.Name) {
				candidates = append(candidates, asset)
			}
		}
	}
	switch len(candidates) {
	case 0:
		return ReleaseAsset{}, fmt.Errorf(
			"no assets found matching the filters '%v'", filters)
	case 1:
		return candidates[0], nil
	default:
		names := []string{}
		for _, cand := range candidates {
			names = append(names, cand.Name)
		}
		return ReleaseAsset{}, fmt.Errorf(
			"multiple assets found match the filter '%v':\n%s",
			filters, strings.Join(names, "\n"))
	}
}
