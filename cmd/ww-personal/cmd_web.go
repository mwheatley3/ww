package main

import (
	"fmt"
	"net/http"
	"os"

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
			port := os.Getenv("PORT")
			if port == "" {
				port = "8081"
			}
			port = ":" + port
			http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
			})
			http.ListenAndServe(port, nil)
		},
	}
}
