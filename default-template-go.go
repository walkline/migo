package migo

var GoDefaultTemplateData = `// file generated with ./migo/tmpl/go/tmpl.go, feel free to edit it :)
package yourpackagename

import (
	"github.com/jinzhu/gorm"
	"github.com/walkline/migo"
	"github.com/walkline/migo/connections/gormconnection"
)

func init() {
	migo.DefaultGoMigrationLoader.Add(&Migration{{.version.verSafe}}{})
}

// Migration{{.version.verSafe}} implements {{.version.ver}}-{{.version.name}} migration :)
type Migration{{.version.verSafe}} struct{
	DB *gorm.DB
}

func (m *Migration{{.version.verSafe}}) Up() error {
	// m.DB ...

	return nil
}

func (m *Migration{{.version.verSafe}}) Down() error {
	// m.DB ...

	return nil
}

func (m *Migration{{.version.verSafe}}) Version() migo.Version {
	v, err := migo.VersionFromString("{{.version.ver}}-{{.version.name}}")
	if err != nil {
		panic(err)
	}
	return *v
}

func (m *Migration{{.version.verSafe}}) SetVersion(v *migo.Version) {
	panic("don't need this")
}

func (m *Migration{{.version.verSafe}}) SetConnection(c migo.Connection) {
	gormDB, casted := c.(*gormconnection.GormConnection)
	if !casted {
		panic("unk DB type")
	}

	m.DB = gormDB.DB
}
`
