package plugin

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/6543/go-version"
)

const tagRefPrefix = "refs/tags/"

var dateRegex = regexp.MustCompile(`^(\d{8}|\d{4}-\d{2}-\d{2})$`)

// DefaultTagSuffix returns a set of default suggested tags
// based on the commit ref with an attached suffix.
func DefaultTagSuffix(ref, defaultTag, suffix string) ([]string, error) {
	tags, err := DefaultTags(ref, defaultTag)
	if err != nil {
		return nil, err
	}
	if len(suffix) == 0 {
		return tags, nil
	}
	for i, tag := range tags {
		if tag == defaultTag {
			tags[i] = suffix
		} else {
			tags[i] = fmt.Sprintf("%s-%s", tag, suffix)
		}
	}
	return tags, nil
}

func splitOff(input, delim string) string {
	parts := strings.SplitN(input, delim, 2)

	if len(parts) == 2 {
		return parts[0]
	}

	return input
}

// DefaultTags returns a set of default suggested tags based on
// the commit ref.
func DefaultTags(ref, defaultTag string) ([]string, error) {
	// check if no tag event
	if !strings.HasPrefix(ref, tagRefPrefix) {
		return []string{defaultTag}, nil
	}

	// else it's an tag event
	tagString := stripTagPrefix(ref)

	// check if date
	if dateRegex.MatchString(tagString) {
		return []string{tagString}, nil
	}

	version, err := version.NewSemver(tagString)
	// if no semversion return default tag and error
	if err != nil {
		return []string{defaultTag}, err
	}

	vParts := version.Segments()
	major, minor, patch := vParts[0], vParts[1], vParts[2]

	// if prerelease or version with metadata, only use this strict version
	if version.Prerelease() != "" || version.Metadata() != "" {
		return []string{
			version.String(),
		}, nil
	}

	// check if version is acutaly a date (%Y%m%d) ... and only return that if so
	if major > 999 && major < 10000 && minor == 0 && patch == 0 {
		return []string{
			fmt.Sprintf("%d", major),
		}, nil
	}

	if major == 0 {
		return []string{
			fmt.Sprintf("%d.%d", major, minor),
			fmt.Sprintf("%d.%d.%d", major, minor, patch),
		}, nil
	}
	return []string{
		fmt.Sprintf("%d", major),
		fmt.Sprintf("%d.%d", major, minor),
		fmt.Sprintf("%d.%d.%d", major, minor, patch),
	}, nil
}

// UseDefaultTag for keep only default branch for latest tag
// return true if tag event or default branch
func UseDefaultTag(ref, defaultBranch string) bool {
	return strings.HasPrefix(ref, tagRefPrefix) ||
		stripHeadPrefix(ref) == defaultBranch
}

func stripHeadPrefix(ref string) string {
	return strings.TrimPrefix(ref, "refs/heads/")
}

func stripTagPrefix(ref string) string {
	return strings.TrimPrefix(ref, tagRefPrefix)
}

func isSingleTag(tag string) bool {
	// currently only filtering for separators, this could be improved...
	return tag == "" || (!regexp.MustCompile(`[,\s]+`).MatchString(tag) && len(tag) <= 128)
}
