package migo

import (
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"sort"
	"testing"
)

func TestSQLMigrationLoader(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	dir, _ = ioutil.TempDir(dir, "tst")

	filesName := []string{
		"1.1.0-.up.sql",
		"1.1.0-.down.sql",
		"1.0.1-.up.sql",
		"1.0.1-.down.sql",
		"1.2.1-.up.sql",
		"1.2.1-.down.sql",
	}
	files := []*os.File{}
	for _, f := range filesName {
		file, err := os.Create(dir + "/" + f)
		if err != nil {
			t.Error(err)
		}

		files = append(files, file)
	}

	defer func() {
		for _, file := range files {
			os.Remove(file.Name())
		}

		os.Remove(dir)
	}()

	migs, _ := NewSQLMigrationLoader(dir).Load()

	sort.SliceStable(migs, func(i, j int) bool {
		left := migs[i].Version()
		right := migs[j].Version()

		return right.GreaterThan(&left)
	})

	if len(migs) != 3 ||
		migs[0].Version().String() != "1.0.1-" ||
		migs[2].Version().String() != "1.2.1-" {

		t.Error("bad migrations")
	}
}
