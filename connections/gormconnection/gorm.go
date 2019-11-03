package gormconnection

import (
	"time"

	"github.com/jinzhu/gorm"
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
	if !c.DB.HasTable(&DBVersion{}) {
		c.DB.AutoMigrate(&DBVersion{})
	}

	var v []DBVersion
	if err := c.DB.Find(&v).Error; err != nil {
		return nil, err
	}

	vers := make([]string, len(v), len(v))
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
