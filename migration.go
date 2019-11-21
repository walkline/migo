package migo

import (
	"os"

	"github.com/walkline/migo/sqlscanner"
)

type Connection interface {
	Exec(sql string, values ...interface{}) error
	LoadVersions() ([]string, error)
	SetVersion(v string) error
}

type Migration interface {
	SetConnection(c Connection)
	Version() Version
	SetVersion(v *Version)
	Up() error
	Down() error
}

type SQLMigration struct {
	c        Connection
	v        Version
	UpFile   *os.File
	DownFile *os.File
}

func (m *SQLMigration) SetConnection(c Connection) {
	m.c = c
}

func (m *SQLMigration) Up() error {
	defer m.UpFile.Close()
	defer m.DownFile.Close()

	scanner := sqlscanner.NewSQLScanner(m.UpFile)
	query := ""
	for scanner.Next(&query) {
		err := m.c.Exec(query)
		if err != nil {
			return err
		}
	}

	return scanner.Error
}

func (m *SQLMigration) Down() error {
	defer m.UpFile.Close()
	defer m.DownFile.Close()

	scanner := sqlscanner.NewSQLScanner(m.DownFile)
	query := ""
	for scanner.Next(&query) {
		err := m.c.Exec(query)
		if err != nil {
			return err
		}
	}

	return scanner.Error
}

func (m *SQLMigration) Version() Version {
	return m.v
}

func (m *SQLMigration) SetVersion(v *Version) {
	m.v = *v
}
