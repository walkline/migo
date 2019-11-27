package migo

import "testing"

func TestGreaterOrEqualVersion(t *testing.T) {
	for _, v := range []struct{ Greater, Less string }{
		{"2.0.0-desc", "1.0.0-desc"},
		{"2.0.1-desc", "2.0.0-desc"},
		{"1.0.0-desc", "0.9.0-desc"},
		{"0.1.1-desc", "0.1.0-desc"},
		{"0.1.1-desc", "0.1.0-desc"},
		{"0.0.2-desc", "0.0.1-desc"},
		{"0.2.2-desc", "0.2.1-desc"},
		{"1.0.0-desc", "1.0.0-desc"},
		{"1.2.0-desc", "1.2.0-hm"},
		{"1.2.1-desc", "1.2.1-hm"},
		{"0.0.1-desc", "0.0.1-desc"},
	} {
		g, err := VersionFromString(v.Greater)
		if err != nil {
			t.Error(err)
		}

		l, err := VersionFromString(v.Less)
		if err != nil {
			t.Error(err)
		}

		if !g.GreaterThanOrEqual(l) {
			t.Error(g, "should be greater then", l, "or equal")
		}
	}
}

func TestVersionString(t *testing.T) {
	source := "1.2.3.4-desc"
	v, err := VersionFromString(source)
	if err != nil {
		t.Error(err)
	}

	if v.String() != source {
		t.Error("wrong response")
	}
}

func TestGreatestVersion(t *testing.T) {
	for _, v := range []struct {
		slice []string
		res   string
	}{
		{[]string{"2.0.0-desc", "1.0.0-desc"}, "2-desc"},
		{[]string{"1.20.1-desc", "2.0.0-desc", "1.0.0-desc"}, "2-desc"},
		{[]string{"1.20.1-desc", "2.0.0-desc", "2.0.0.1-desc"}, "2.0.0.1-desc"},
		{[]string{"1.20.1-desc", "2.0.0-desc", "2.0.0.1-desc", "2.0.0.0-desc"}, "2.0.0.1-desc"},
		{[]string{"1.20.1-desc", "2.0.0-desc", "2.0.0.21-desc", "2.0.0.1-desc"}, "2.0.0.21-desc"},
	} {
		versions := []Version{}
		for _, v := range v.slice {
			v, err := VersionFromString(v)
			if err != nil {
				t.Error(err)
			}
			versions = append(versions, *v)
		}

		result := GreatestVersion(versions)
		if result.String() != v.res {
			t.Error("expected", v.res, "have", result)
		}
	}
}
