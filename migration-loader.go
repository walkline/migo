package migo

type MigrationsLoader interface {
	Load() ([]Migration, error)
}
