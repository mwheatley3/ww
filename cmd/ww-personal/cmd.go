package main

import (
	"fmt"
	"os"

	// "github.com/conversable/woodhouse/cmd/internal"
	// "github.com/conversable/woodhouse/server/portal/api/db/migrations"
	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use: "personal",
	}

	cmd.PersistentFlags().StringVar(&confPath, "config", DefaultConfPath, "path to config file")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	}

	cmd.AddCommand(
		webCmd(),
		// createUserCmd(),
		// createProjectCmd(),
		// internal.PsqlCmd(pgLoadConfig),
		// internal.MigrateCmd(pgLoadConfig, "server/portal/api/db/migrations", migrations.Migrations),
	)

	if err := cmd.Execute(); err != nil {
		fmt.Printf("Run error: %s\n", err)
		os.Exit(1)
	}
}
