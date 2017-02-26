package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/mwheatley3/ww/server/log"
	"github.com/mwheatley3/ww/server/personal/api/db"
	"github.com/mwheatley3/ww/server/personal/api/service"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

func createUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "create-user",
	}

	cmd.Run = func(cmd *cobra.Command, args []string) {
		var (
			c   = loadConfig()
			l   = log.New(c.Log)
			db  = db.NewFromConfig(l, c.Postgres)
			svc = service.New(l, db)
			r   = bufio.NewReader(os.Stdin)
		)

		fmt.Fprint(os.Stderr, "Email: ")
		email, err := r.ReadString('\n')

		if err != nil {
			l.Fatalf("Email read err: %s", err)
		}

		fmt.Fprint(os.Stderr, "Password: ")
		pw, err := terminal.ReadPassword(int(os.Stdin.Fd()))

		if err != nil {
			l.Fatalf("Password read err: %s", err)
		}

		if err := svc.Init(); err != nil {
			l.Fatalf("init err: %s", err)
		}

		user, err := svc.CreateUser(strings.TrimSpace(email), strings.TrimSpace(string(pw)))

		if err != nil {
			l.Fatalf("create user err: %s", err)
		}

		b, _ := json.MarshalIndent(user, "", "\t")
		fmt.Fprintf(os.Stdout, "\n%s", b)
	}

	return cmd
}
