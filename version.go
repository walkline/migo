package migo

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	version "github.com/blang/semver"
)

type Version struct {
	v         *version.Version
	timestamp int64
	Name      string
}

func (v Version) String() string {
	skipSegs := true
	if v.v.Patch != 0 || v.v.Minor != 0 || v.timestamp != 0 {
		skipSegs = false
	}

	ver := ""
	if skipSegs {
		ver = fmt.Sprintf("%v", v.v.Major)
	} else {
		ver = v.v.String()
		if v.timestamp > 0 {
			ver += fmt.Sprintf(".%d", v.timestamp)
		}
	}

	return fmt.Sprintf("%s-%s", ver, v.Name)
}

func VersionFromString(s string) (*Version, error) {
	strs := strings.Split(s, "-")
	if len(strs) < 2 {
		return nil, errors.New("bad version format")
	}

	semver := strs[0]
	var timestamp int64
	versions := strings.Split(strs[0], ".")
	if len(versions) == 4 {
		semver = fmt.Sprintf("%s.%s.%s", versions[0], versions[1], versions[2])
		var err error
		timestamp, err = strconv.ParseInt(versions[3], 10, 64)
		if err != nil {
			return nil, err
		}
	}

	v, err := version.ParseTolerant(semver)
	if err != nil {
		return nil, err
	}

	strAsByte := []byte(s)
	name := string(strAsByte[len(strs[0])+1:])

	return &Version{
		v:         &v,
		Name:      name,
		timestamp: timestamp,
	}, nil
}

func (v Version) GreaterThan(o *Version) bool {
	r := v.v.Compare(*o.v)
	if r == 0 {
		return v.timestamp > o.timestamp
	}

	return r == 1
}

func (v Version) GreaterThanOrEqual(o *Version) bool {
	r := v.v.Compare(*o.v)
	if r == 0 {
		return v.timestamp >= o.timestamp
	}

	return r >= 0
}

func StringsToVersions(strs []string) ([]Version, error) {
	result := make([]Version, len(strs), len(strs))
	for i, str := range strs {
		ver, err := VersionFromString(str)
		if err != nil {
			return nil, err
		}

		result[i] = *ver
	}

	return result, nil
}

func GreatestVersion(versions []Version) *Version {
	sort.SliceStable(versions, func(i, j int) bool {
		left := versions[i]
		right := versions[j]
		return right.GreaterThan(&left)
	})

	return &versions[len(versions)-1]
}
