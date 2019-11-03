package migo

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type MigrationLoader interface {
	Load() ([]Migration, error)
}

type SQLMigrationsLoader struct {
	path string
}

func NewSQLMigrationLoader(path string) *SQLMigrationsLoader {
	return &SQLMigrationsLoader{
		path: path,
	}
}

func (l *SQLMigrationsLoader) Load() ([]Migration, error) {
	files := map[string]os.FileInfo{}

	removeSuffix := func(s string) string {
		s = strings.Replace(s, ".up.sql", "", -1)
		return strings.Replace(s, ".down.sql", "", -1)
	}

	err := filepath.Walk(l.path, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".up.sql") || strings.HasSuffix(path, ".down.sql") {
			files[removeSuffix(path)] = info
		}

		return nil
	})
	if err != nil {
		panic(err)
	}

	migrations := []Migration{}
	for fileName, file := range files {
		version, err := VersionFromString(removeSuffix(file.Name()))
		if err != nil {
			panic(err)
		}
		migration := SQLMigration{
			v: *version,
		}

		upData, err := ioutil.ReadFile(fileName + ".up.sql")
		if err != nil {
			panic(err)
		}

		downData, err := ioutil.ReadFile(fileName + ".down.sql")
		if err != nil {
			panic(err)
		}

		migration.UpBuffer = upData
		migration.DownBuffer = downData

		migrations = append(migrations, &migration)
	}

	return migrations, nil
}
