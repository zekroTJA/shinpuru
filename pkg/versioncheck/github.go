package versioncheck

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/valyala/fasthttp"
)

const endpoint = "https://api.github.com/repos/%s/%s/tags?per_page=1"

var (
	ErrNoTags = errors.New("no tags available")
)

// GitHubProvider implements Provider using a GitHub
// repositories tag list as source for the current
// version.
type GitHubProvider struct {
	owner, repo string

	lastChecked time.Time
	lastResult  *Semver
}

var _ Provider = (*GitHubProvider)(nil)

// NewGitHubProvider returns a new instance of GitHubProvider
// with the passed repository owner and name as source.
func NewGitHubProvider(owner, repo string) *GitHubProvider {
	return &GitHubProvider{
		owner: owner,
		repo:  repo,
	}
}

type tagsReponse []struct {
	Name string `json:"name"`
}

func (gh *GitHubProvider) GetLatestVersion() (v Semver, err error) {
	if gh.lastResult != nil && time.Since(gh.lastChecked) < 15*time.Minute {
		v = *gh.lastResult
		return
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)

	req.Header.SetMethod("GET")
	req.SetRequestURI(fmt.Sprintf(endpoint, gh.owner, gh.repo))

	if err = fasthttp.Do(req, res); err != nil {
		return
	}

	if res.StatusCode() != 200 {
		err = fmt.Errorf("response error: %d", res.StatusCode())
		return
	}

	var tags tagsReponse
	if err = json.Unmarshal(res.Body(), &tags); err != nil {
		return
	}

	if len(tags) == 0 {
		err = ErrNoTags
		return
	}

	if v, err = ParseSemver(tags[0].Name); err != nil {
		return
	}

	gh.lastChecked = time.Now()
	gh.lastResult = &v

	return
}
