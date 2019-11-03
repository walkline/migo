package migo

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	version "github.com/hashicorp/go-version"
)

type Version struct {
	v    *version.Version
	Name string
}

func (v Version) String() string {
	segs := v.v.Segments64()
	skipSegs := true
	for i := 1; i < len(segs); i++ {
		if segs[i] != 0 {
			skipSegs = false
			break
		}

	}

	ver := ""
	if skipSegs {
		ver = fmt.Sprintf("%v", segs[0])
	} else {
		ver = v.v.String()
	}

	return fmt.Sprintf("%s-%s", ver, v.Name)
}

func VersionFromString(s string) (*Version, error) {
	strs := strings.Split(s, "-")
	if len(strs) < 2 {
		return nil, errors.New("bad version format")
	}

	v, err := version.NewVersion(strs[0])
	if err != nil {
		return nil, err
	}

	strAsByte := []byte(s)
	name := string(strAsByte[len(strs[0])+1:])

	return &Version{
		v:    v,
		Name: name,
	}, nil
}

func (v Version) GreaterThan(o *Version) bool {
	return v.v.GreaterThan(o.v)
}

func (v Version) GreaterThanOrEqual(o *Version) bool {
	return v.v.GreaterThanOrEqual(o.v)
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
