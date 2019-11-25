package migo

import (
	"os"

	"github.com/walkline/migo/sqlscanner"
)

type Connection interface {
	Exec(sql string, values ...interface{}) error
	LoadVersions() ([]string, error)
	SetVersion(v string) error
	Tx() (Transaction, error)
}

type Transaction interface {
	Exec(sql string, values ...interface{}) error
	Commit() error
	Rollback() error
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

	tx, err := m.c.Tx()
	if err != nil {
		return err
	}

	scanner := sqlscanner.NewSQLScanner(m.UpFile)
	query := ""
	for scanner.Next(&query) {
		err := tx.Exec(query)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if scanner.Error != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (m *SQLMigration) Down() error {
	defer m.UpFile.Close()
	defer m.DownFile.Close()

	tx, err := m.c.Tx()
	if err != nil {
		return err
	}

	scanner := sqlscanner.NewSQLScanner(m.DownFile)
	query := ""
	for scanner.Next(&query) {
		err := m.c.Exec(query)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if scanner.Error != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (m *SQLMigration) Version() Version {
	return m.v
}

func (m *SQLMigration) SetVersion(v *Version) {
	m.v = *v
}
