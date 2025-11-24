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

type assetMatch struct {
	asset      ReleaseAsset
	matchCount int
	rank       int
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
		regexFilters = append(regexFilters, regexp.MustCompile("(?i)"+filter))
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

func matchAnyFilter(assets []ReleaseAsset, regexFilters []*regexp.Regexp) []assetMatch {
	candidates := []assetMatch{}
	for _, asset := range assets {
		matchCount := 0
		rank := -1
		for i, re := range regexFilters {
			if re.MatchString(asset.Name) {
				matchCount++
				if rank == -1 {
					rank = i
				}
			}
		}
		if matchCount > 0 {
			candidates = append(candidates, assetMatch{
				asset:      asset,
				matchCount: matchCount,
				rank:       rank,
			})
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
		optionalMatches := matchAnyFilter(candidates,
			makeRegexFilters(filters.Optional))

		switch len(optionalMatches) {
		case 0:
			return ReleaseAsset{}, fmt.Errorf("multiple assets matched the " +
				"required filters but non matched the optional filters")
		case 1:
			candidate = optionalMatches[0].asset
		default:
			// Find the highest match count
			highestMatchCount := 0
			for _, match := range optionalMatches {
				if match.matchCount > highestMatchCount {
					highestMatchCount = match.matchCount
				}
			}

			// Collect all candidates with the highest match count
			bestCandidates := []assetMatch{}
			for _, match := range optionalMatches {
				if match.matchCount == highestMatchCount {
					bestCandidates = append(bestCandidates, match)
				}
			}

			if len(bestCandidates) == 1 {
				candidate = bestCandidates[0].asset
			} else {
				// Tie-break with ranking
				lowestRank := -1
				for _, match := range bestCandidates {
					if lowestRank == -1 || match.rank < lowestRank {
						lowestRank = match.rank
					}
				}

				tiebreakCandidates := []assetMatch{}
				for _, match := range bestCandidates {
					if match.rank == lowestRank {
						tiebreakCandidates = append(tiebreakCandidates, match)
					}
				}

				if len(tiebreakCandidates) == 1 {
					candidate = tiebreakCandidates[0].asset
				} else {
					names := []string{}
					for _, cand := range tiebreakCandidates {
						names = append(names, cand.asset.Name)
					}
					return ReleaseAsset{}, fmt.Errorf("multiple assets matched both "+
						"the required and optional filters with the same priority:\n  %s\n",
						strings.Join(names, "\n  "))
				}
			}
		}
	}
	log.Infoln("selected release asset:", candidate.Name)
	return candidate, nil
}
