package gormconnection

import (
	"time"

	"github.com/walkline/migo"
	"gorm.io/gorm"
)

const (
	NoVer = "0-null"
)

type DBVersion struct {
	Date    time.Time
	Version string
}

type GormConnection struct {
	DB *gorm.DB
}

func NewConnection(c *gorm.DB) *GormConnection {
	return &GormConnection{
		DB: c,
	}
}

func (c *GormConnection) Exec(sql string, values ...interface{}) error {
	return c.DB.Exec(sql, values...).Error
}

func (c *GormConnection) LoadVersions() ([]string, error) {
	var v []DBVersion
	if err := c.DB.Find(&v).Error; err != nil {
		// maybe we need to create table...
		if err = c.DB.AutoMigrate(&DBVersion{}); err != nil {
			return nil, err
		}

		if err = c.DB.Find(&v).Error; err != nil {
			return nil, err
		}
	}

	vers := make([]string, len(v))
	for i, ver := range v {
		vers[i] = ver.Version
	}

	if len(vers) == 0 {
		vers = append(vers, NoVer)
	}

	return vers, nil
}

func (c *GormConnection) SetVersion(ver string) error {
	return c.DB.Save(&DBVersion{
		Date:    time.Now(),
		Version: ver,
	}).Error
}

func (c *GormConnection) Tx() (migo.Transaction, error) {
	return NewTransaction(c.DB)
}

type GormTransaction struct {
	DB *gorm.DB
}

func NewTransaction(db *gorm.DB) (*GormTransaction, error) {
	return &GormTransaction{
		DB: db.Begin(),
	}, nil
}

func (tx *GormTransaction) Exec(sql string, values ...interface{}) error {
	return tx.DB.Exec(sql, values...).Error
}

func (tx *GormTransaction) Commit() error {
	return tx.DB.Commit().Error
}

func (tx *GormTransaction) Rollback() error {
	return tx.DB.Rollback().Error
}
