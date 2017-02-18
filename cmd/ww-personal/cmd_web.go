package main

import (
	"fmt"

	"github.com/mwheatley3/ww/server/log"

	"github.com/mwheatley3/ww/server/personal/api/db"
	"github.com/mwheatley3/ww/server/personal/api/service"
	"github.com/mwheatley3/ww/server/personal/web"
	"github.com/spf13/cobra"
)

func webCmd() *cobra.Command {
	return &cobra.Command{
		Use: "web",
		Run: func(cmd *cobra.Command, args []string) {
			var (
				c   = loadConfig()
				l   = log.New(c.Log)
				db  = db.NewFromConfig(l, c.Postgres)
				svc = service.New(l, db)
				srv = web.NewServer(l, svc, c.Web.Config)
			)

			fmt.Printf("%#+v\n", c)

			if err := srv.Listen(); err != nil {
				l.Fatalf("server listen error: %s", err)
			}
			// port := os.Getenv("PORT")
			// if port == "" {
			// 	port = "8081"
			// }
			// port = ":" + port
			// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			// 	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
			// })
			// http.ListenAndServe(port, nil)
		},
	}
}
