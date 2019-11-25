package migo

import (
	"os"
	"testing"
)

type ConnectionMock struct {
	v    string
	sqls []string
}

func (c *ConnectionMock) Exec(sql string, values ...interface{}) error {
	c.sqls = append(c.sqls, sql)
	return nil
}

func (c *ConnectionMock) LoadVersions() ([]string, error) {
	return []string{c.v}, nil
}

func (c *ConnectionMock) SetVersion(v string) error {
	c.v = v
	return nil
}

func (c *ConnectionMock) Commit() error {
	return nil
}

func (c *ConnectionMock) Rollback() error {
	return nil
}

func (c *ConnectionMock) Tx() (Transaction, error) {
	return c, nil
}

func TestUpToLatestSQLMigration(t *testing.T) {
	up1File, err := os.Create("up1.sql")
	if err != nil {
		t.Error(err)
	}

	down1File, err := os.Create("down1.sql")
	if err != nil {
		t.Error(err)
	}

	up2File, err := os.Create("up2.sql")
	if err != nil {
		t.Error(err)
	}

	down2File, err := os.Create("down2.sql")
	if err != nil {
		t.Error(err)
	}

	defer func() {
		os.Remove("up1.sql")
		os.Remove("down1.sql")
		os.Remove("up2.sql")
		os.Remove("down2.sql")
	}()

	up1File.WriteString("UP SQL 1;")
	up2File.WriteString("UP SQL 21; new line;")
	up1File.Seek(0, os.SEEK_SET)
	up2File.Seek(0, os.SEEK_SET)

	c := &ConnectionMock{}
	c.SetVersion("0-null")

	m := Migrate{}
	m.SetConncetion(c)

	m1 := &SQLMigration{
		UpFile:   up1File,
		DownFile: down1File,
	}

	v, _ := VersionFromString("1-name")
	m1.SetVersion(v)
	m.Add(m1)

	m2 := &SQLMigration{
		UpFile:   up2File,
		DownFile: down2File,
	}

	v, _ = VersionFromString("2-name")
	m2.SetVersion(v)
	m.Add(m2)

	err = m.UpToLatest()
	if err != nil {
		t.Error(err)
	}

	if c.sqls[0] != "UP SQL 1;" {
		t.Error("bad sql")
	}

	if c.sqls[1] != "UP SQL 21;" {
		t.Error("bad sql")
	}

	if c.sqls[2] != "new line;" {
		t.Error("bad sql")
	}
}

func TestDownWithStepsSQLMigration(t *testing.T) {
	up1File, err := os.Create("up1.sql")
	if err != nil {
		t.Error(err)
	}

	down1File, err := os.Create("down1.sql")
	if err != nil {
		t.Error(err)
	}

	up2File, err := os.Create("up2.sql")
	if err != nil {
		t.Error(err)
	}

	down2File, err := os.Create("down2.sql")
	if err != nil {
		t.Error(err)
	}

	defer func() {
		os.Remove("up1.sql")
		os.Remove("down1.sql")
		os.Remove("up2.sql")
		os.Remove("down2.sql")
	}()

	down1File.WriteString("Down SQL 1;")
	down2File.WriteString("Down SQL 2;")
	down1File.Seek(0, os.SEEK_SET)
	down2File.Seek(0, os.SEEK_SET)

	c := &ConnectionMock{}
	c.SetVersion("0-null")

	m := Migrate{}
	m.SetConncetion(c)

	m1 := &SQLMigration{
		UpFile:   up1File,
		DownFile: down1File,
	}

	v, _ := VersionFromString("1-name")
	m1.SetVersion(v)
	m.Add(m1)

	m2 := &SQLMigration{
		UpFile:   up2File,
		DownFile: down2File,
	}

	v, _ = VersionFromString("2-name")
	m2.SetVersion(v)
	m.Add(m2)

	c.SetVersion(v.String())

	err = m.DownWithSteps(1)
	if err != nil {
		t.Error(err)
	}

	if len(c.sqls) != 1 {
		t.Error("expected only 1 migration")
	}

	if c.sqls[0] != "Down SQL 2;" {
		t.Error("bad sql")
	}
}
