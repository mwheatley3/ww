package migrations

import (
	"fmt"
	"os"

	"github.com/mwheatley3/ww/server/pg"
	"github.com/spf13/cobra"
)

// Cmd runs the migrate tool
func Cmd(dir, pkg string, fn func() (*Set, *pg.Db, error)) *cobra.Command {
	cmd := &cobra.Command{
		Use: "migrate",
	}

	var (
		rollback = cmd.Flags().Int("rollback", -10000, "The number of migrations to rollback")
		migrate  = cmd.Flags().Int("migrate", -10000, "The number of migrations to run")
		generate = cmd.Flags().String("generate", "", "The name of the migration to generate")
	)

	cmd.Run = func(c *cobra.Command, args []string) {
		if *generate != "" {
			f, err := gen(dir, pkg, *generate)

			if err != nil {
				fmt.Printf("Gen Error: %s\n", err)
				os.Exit(1)
			}

			fmt.Printf("New migration generated at %s\n", f)
			os.Exit(0)
		}

		up := true
		count := -1

		if *rollback != -10000 {
			up = false
			count = *rollback
		} else if *migrate != -10000 {
			count = *migrate
		}

		set, db, err := fn()

		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}

		if err := run(set, db, up, count); err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
	}

	return cmd
}

func run(set *Set, db *pg.Db, up bool, count int) error {
	if err := set.Init(db); err != nil {
		return err
	}

	if up {
		return set.Migrate(count)
	}

	return set.Rollback(count)
}
