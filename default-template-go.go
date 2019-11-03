package migo

var GoDefaultTemplateData = `// file generated with ./migo/tmpl/go/tmpl.go, feel free to edit it :)
package yourpackagename

import (
	"github.com/jinzhu/gorm"
	"lab.mocintra.com/rostyslav.boyko/MyHomeBackend/src/migo"
	"lab.mocintra.com/rostyslav.boyko/MyHomeBackend/src/migo/connections/gormconnection"
)

func init() {
	migo.DefaultGoMigrationLoader.Add(&Migartion{{.version.verSafe}}{})
}

// Migartion{{.version.verSafe}} implements {{.version.ver}}-{{.version.name}} migration :)
type Migartion{{.version.verSafe}} struct{
	DB *gorm.DB
}

func (m *Migartion{{.version.verSafe}}) Up() error {
	// m.DB ...

	return nil
}

func (m *Migartion{{.version.verSafe}}) Down() error {
	// m.DB ...

	return nil
}

func (m *Migartion{{.version.verSafe}}) Version() migo.Version {
	v, err := migo.VersionFromString("{{.version.ver}}-{{.version.name}}")
	if err != nil {
		panic(err)
	}
	return *v
}

func (m *Migartion{{.version.verSafe}}) SetVersion(v *migo.Version) {
	panic("don't need this")
}

func (m *Migartion{{.version.verSafe}}) SetConnection(c migo.Connection) {
	gormDB, casted := c.(*gormconnection.GormConnection)
	if !casted {
		panic("unk DB type")
	}

	m.DB = gormDB.DB
}
`
