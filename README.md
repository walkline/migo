## Migo

`migo` is database migration tool.
It supports `.sql` and `.go` type migrations.
For now it have only `gorm` database connector.

## Installation

`migo` has command line tool that creates new migrations. For installation use the next command: `$ go get github.com/walkline/migo/cmd/migo`

## Usage

First of all you need to create a new migration.
```
cd path_with_migrations
migo new sql "create user table"
```

As result you will have two files:
```
1.0.0-create-user-table.up.sql # sql migration that will be used when we will want update database
1.0.0-create-user-table.down.sql # sql migration that will be used when we will want downgrade database
```

Than inside of your project you would need something simmilar to this:
```
package main

import (
	"os"

	"github.com/walkline/migo"
	"github.com/walkline/migo/connections/gormconnection"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func Migrate(db *gorm.DB) {
	m := migo.NewMigrate(
		gormconnection.NewConnection(db),
		migo.NewSQLMigrationLoader(path),
		migo.DefaultGoMigrationLoader,
	)

	err := m.UpToLatest()
	if err != nil {
		panic(err)
	}

	migo.DefaultGoMigrationLoader.Clear()
}

func main() {
	db, _ := gorm.Open("postgres", "host=...")
	Migrate(db)
}
```
