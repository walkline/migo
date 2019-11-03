package migo

type GoMigrationLoader struct {
	m []Migration
}

var DefaultGoMigrationLoader = &GoMigrationLoader{
	m: []Migration{},
}

func (l *GoMigrationLoader) Add(m Migration) {
	l.m = append(l.m, m)
}

func (l *GoMigrationLoader) Load() ([]Migration, error) {
	return l.m, nil
}

func (l *GoMigrationLoader) Clear() {
	l.m = []Migration{}
}
