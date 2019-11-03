package migo

import (
	"strings"
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
	c          Connection
	v          Version
	UpBuffer   []byte
	DownBuffer []byte
}

func (m *SQLMigration) SetConnection(c Connection) {
	m.c = c
}

func (m *SQLMigration) Up() error {
	lastI := 0
	potentionalLastI := 0
	for i := 1; i < len(m.UpBuffer); i++ {
		if potentionalLastI <= 0 && (m.UpBuffer[i] == byte(';') || string(m.UpBuffer[i-1:i+1]) == "\\.") {
			potentionalLastI = i
			if i == len(m.UpBuffer)-1 {
				err := m.c.Exec(string(m.UpBuffer[lastI:potentionalLastI]))
				if err != nil {
					return err
				}
			}
			continue
		}

		if potentionalLastI > 0 && (m.UpBuffer[i] == byte('\n')) {
			if strings.TrimSpace(string(m.UpBuffer[potentionalLastI+1:i])) == "" {
				err := m.c.Exec(string(m.UpBuffer[lastI:potentionalLastI]))
				if err != nil {
					return err
				}
				lastI = i
				potentionalLastI = -1
			} else {
				potentionalLastI = -1
			}
		}
	}

	return nil
}

func (m *SQLMigration) Down() error {
	return m.c.Exec(string(m.DownBuffer))
}

func (m *SQLMigration) Version() Version {
	return m.v
}

func (m *SQLMigration) SetVersion(v *Version) {
	m.v = *v
}
