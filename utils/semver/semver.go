package semver

import (
	"regexp"
	"strconv"
	"strings"
)

type Version struct {
	Major int
	Minor int
	Patch int
	Type  ReleaseType
	Count int
}

type ReleaseType string

const (
	FULL_RELEASE = ""
	RC_RELEASE = "rc"
	BETA_RELEASE = "beta"
	ALPHA_RELEASE = "alpha"
	DEV_RELEASE = "dev"
)

var semverRegex = regexp.MustCompile(
	`^v?` +								// optional v as prefix
	`(0|[1-9]\d*)\.` +                  // major
	`(0|[1-9]\d*)\.` +                  // minor
	`(0|[1-9]\d*)` +                    // patch
	`(?:-([0-9A-Za-z-]+?)(\d*)` +       // release type + optional numeric suffix
	`(?:\.[0-9A-Za-z-]+)*)?$`,          // allow dots in release type
)

func (t ReleaseType) Long() string {
	switch (t) {
	case RC_RELEASE:
		return "release candidate"
	case DEV_RELEASE:
		return "development"
	case "":
		return "full"
	default:
		return ""
	}
}

func (v Version) String() string {
	res := "v" + strings.Join([]string{strconv.Itoa(v.Major), strconv.Itoa(v.Minor), strconv.Itoa(v.Patch)}, ".")
	
	if v.Type != "" {
		res += "-" + string(v.Type) + strconv.Itoa(v.Count)
	}

	return res
}

func ParseSemver(str string) Version {
	matches := semverRegex.FindStringSubmatch(str)

	if len(matches) == 0 {
		return Version{}
	}

	major, _ := strconv.Atoi(matches[1])
	minor, _ := strconv.Atoi(matches[2])
	patch, _ := strconv.Atoi(matches[3])
	count, _ := strconv.Atoi(matches[5])

	return Version{
		Major: major,
		Minor: minor,
		Patch: patch,
		Type: ParseReleaseType(matches[4]),
		Count: count,
	}
}

func IsValid(str string) bool {
	return semverRegex.MatchString(str)
}

func ParseReleaseType(str string) ReleaseType {
	switch (str) {
	case "rc":
		return RC_RELEASE
	case "beta":
		return BETA_RELEASE
	case "alpha":
		return ALPHA_RELEASE
	case "dev":
		return DEV_RELEASE
	case "":
		return FULL_RELEASE
	default:
		return ""
	}
}