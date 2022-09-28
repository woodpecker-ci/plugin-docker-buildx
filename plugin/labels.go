package plugin

import (
	"fmt"
	"strings"
	"time"

	"github.com/coreos/go-semver/semver"
)

// Labels returns list of labels to use for image
func (p *Plugin) Labels() []string {
	l := p.settings.Build.Labels.Value()
	// As described in https://github.com/opencontainers/image-spec/blob/main/annotations.md
	l = append(l, fmt.Sprintf("org.opencontainers.image.created=%s", time.Now().UTC().Format(time.RFC3339)))
	if p.settings.Build.Remote != "" {
		l = append(l, fmt.Sprintf("org.opencontainers.image.source=%s", p.settings.Build.Remote))
	}
	if p.pipeline.Repo.Link != "" {
		l = append(l, fmt.Sprintf("org.opencontainers.image.url=%s", p.pipeline.Repo.Link))
	}
	if p.pipeline.Commit.SHA != "" {
		l = append(l, fmt.Sprintf("org.opencontainers.image.revision=%s", p.pipeline.Commit.SHA))
	}
	if p.settings.Build.Ref != "" && strings.HasPrefix(p.settings.Build.Ref, "refs/tags/") {
		v, err := semver.NewVersion(strings.TrimPrefix(p.settings.Build.Ref[10:], "v"))
		if err == nil && v != nil {
			l = append(l, fmt.Sprintf("org.opencontainers.image.version=%s", v.String()))
		}
	}
	return l
}
