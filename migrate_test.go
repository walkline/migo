package migo

import (
	"fmt"
	"testing"
)

type ConnectionMock struct {
	v string
}

func (c *ConnectionMock) Exec(sql string, values ...interface{}) error {
	fmt.Printf("--- Exec sql: %s. Values: %v.\n", sql, values)
	return nil
}

func (c *ConnectionMock) LoadVersions() ([]string, error) {
	return []string{c.v}, nil
}

func (c *ConnectionMock) SetVersion(v string) error {
	c.v = v
	return nil
}

func TestSQLMigrations(t *testing.T) {
	c := &ConnectionMock{}
	c.SetVersion("0-null")

	m := Migrate{}
	m.SetConncetion(c)

	m1 := &SQLMigration{
		UpBuffer:   []byte("UP SQL 1"),
		DownBuffer: []byte("Down SQL 1"),
	}

	v, _ := VersionFromString("1-name")
	m1.SetVersion(v)
	m.Add(m1)

	m2 := &SQLMigration{
		UpBuffer: []byte(`UP SQL 21;
new line;`),
		DownBuffer: []byte("Down SQL 2"),
	}

	v, _ = VersionFromString("2-name")
	m2.SetVersion(v)
	m.Add(m2)

	m.UpToLatest()
}
