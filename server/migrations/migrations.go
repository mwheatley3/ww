package migrations

import (
	"fmt"
	"sort"
	"time"

	"github.com/mwheatley3/ww/server/pg"
)

const (
	migrationTable = "schema_migrations"
)

// Func is a migration action
type Func func(DB) error

// M is a migration.  Implementations should have a unique name and a function each for Up and Down
type m struct {
	Name string
	Up   func(DB) error
	Down func(DB) error
}

type migrations []*m

func (m migrations) Len() int {
	return len(m)
}

func (m migrations) Less(i, j int) bool {
	return m[i].Name < m[j].Name
}

func (m migrations) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

// NewSet constructs a new Set
func NewSet() *Set {
	return &Set{
		migrations: make(migrations, 0),
	}
}

// A Set is an ordered collection of migrations
type Set struct {
	migrations migrations
	db         *pg.Db
}

// Add a new Migration to the Set
func (ms *Set) Add(name string, up Func, down Func) {
	ms.migrations = append(ms.migrations, &m{name, up, down})
}

// AddSet clones a Set's migrations
func (ms *Set) AddSet(s *Set) {
	ms.migrations = append(ms.migrations, s.migrations...)
}

// Get finds a Migration by name
func (ms *Set) Get(name string) (Func, Func) {
	for _, m := range ms.migrations {
		if m.Name == name {
			return m.Up, m.Down
		}
	}

	return nil, nil
}

// Init will prep the Set for running (Migrate/Rollback). This method will
//  - connect to the db
//  - make sure that the schema record table exists
func (ms *Set) Init(db *pg.Db) error {
	ms.db = db

	if err := ms.db.Connect(); err != nil {
		return err
	}

	if err := ms.ensureSchemaTable(); err != nil {
		return err
	}

	return nil
}

// Migrate runs `count` migrations starting from one past the last run migration
func (ms *Set) Migrate(count int) error {
	return ms.db.WithTxn(func(d *pg.Db) error {
		run, err := ms.toMigrate(d)

		if err != nil {
			return err
		}

		if count > 0 {
			if count > len(run) {
				count = len(run)
			}

			run = run[0:count]
		}

		for _, m := range run {
			fmt.Printf("Processing [%s] [%s]\n", "migrate", m.Name)

			if err := m.Up((*pgDB)(d)); err != nil {
				return err
			}

			if err := insertMigration(d, m.Name); err != nil {
				return err
			}
		}

		return nil
	})
}

// Rollback runs `count` migrations starting from the last migration run and working backwards
func (ms *Set) Rollback(count int) error {
	return ms.db.WithTxn(func(d *pg.Db) error {
		run, err := ms.toRollback(d)

		if err != nil {
			return err
		}

		if count > 0 {
			if count > len(run) {
				count = len(run)
			}

			run = run[0:count]
		}

		for _, m := range run {
			fmt.Printf("Processing [%s] [%s]\n", "rollback", m.Name)

			if err := m.Down((*pgDB)(d)); err != nil {
				return err
			}

			if err := removeMigration(d, m.Name); err != nil {
				return err
			}
		}

		return nil
	})
}

func (ms *Set) toMigrate(d *pg.Db) (migrations, error) {
	old, err := getOld(d)

	if err != nil {
		return nil, err
	}

	oldNames := make(map[string]struct{})

	for _, m := range old {
		oldNames[m.Version] = struct{}{}
	}

	torun := make(migrations, 0)

	for _, m := range ms.migrations {
		if _, ok := oldNames[m.Name]; !ok {
			torun = append(torun, m)
		}
	}

	sort.Sort(torun)

	return torun, nil
}

func (ms *Set) toRollback(d *pg.Db) (migrations, error) {
	old, err := getOld(d)

	if err != nil {
		return nil, err
	}

	oldNames := make(map[string]struct{})

	for _, m := range old {
		oldNames[m.Version] = struct{}{}
	}

	torun := make(migrations, 0)

	for _, m := range ms.migrations {
		if _, ok := oldNames[m.Name]; ok {
			torun = append(torun, m)
		}
	}

	sort.Sort(sort.Reverse(torun))

	return torun, nil
}

func insertMigration(dc *pg.Db, name string) error {
	_, err := dc.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", pg.NewParams(name))
	return err
}

func removeMigration(dc *pg.Db, name string) error {
	_, err := dc.Exec("DELETE FROM schema_migrations where version = $1", pg.NewParams(name))
	return err

}

func getOld(d *pg.Db) ([]*migration, error) {
	var m []*migration
	err := d.GetMany(&m, "SELECT * FROM "+migrationTable+" ORDER BY version ASC", nil)

	if err != nil {
		return nil, err
	}

	return m, nil
}

func (ms *Set) ensureSchemaTable() error {
	count := 0

	err := ms.db.QueryRowx(`
SELECT count(*)
FROM information_schema.tables
WHERE
	table_catalog = current_database() AND
	table_schema = 'public' AND
	table_name = $1
`, migrationTable).Scan(&count)

	if err != nil {
		return err
	}

	if count == 0 {
		_, err := ms.db.Exec(`
CREATE TABLE `+migrationTable+` (
	version text primary key,
	created_at timestamptz not null default now()
)
`, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

type migration struct {
	Version   string    `db:"version"`
	CreatedAt time.Time `db:"created_at"`
}
