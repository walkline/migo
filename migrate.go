package migo

import (
	"errors"
	"fmt"
	"sort"
	"time"
)

// Migrate is entry point struct that can trigger
// database updating or downgrading
type Migrate struct {
	migrations       []Migration
	c                Connection
	loaders          []MigrationLoader
	migrationsLoaded bool
}

// NewMigrate creates new struct that can start migration.
// `c` is a connection to database, you can use gorm connection.
// `loaders` is migrations loaders, it can be `sql` or `go` migrations loader
func NewMigrate(c Connection, loaders ...MigrationLoader) *Migrate {
	return &Migrate{
		c:       c,
		loaders: loaders,
	}
}

// SetConncetion sets database connection that will be used for migration
func (m *Migrate) SetConncetion(c Connection) {
	m.c = c
}

// Add adds migrations to the pool of pending migrations
func (m *Migrate) Add(migration Migration) error {
	m.migrations = append(m.migrations, migration)
	return nil
}

// UpToLatest loads all needed migrations,
// filters migrations that already applied,
// and starts migration process.
func (m *Migrate) UpToLatest() error {

	err := m.loadMigrations()
	if err != nil {
		return errors.New("can't load migrations " + err.Error())
	}

	verStrs, err := m.c.LoadVersions()
	if err != nil {
		return err
	}

	appliedVers, err := StringsToVersions(verStrs)
	if err != nil {
		return err
	}

	lastVer := GreatestVersion(appliedVers)

	m.ensureThatMigrationsNotLost(lastVer, appliedVers)

	migrationsToApply := m.migrationsAfter(lastVer)
	fmt.Printf("Going to apply %d migration(s)...\n", len(migrationsToApply))

	for _, migration := range migrationsToApply {
		fmt.Printf("Applying '%s' migration... \n", migration.Version())
		start := time.Now()

		migration.SetConnection(m.c)
		err := migration.Up()
		if err != nil {
			return err
		}

		err = m.c.SetVersion(migration.Version().String())
		if err != nil {
			return err
		}

		fmt.Printf("Applied '%s'! Duration: %v.\n\n", migration.Version(), time.Since(start))
	}

	fmt.Println("Database up to date!")

	m.migrations = []Migration{}

	return nil
}

// DownWithSteps runs migrations to downgrade database version.
// `steps` is number of latest migrations that needs to be unapplied.
func (m *Migrate) DownWithSteps(steps int) error {

	err := m.loadMigrations()
	if err != nil {
		return errors.New("can't load migrations " + err.Error())
	}

	verStrs, err := m.c.LoadVersions()
	if err != nil {
		return err
	}

	appliedVers, err := StringsToVersions(verStrs)
	if err != nil {
		return err
	}

	lastVer := GreatestVersion(appliedVers)

	migrationsToApply := m.migrationsBefore(lastVer)
	migrationsToApplyCount := steps
	if steps > len(migrationsToApply) {
		migrationsToApplyCount = len(migrationsToApply)
	}

	fmt.Printf("Going to discard %d migration(s)...\n", len(migrationsToApply))
	for i := 0; i < migrationsToApplyCount; i++ {
		migration := migrationsToApply[i]

		fmt.Printf("Discarding '%s' migration... \n", migration.Version())

		migration.SetConnection(m.c)
		err := migration.Down()
		if err != nil {
			return err
		}

		err = m.c.SetVersion(migration.Version().String())
		if err != nil {
			return err
		}

		fmt.Printf("Discarded '%s'!\n\n", migration.Version())
	}

	return nil
}

func (m *Migrate) loadMigrations() error {
	if m.migrationsLoaded {
		return nil
	}

	for _, loader := range m.loaders {
		migrations, err := loader.Load()
		if err != nil {
			return err
		}

		for _, migration := range migrations {
			m.Add(migration)
		}
	}

	m.migrationsLoaded = true

	return nil
}

func (m *Migrate) migrationsAfter(v *Version) []Migration {
	ms := m.sort(m.migrations, true)
	for i := range ms {
		if m.migrations[i].Version().GreaterThan(v) {
			return m.migrations[i:]
		}
	}

	return []Migration{}
}

func (m *Migrate) migrationsBefore(v *Version) []Migration {
	ms := m.sort(m.migrations, false)
	for i := range ms {
		otherV := m.migrations[i].Version()
		if v.GreaterThanOrEqual(&otherV) {
			return m.migrations[i:]
		}
	}

	return []Migration{}
}

func (m *Migrate) sort(ms []Migration, asc bool) []Migration {
	sort.SliceStable(ms, func(i, j int) bool {
		left := m.migrations[i].Version()
		right := m.migrations[j].Version()

		if asc {
			return right.GreaterThan(&left)
		}
		return left.GreaterThan(&right)
	})

	return ms
}

func (m *Migrate) ensureThatMigrationsNotLost(lastVersion *Version, appliedVersions []Version) {
	migsBefore := m.migrationsBefore(lastVersion)

	for _, mig := range migsBefore {
		found := false
		for _, v := range appliedVersions {
			if v.String() == mig.Version().String() {
				found = true
				break
			}
		}

		if !found {
			panic(fmt.Sprintf("version '%s' not found in applied migrations", mig.Version()))
		}
	}
}

func (m *Migrate) lostMigrations(lastVersion *Version, appliedVersions []Version) []Migration {
	migsBefore := m.migrationsBefore(lastVersion)
	lost := []Migration{}
	for _, mig := range migsBefore {
		found := false
		for _, v := range appliedVersions {
			if v.String() == mig.Version().String() {
				found = true
				break
			}
		}

		if !found {
			lost = append(lost, mig)
		}
	}

	lost = m.sort(lost, true)

	return lost
}

func (m *Migrate) ApplyLostAndPanic() error {
	err := m.applyLost(5 * time.Second)
	if err != nil {
		return err
	}

	panic("remove ApplyLostAndPanic code")
}

func (m *Migrate) applyLost(delayBetweenMigrations time.Duration) error {
	err := m.loadMigrations()
	if err != nil {
		return errors.New("can't load migrations " + err.Error())
	}

	verStrs, err := m.c.LoadVersions()
	if err != nil {
		return err
	}

	appliedVers, err := StringsToVersions(verStrs)
	if err != nil {
		return err
	}

	lastVer := GreatestVersion(appliedVers)

	migrationsToApply := m.lostMigrations(lastVer, appliedVers)

	fmt.Printf("Found %d lost migration(s)...\n", len(migrationsToApply))

	for _, migration := range migrationsToApply {
		fmt.Printf("Applying '%s' migration... \n", migration.Version())
		if delayBetweenMigrations > 0 {
			time.Sleep(delayBetweenMigrations)
		}

		start := time.Now()

		migration.SetConnection(m.c)
		err := migration.Up()
		if err != nil {
			return err
		}

		err = m.c.SetVersion(migration.Version().String())
		if err != nil {
			return err
		}

		fmt.Printf("Applied '%s'! Duration: %v.\n\n", migration.Version(), time.Since(start))
	}

	fmt.Println("Database up to date!")

	m.migrations = []Migration{}

	return nil
}
