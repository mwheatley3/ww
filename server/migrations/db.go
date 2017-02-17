package migrations

import "github.com/mwheatley3/ww/server/pg"

// A DB abstracts away the underlying db resource
type DB interface {
	Exec(string, ...interface{}) error
	ExecMany(...Query) error
}

type pgDB pg.Db

func (d *pgDB) Exec(q string, p ...interface{}) error {
	_, err := (*pg.Db)(d).Exec(q, pg.NewParams(p...))
	return err
}

func (d *pgDB) ExecMany(queries ...Query) error {
	for _, q := range queries {
		if err := d.Exec(q.Query(), q.Params()...); err != nil {
			return err
		}
	}

	return nil
}

// A Query instance represents a single db
// query
type Query interface {
	Query() string
	Params() []interface{}
}

type q struct {
	query  string
	params []interface{}
}

func (q q) Query() string         { return q.query }
func (q q) Params() []interface{} { return q.params }

// Q is a helper function to generate a Query
func Q(query string, params ...interface{}) Query {
	return q{query, params}
}
