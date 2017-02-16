package main

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

func webCmd() *cobra.Command {
	return &cobra.Command{
		Use: "web",
		Run: func(cmd *cobra.Command, args []string) {
			// var (
			// 	c   = loadConfig()
			// 	l   = log.New(c.Log)
			// 	db  = db.NewFromConfig(l, c.Postgres)
			// 	svc = service.New(l, db)
			// 	srv = web.NewServer(l, svc, c.Web.Config)
			// )

			// if err := srv.Listen(); err != nil {
			// 	l.Fatalf("server listen error: %s", err)
			// }
			http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
			})
			http.ListenAndServe(":8080", nil)
		},
	}
}
