package pg

import (
	"database/sql"
	"github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx"
	"time"
)

type queryer struct {
	*logrus.Logger
	sqlx.Queryer
	sqlx.Execer
	sqlx.Preparer
	slowThreshold time.Duration
}

func (d *queryer) logSlow(prefix, query string, params *Params) func() {
	st := time.Now()

	return func() {
		dur := time.Since(st)

		if dur > d.slowThreshold {
			d.WithFields(logrus.Fields{
				"prefix": prefix,
				"query":  query,
				"params": params.Values(),
			}).Warnf("[PG DB SLOW QUERY] %s > threshold %s", dur, d.slowThreshold)
		}
	}
}

func (d *queryer) logErr(err error, prefix, query string, params *Params) {
	if err != nil {
		// Suppress no rows error if it occurs it's usually handled correctly
		if err == sql.ErrNoRows {
			return
		}

		d.WithFields(logrus.Fields{
			"prefix": prefix,
			"query":  query,
			"params": params.Values(),
		}).Errorf("[PG DB ERROR] %s", err)
	}
}

func (d *queryer) Prepare(query string) (*sqlx.Stmt, error) {
	defer d.logSlow("prepare", query, nil)()

	st, err := sqlx.Preparex(d.Preparer, query)
	d.logErr(err, "prepare", query, nil)

	return st, err
}

func (d *queryer) Get(dest interface{}, query string, params *Params) error {
	defer d.logSlow("get", query, params)()

	var err error
	if qp, ok := dest.(MapProxyer); ok {
		err = d.mapProxyGet(qp, query, params)
	} else {
		err = sqlx.Get(d.Queryer, dest, query, params.Values()...)
	}

	d.logErr(err, "get", query, params)

	return err
}

func (d *queryer) GetMany(dest interface{}, query string, params *Params) error {
	defer d.logSlow("getmany", query, params)()

	var err error
	if smp, ok := dest.(SliceMapProxyer); ok {
		err = d.mapProxyGetMany(smp, query, params)
	} else {
		err = sqlx.Select(d.Queryer, dest, query, params.Values()...)
	}

	d.logErr(err, "getmany", query, params)

	return err
}

func (d *queryer) Exec(query string, params *Params) (sql.Result, error) {
	defer d.logSlow("exec", query, params)()

	res, err := d.Execer.Exec(query, params.Values()...)
	d.logErr(err, "exec", query, params)

	return res, err
}

func (d *queryer) Query(query string, params *Params) (*sqlx.Rows, error) {
	defer d.logSlow("query", query, params)()

	rows, err := d.Queryer.Queryx(query, params.Values()...)
	d.logErr(err, "query", query, params)

	return rows, err
}

func (d *queryer) QueryRow(query string, params *Params) *sqlx.Row {
	return d.Queryer.QueryRowx(query, params.Values()...)
}

func (d *queryer) mapProxyGet(dest MapProxyer, query string, params *Params) error {
	p := dest.MapProxy()
	return mapScan(d.QueryRow(query, params), p)
}

func (d *queryer) mapProxyGetMany(dest SliceMapProxyer, query string, params *Params) error {
	rows, err := d.Query(query, params)

	if err != nil {
		return err
	}

	columns, err := rows.Columns()

	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		p := dest.Append().MapProxy()

		values := make([]interface{}, len(columns))

		for i := range values {
			values[i] = p[columns[i]]
		}

		if err = rows.Scan(values...); err != nil {
			return err
		}
	}

	return nil
}
