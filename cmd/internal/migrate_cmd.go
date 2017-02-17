package internal

import (
	"github.com/mwheatley3/ww/server/log"
	m "github.com/mwheatley3/ww/server/migrations"
	"github.com/mwheatley3/ww/server/pg"
	"github.com/spf13/cobra"
)

// MigrateCmd returns a command for running migrations
func MigrateCmd(loadConf PGLoadConf, path string, set *m.Set) *cobra.Command {
	return m.Cmd(path, "migrations", func() (*m.Set, *pg.Db, error) {
		var (
			pgConf, logConf = loadConf()
			l               = log.New(logConf)
		)

		return set, pg.NewDb(l, pgConf), nil
	})
}
